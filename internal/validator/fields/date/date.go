package date

import (
	"encoding/json"
	"fmt"
	"time"
	"validation_service/internal/validator/fields"
	"validation_service/pkg/log"
)

type DateDateValue string

type DateDependingValue struct {
	Scope string `json:"scope"`
	Key   string `json:"key"`
}

type dependency struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func (d *dependency) getInitialDate(fieldsMap map[string]interface{}) (time.Time, error) {
	switch d.Type {
	case "now":
		return time.Now(), nil
	case "depending":
		dependingScope := d.Value.(map[string]interface{})
		value, err := waitForValue(dependingScope["key"].(string), fieldsMap)
		if !err {
			return time.Time{}, fmt.Errorf("depending field not found: %v", dependingScope["key"])
		}
		return value.(time.Time), nil
	default:
		return time.Time{}, fmt.Errorf("no logic for dependency: %v", d.Type)
	}
}

type DateDependingFormulaValue struct {
	Dependency dependency `json:"depending"`
	Operation  string     `json:"operation"`
	Value      int        `json:"value"`
	Unit       string     `json:"unit"`
}

func (f *DateDependingFormulaValue) getExpectedDate(fieldsMap map[string]interface{}) (time.Time, error) {
	var (
		years, months, days int
	)

	initialDate, err := f.Dependency.getInitialDate(fieldsMap)
	if err != nil {
		return time.Time{}, err
	}

	switch f.Unit {
	case "year":
		years = f.Value
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

func init() {

}

func Validate(field interface{}, fieldValidator *fields.FieldValidator, fieldsMap map[string]interface{}) (interface{}, bool) {
	var (
		fieldDate    time.Time
		datePatterns []*DatePattern
		err          error
		ok           bool
		strField     string
	)

	if strField, ok = field.(string); !ok {
		log.Logger.Error("type conversion failed")
		return nil, false
	}

	if err = json.Unmarshal([]byte(fieldValidator.Patterns), &datePatterns); err != nil {
		return nil, false
	}

	if fieldDate, err = time.Parse("2006-01-02", strField); err != nil {
		return nil, false
	}

	pattern := datePatterns[0]

	for _, maxPattern := range pattern.Max {
		switch maxPattern.PatternType {
		case "date":
			if !validateMaxDate(fieldDate, maxPattern.Value) {
				return nil, false
			}
		case "now":
			if !validateMaxNow(fieldDate, maxPattern.Value) {
				return nil, false
			}
		case "depending_formula":
			if !validateMaxDependingFormula(fieldDate, maxPattern.Value, fieldsMap) {
				return nil, false
			}
		default:
			log.Logger.Errorf("unknown date type: %s", maxPattern.PatternType)
			return nil, false
		}
	}
	for _, minPattern := range pattern.Min {
		switch minPattern.PatternType {
		case "date":
			if !validateMinDate(fieldDate, minPattern.Value) {
				return nil, false
			}
		case "now":
			if !validateMinNow(fieldDate, minPattern.Value) {
				return nil, false
			}
		case "depending":
			if !validateMinDepending(fieldDate, minPattern.Value, fieldsMap) {
				return nil, false
			}
		case "depending_formula":
			// кейс для валидации минимального паттерна с формулой и зависимостью
			if !validateMinDependingFormula(fieldDate, minPattern.Value, fieldsMap) {
				return nil, false
			}
		default:
			log.Logger.Errorf("unknown date type: %s", minPattern.PatternType)
			return nil, false
		}
	}

	return fieldDate, true
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

	return fieldDate.After(expectedDate)
}

func validateMinNow(fieldDate time.Time, rawValue json.RawMessage) bool {
	return fieldDate.After(time.Now())
}

func validateMinDepending(fieldDate time.Time, rawValue json.RawMessage, fieldsMap map[string]interface{}) bool {
	var value DateDependingValue
	json.Unmarshal(rawValue, &value)

	dependingValue, ok := waitForValue(value.Key, fieldsMap)
	if !ok {
		log.Logger.Errorf("depending field not found: %s", value.Key)
		return false
	}
	var expectedDate time.Time
	if dependingValue != nil {
		expectedDate = dependingValue.(time.Time)
	} else {
		log.Logger.Errorf("depending field not validated: %s", value.Key)
		return false
	}

	return fieldDate.After(expectedDate)
}

func validateMinDependingFormula(fieldDate time.Time, rawValue json.RawMessage, fieldsMap map[string]interface{}) bool {
	var formula DateDependingFormulaValue
	json.Unmarshal(rawValue, &formula)

	expectedDate, err := formula.getExpectedDate(fieldsMap)
	if err != nil {
		log.Logger.Errorf("error on with calculating formula: %s", err.Error())
		return false
	}

	return fieldDate.After(expectedDate)
}

func validateMaxDate(fieldDate time.Time, rawValue json.RawMessage) bool {
	var value DateDateValue
	json.Unmarshal(rawValue, &value)

	expectedDate, err := time.Parse("2006-01-02", string(value))
	if err != nil {
		return false
	}

	return fieldDate.Before(expectedDate)
}

func validateMaxNow(fieldDate time.Time, rawValue json.RawMessage) bool {
	return fieldDate.Before(time.Now())
}

func validateMaxDependingFormula(fieldDate time.Time, rawValue json.RawMessage, fieldsMap map[string]interface{}) bool {
	var formula DateDependingFormulaValue
	json.Unmarshal(rawValue, &formula)

	expectedDate, err := formula.getExpectedDate(fieldsMap)
	if err != nil {
		log.Logger.Errorf("error on with calculating formula: %s", err.Error())
		return false
	}

	return fieldDate.Before(expectedDate)
}

func waitForValue(key string, fieldsMap map[string]interface{}) (interface{}, bool) {
	// ожидание получения значения из мапы
	var beggining time.Time
	beggining = time.Now()
	for true {
		if value, ok := fieldsMap[key]; ok {
			return value, ok
		}
		current := time.Now().Sub(beggining)
		if current.Seconds() > 5 {
			break
		}
		time.Sleep(time.Millisecond * 10)
	}
	return nil, false
}
