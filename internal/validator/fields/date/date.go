package date

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"validation_service/internal/models"
	"validation_service/internal/validator/fields"
	"validation_service/pkg/log"

	"github.com/oleiade/reflections"
)

type DateDateValue string

type DateDependingValue struct {
	Scope string `json:"scope"`
	Key   string `json:"key"`
}

type dependency struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (d *dependency) getInitialDate() (time.Time, error) {
	switch d.Type {
	case "now":
		return time.Now(), nil
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

func (f *DateDependingFormulaValue) getExpectedDate() (time.Time, error) {
	var (
		years, months, days int
	)

	initialDate, err := f.Dependency.getInitialDate()
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

var car models.Car

func init() {
	car = models.Car{} // TODO: убрать это. Одно ТС на один запрос, а не навсегда
}

func Validate(field interface{}, fieldValidator *fields.FieldValidator) bool {
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
		return false
	}

	if fieldDate, err = time.Parse("2006-01-02", strField); err != nil {
		return false
	}

	pattern := datePatterns[0]

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
			if !validateMaxDependingFormula(fieldDate, maxPattern.Value) {
				return false
			}
		default:
			log.Logger.Errorf("unknown date type: %s", maxPattern.PatternType)
			return false
		}
	}
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
			if !validateMinDepending(fieldDate, minPattern.Value) {
				return false
			}
		default:
			log.Logger.Errorf("unknown date type: %s", minPattern.PatternType)
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

	return fieldDate.After(expectedDate)
}

func validateMinNow(fieldDate time.Time, rawValue json.RawMessage) bool {
	return fieldDate.After(time.Now())
}

func validateMinDepending(fieldDate time.Time, rawValue json.RawMessage) bool {
	var value DateDependingValue
	json.Unmarshal(rawValue, &value)

	if value.Scope == "car" {
		expectedDate, err := reflections.GetField(&car, strings.Title(value.Key))
		if err != nil {
			return false
		}
		return fieldDate.After(expectedDate.(time.Time))
	} else {
		log.Logger.Errorf("unknown scope of depending: %s", value.Scope)
		return false
	}
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

func validateMaxDependingFormula(fieldDate time.Time, rawValue json.RawMessage) bool {
	var formula DateDependingFormulaValue
	json.Unmarshal(rawValue, &formula)

	expectedDate, err := formula.getExpectedDate()
	if err != nil {
		log.Logger.Errorf("error on with calculating formula: %s", err.Error())
		return false
	}

	return fieldDate.Before(expectedDate)
}
