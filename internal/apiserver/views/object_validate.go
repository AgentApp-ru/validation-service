package views

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	validator_module "validation_service/internal/validator"
	"validation_service/pkg/log"

	"github.com/oleiade/reflections"
)

type Pattern struct {
	Chars string `json:"chars"`
	Min   int    `json:"min"`
	Max   int    `json:"max"`
}

type StringPattern struct {
	Name               string     `json:"name"`
	Allow_white_spaces bool       `json:"allow_white_spaces"`
	Patterns           []*Pattern `json:"patterns"`
}

type IntPattern struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

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

type FieldValidator struct {
	FieldName string          `json:"field"`
	FieldType string          `json:"type"`
	Patterns  json.RawMessage `json:"patterns"`
	// StringPatterns []*StringPattern `json:"patterns"`
}

type validatorClass struct {
	Schema             string            `json:"$schema"`
	FieldValidators    []*FieldValidator `json:"validators"`
	FieldValidatorsMap map[string]*FieldValidator
}

type Car struct {
	Manufacturing_year    time.Time
	Credential_issue_date time.Time
}

var car Car

func ValidateCar(data map[string]interface{}) (bool, []string) {
	car = Car{}

	fieldsWithErrors := []string{}

	rawValidator, err := validator_module.Validator.GetRaw("car")
	if err != nil {
		return false, fieldsWithErrors
	}

	validator := getValidator(rawValidator)
	if validator == nil {
		return false, fieldsWithErrors
	}

	for k, v := range data {
		println("validating ", k, v)
		fieldValidator, ok := validator.FieldValidatorsMap[k]
		if !ok {
			log.Logger.Info("TODO: что-то сделать")
		}

		ok = validate(v, fieldValidator)
		println(ok)
		if !ok {
			fieldsWithErrors = append(fieldsWithErrors, k)
		}
	}

	return len(fieldsWithErrors) == 0, fieldsWithErrors
}

func getValidator(data []byte) *validatorClass {
	v := &validatorClass{}

	err := json.Unmarshal(data, v)
	if err != nil {
		log.Logger.Error(err)
	}

	v.FieldValidatorsMap = map[string]*FieldValidator{}
	for _, field := range v.FieldValidators {
		v.FieldValidatorsMap[field.FieldName] = field
	}

	return v
}

func validate(field interface{}, fieldValidator *FieldValidator) bool {
	var (
		ok       bool
		err      error
		strField string
		intField int
	)

	strField, ok = field.(string)
	if !ok {
		log.Logger.Error("type conversion failed")
		return false
	}

	if fieldValidator.FieldType == "string" {
		return validateString(strField, fieldValidator)
	} else if fieldValidator.FieldType == "number" {
		intField, err = strconv.Atoi(strField)
		if err != nil {
			return false
		}
		return validateNumber(intField, fieldValidator)
	} else if fieldValidator.FieldType == "date" {
		return validateDate(strField, fieldValidator)
	} else {
		log.Logger.Errorf("unknown type: %s for field: %s", fieldValidator.FieldType, fieldValidator.FieldName)
		return false
	}
}

func validateString(field string, fieldValidator *FieldValidator) bool {
	var (
		ok             bool
		stringPatterns []*StringPattern
	)

	err := json.Unmarshal([]byte(fieldValidator.Patterns), &stringPatterns)
	if err != nil {
		return false
	}

	for _, stringPattern := range stringPatterns {
		println(stringPattern.Name)
		ok = true
		leftBody := []rune(field)

		for _, pattern := range stringPattern.Patterns {
			// println(pattern.Chars)
			if pattern.Min == 0 {
				pattern.Min = pattern.Max
			}

			// check-size

			if len(leftBody) < pattern.Min {
				ok = false
				break
			}

			asd := int(math.Min(float64(len(leftBody)), float64(pattern.Max)))
			stringToCheck := []byte(string(leftBody[:asd]))
			// println("regexp ", fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), string(leftBody), string(stringToCheck))

			matched, err := regexp.Match(fmt.Sprintf("%s{%d,%d}", pattern.Chars, pattern.Min, pattern.Max), stringToCheck)
			if !matched || err != nil {
				ok = false
			}

			leftBody = leftBody[asd:]
		}

		if len(leftBody) > 0 {
			ok = false
		}

		if ok {
			return true
		}
	}

	return false
}

func validateNumber(field int, fieldValidator *FieldValidator) bool {
	var (
		intPatterns []*IntPattern
	)

	err := json.Unmarshal([]byte(fieldValidator.Patterns), &intPatterns)
	if err != nil {
		return false
	}

	pattern := intPatterns[0]

	return pattern.Min <= field && field <= pattern.Max
}

func validateDate(field string, fieldValidator *FieldValidator) bool {
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
		if maxPattern.PatternType == "date" {
			if !validateMaxDate(fieldDate, maxPattern.Value) {
				return false
			}
		} else if maxPattern.PatternType == "now" {
			if !validateMaxNow(fieldDate, maxPattern.Value) {
				return false
			}
		} else {
			log.Logger.Errorf("unknown date type: %s", maxPattern.PatternType)
		}
	}
	for _, minPattern := range pattern.Min {
		if minPattern.PatternType == "date" {
			if !validateMinDate(fieldDate, minPattern.Value) {
				return false
			}
		} else if minPattern.PatternType == "now" {
			if !validateMinNow(fieldDate, minPattern.Value) {
				return false
			}
		} else if minPattern.PatternType == "depending" {
			if !validateMinDepending(fieldDate, minPattern.Value) {
				return false
			}
		} else {
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
