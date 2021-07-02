package date

import (
	"context"
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
		value := waitingForValue(dependingScope["key"].(string), fieldsMap)
		if value == nil {
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

func Validate(field interface{}, fieldValidator *fields.FieldValidator, fieldsMap map[string]interface{}) interface{} {
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

	if err = json.Unmarshal([]byte(fieldValidator.Patterns), &datePatterns); err != nil {
		return nil
	}

	if fieldDate, err = time.Parse("2006-01-02", strField); err != nil {
		return nil
	}

	pattern := datePatterns[0]

	for _, maxPattern := range pattern.Max {
		switch maxPattern.PatternType {
		case "date":
			if !validateMaxDate(fieldDate, maxPattern.Value) {
				return nil
			}
		case "now":
			if !validateMaxNow(fieldDate, maxPattern.Value) {
				return nil
			}
		case "depending_formula":
			if !validateMaxDependingFormula(fieldDate, maxPattern.Value, fieldsMap) {
				return nil
			}
		default:
			log.Logger.Errorf("unknown date type: %s", maxPattern.PatternType)
			return nil
		}
	}
	for _, minPattern := range pattern.Min {
		switch minPattern.PatternType {
		case "date":
			if !validateMinDate(fieldDate, minPattern.Value) {
				return nil
			}
		case "now":
			if !validateMinNow(fieldDate, minPattern.Value) {
				return nil
			}
		case "depending":
			if !validateMinDepending(fieldDate, minPattern.Value, fieldsMap) {
				return nil
			}
		case "depending_formula":
			// кейс для валидации минимального паттерна с формулой и зависимостью
			if !validateMinDependingFormula(fieldDate, minPattern.Value, fieldsMap) {
				return nil
			}
		default:
			log.Logger.Errorf("unknown date type: %s", minPattern.PatternType)
			return nil
		}
	}

	return fieldDate
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
	dependingValue := waitingForValue(value.Key, fieldsMap)
	if dependingValue == nil {
		log.Logger.Errorf("depending field not found: %s", value.Key)
		return false
	}
	expectedDate := dependingValue.(time.Time)

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

func waitingForValue(key string, fieldsMap map[string]interface{}) interface{} {
	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	dependingValue := waitForValueWithTimeout(ctx, key, fieldsMap)
	cancel()
	return dependingValue
}

func waitForValueWithTimeout(ctx context.Context, key string, fieldsMap map[string]interface{}) interface{} {
	waitingChannel := make(chan interface{})
	go waitForValue(key, fieldsMap, waitingChannel)
	select {
	case <-ctx.Done():
		return nil
	case value := <-waitingChannel:
		return value
	}
}

func waitForValue(key string, fieldsMap map[string]interface{}, waitingChannel chan interface{}) {
	// ожидание получения значения из мапы
	ticker := time.NewTimer(10 * time.Millisecond)
	checkingChannel := make(chan interface{})
	go checkField(key, fieldsMap, checkingChannel)
	for {
		select {
		case value := <-checkingChannel:
			waitingChannel <- value
			return
		case <-ticker.C:
			continue
		}
	}
}

func checkField(key string, fieldsMap map[string]interface{}, checkingChannel chan interface{}) {
	if value, ok := fieldsMap[key]; ok {
		checkingChannel <- value
	}
}
