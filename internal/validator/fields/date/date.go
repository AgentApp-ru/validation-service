package date

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"validation_service/pkg/log"
)

type (
	DateDateValue  string
	DependingValue struct {
		Scope string `json:"scope"`
		Key   string `json:"key"`
	}

	dependency struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}

	DateDependingFormulaValue struct {
		Dependency dependency `json:"depending"`
		Operation  string     `json:"operation"`
		Value      int        `json:"value"`
		Unit       string     `json:"unit"`
	}

	condition struct {
		Type    string         `json:"type"`
		Items   map[string]int `json:"items"`
		Default int            `json:"default"`
	}

	conditionValue struct {
		Field     DependingValue `json:"field"`
		Condition condition      `json:"condition"`
	}

	DateDependingConditionFormulaValue struct {
		Dependency     dependency     `json:"depending"`
		Operation      string         `json:"operation"`
		ConditionValue conditionValue `json:"value"`
		Unit           string         `json:"unit"`
	}

	Validator struct {
		objectMap    *sync.Map
		allFieldsMap *sync.Map
		errors       chan string
	}
)

func (c *condition) getItem(searchedValue string) (int, error) {
	switch c.Type {
	case "equals":
		item, ok := c.Items[searchedValue]
		if !ok {
			return c.Default, nil
		}
		return item, nil
	default:
		log.Logger.Errorf("unknown condition type: %s", c.Type)
		return c.Default, nil
	}
}

func (c *conditionValue) getItem(selfMap, fieldsMap *sync.Map) (int, error) {
	value := waitingForValue(c.Field.Scope, c.Field.Key, selfMap, fieldsMap)
	if value == nil {
		value = "default"
	}
	searchedValue := fmt.Sprintf("%v", value)

	return c.Condition.getItem(searchedValue)
}

func New(objectMap, allFieldsMap *sync.Map, errors chan string) *Validator {
	return &Validator{
		objectMap:    objectMap,
		allFieldsMap: allFieldsMap,
		errors:       errors,
	}
}

func (d *dependency) getInitialDate(selfMap, fieldsMap *sync.Map) (time.Time, error) {
	switch d.Type {
	case "now":
		return time.Now(), nil
	case "depending":
		var dependingScope DependingValue
		err := json.Unmarshal(d.Value, &dependingScope)
		if err != nil {
			return time.Time{}, err
		}

		value := waitingForValue(dependingScope.Scope, dependingScope.Key, selfMap, fieldsMap)
		if value == nil {
			return time.Time{}, fmt.Errorf("depending field not found: %s/%s", dependingScope.Scope, dependingScope.Key)
		}
		expectedDateRaw := value.(string)
		expectedDate, err := time.Parse("2006-01-02", expectedDateRaw)
		return expectedDate, err
	default:
		return time.Time{}, fmt.Errorf("no logic for dependency: %v", d.Type)
	}
}

func (f *DateDependingFormulaValue) getExpectedDate(selfMap, fieldsMap *sync.Map) (time.Time, error) {
	var (
		years, months, days int
	)

	initialDate, err := f.Dependency.getInitialDate(selfMap, fieldsMap)
	if err != nil {
		return time.Time{}, err
	}

	switch f.Unit {
	case "year":
		years = f.Value
	case "month":
		months = f.Value
	case "day":
		days = f.Value
	}

	switch f.Operation {
	case "subtract":
		years = 0 - years
		months = 0 - months
		days = 0 - days
	case "add":
		years = 0 + years
		months = 0 + months
		days = 0 + days
	}

	if days == 0 {
		days = -1
	}

	return initialDate.AddDate(years, months, days), nil
}

func (f *DateDependingConditionFormulaValue) getExpectedDate(selfMap, fieldsMap *sync.Map) (time.Time, error) {
	var (
		years, months, days int
	)

	initialDate, err := f.Dependency.getInitialDate(selfMap, fieldsMap)
	if err != nil {
		return time.Time{}, err
	}

	item, err := f.ConditionValue.getItem(selfMap, fieldsMap)
	if err != nil {
		return time.Time{}, err
	}

	switch f.Unit {
	case "year":
		years = item
	case "month":
		months = item
	case "day":
		days = item
	}

	switch f.Operation {
	case "subtract":
		years = 0 - years
		months = 0 - months
		days = 0 - days
	case "add":
		years = 0 + years
		months = 0 + months
		days = 0 + days
	}

	if days == 0 {
		days = -1
	}

	return initialDate.AddDate(years, months, days), nil
}

type DateMinMaxPattern struct {
	PatternType string          `json:"type"`
	Value       json.RawMessage `json:"value"`
}

type DatePattern struct {
	Min []*DateMinMaxPattern `json:"min"`
	Max []*DateMinMaxPattern `json:"max"`
}

