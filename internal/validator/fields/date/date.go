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

func Validate(field string, fieldValidator *fields.FieldValidator) bool {

	fmt.Printf("man_year: %v\n", car.Manufacturing_year)
	fmt.Printf("cred_year: %v\n", car.Credential_issue_date)

	var (
		fieldDate    time.Time
		datePatterns []*DatePattern
		err          error
	)

	err = json.Unmarshal([]byte(fieldValidator.Patterns), &datePatterns)
	if err != nil {
		// println(err)
		return false
	}

	fieldDate, err = time.Parse("2006-01-02", field)
	if err != nil {
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

	err = reflections.SetField(&car, strings.Title(fieldValidator.FieldName), fieldDate)
	if err != nil {
		log.Logger.Error(err)
	}

	return true
}

func validateMinDate(fieldDate time.Time, rawValue json.RawMessage) bool {
	var value DateDateValue
	json.Unmarshal(rawValue, &value)

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