func (v *Validator) Validate(
	field interface{},
	transformers *json.RawMessage,
	patterns json.RawMessage,
	allowWhiteSpaces bool,
) bool {
	var (
		fieldDate    time.Time
		datePatterns []*DatePattern
		err          error
		ok           bool
		strField     string
	)

	if strField, ok = field.(string); !ok {
		log.Logger.Error("type conversion failed")
		return false
	}

	if err = json.Unmarshal([]byte(patterns), &datePatterns); err != nil {
		return false
	}

	if fieldDate, err = time.Parse("2006-01-02", strField); err != nil {
		return false
	}

	pattern := datePatterns[0]

	for _, minPattern := range pattern.Min {
		switch minPattern.PatternType {
		case "date":
			if !validateMinDate(fieldDate, minPattern.Value) {
				return false
			}
		case "now":
			if !validateMinNow(fieldDate, minPattern.Value) {
				return false
			}
		case "depending":
			if !validateMinDepending(fieldDate, minPattern.Value, v.objectMap, v.allFieldsMap) {
				return false
			}
		case "depending_formula":
			if !validateMinDependingFormula(fieldDate, minPattern.Value, v.objectMap, v.allFieldsMap) {
				return false
			}
		case "depending_condition_formula":
			if !validateMinDependingConditionFormula(fieldDate, minPattern.Value, v.objectMap, v.allFieldsMap) {
				return false
			}
		default:
			log.Logger.Errorf("unknown date type: %s", minPattern.PatternType)
			return false
		}
	}
	for _, maxPattern := range pattern.Max {
		switch maxPattern.PatternType {
		case "date":
			if !validateMaxDate(fieldDate, maxPattern.Value) {
				return false
			}
		case "now":
			if !validateMaxNow(fieldDate, maxPattern.Value) {
				return false
			}
		case "depending_formula":
			if !validateMaxDependingFormula(fieldDate, maxPattern.Value, v.objectMap, v.allFieldsMap) {
				return false
			}
		case "depending_condition_formula":
			if !validateMaxDependingConditionFormula(fieldDate, maxPattern.Value, v.objectMap, v.allFieldsMap) {
				return false
			}
		default:
			log.Logger.Errorf("unknown date type: %s", maxPattern.PatternType)
			return false
		}
	}
	return true
}

func validateMinDate(fieldDate time.Time, rawValue json.RawMessage) bool {
	var value DateDateValue
	err := json.Unmarshal(rawValue, &value)
	if err != nil {
		return false
	}

	expectedDate, err := time.Parse("2006-01-02", string(value))
	if err != nil {
		return false
	}

	return !fieldDate.Before(expectedDate)
}

func validateMinNow(fieldDate time.Time, rawValue json.RawMessage) bool {
	return fieldDate.After(time.Now())
}

func validateMinDepending(fieldDate time.Time, rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) bool {
	var value DependingValue
	json.Unmarshal(rawValue, &value)

	dependingValue := waitingForValue(value.Scope, value.Key, selfMap, fieldsMap)
	if dependingValue == nil {
		log.Logger.Errorf("depending field not found: %s/%s", value.Scope, value.Key)
		return false
	}
	expectedDateRaw := dependingValue.(string)
	expectedDate, err := time.Parse("2006-01-02", expectedDateRaw)
	if err != nil {
		return false
	}

	return !fieldDate.Before(expectedDate)
}

func validateMinDependingFormula(fieldDate time.Time, rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) bool {
	var formula DateDependingFormulaValue
	json.Unmarshal(rawValue, &formula)

	expectedDate, err := formula.getExpectedDate(selfMap, fieldsMap)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return false
	}

	return !fieldDate.Before(expectedDate)
}

func validateMinDependingConditionFormula(fieldDate time.Time, rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) bool {
	var formula DateDependingConditionFormulaValue
	json.Unmarshal(rawValue, &formula)

	expectedDate, err := formula.getExpectedDate(selfMap, fieldsMap)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return false
	}

	return !fieldDate.Before(expectedDate)
}

func validateMaxDependingConditionFormula(fieldDate time.Time, rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) bool {
	var formula DateDependingConditionFormulaValue
	json.Unmarshal(rawValue, &formula)

	expectedDate, err := formula.getExpectedDate(selfMap, fieldsMap)
	println("!!!", fmt.Sprintf("%v", expectedDate), err)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return false
	}

	return !fieldDate.After(expectedDate)
}

func validateMaxDate(fieldDate time.Time, rawValue json.RawMessage) bool {
	var value DateDateValue
	json.Unmarshal(rawValue, &value)

	expectedDate, err := time.Parse("2006-01-02", string(value))
	if err != nil {
		return false
	}

	return !fieldDate.After(expectedDate)
}

func validateMaxNow(fieldDate time.Time, rawValue json.RawMessage) bool {
	return !fieldDate.After(time.Now())
}

func validateMaxDependingFormula(fieldDate time.Time, rawValue json.RawMessage, selfMap, fieldsMap *sync.Map) bool {
	var formula DateDependingFormulaValue
	json.Unmarshal(rawValue, &formula)

	expectedDate, err := formula.getExpectedDate(selfMap, fieldsMap)
	if err != nil {
		log.Logger.Errorf("Ошибка при расчёте формулы: %s", err.Error())
		return false
	}

	return !fieldDate.After(expectedDate)
}

func waitingForValue(scope, key string, selfMap, fieldsMap *sync.Map) interface{} {
	var scopeObjectMap *sync.Map

	if scope == "self" {
		scopeObjectMap = selfMap
	} else {
		scopeObjectMapLoaded, ok := fieldsMap.Load(scope)
		if !ok {
			return nil
		}
		scopeObjectMap = scopeObjectMapLoaded.(*sync.Map)
	}

	end := time.Now().Add(2 * time.Second)
	for time.Now().Before(end) {
		if value, ok := scopeObjectMap.Load(key); ok {
			return value
		}
		time.Sleep(300 * time.Millisecond)
	}

	return nil
}
