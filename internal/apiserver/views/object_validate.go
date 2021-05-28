package views

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	validator_module "validation_service/internal/validator"
	"validation_service/pkg/log"
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

func ValidateCar(data map[string]interface{}) (bool, []string) {
	fieldsWithErrors := []string{}

	// return true, fieldsWithErrors

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
		strField string
		intField int
	)

	strField, ok = field.(string)
	if fieldValidator.FieldType == "string" {
		return validateString(strField, fieldValidator)
	} else if fieldValidator.FieldType == "number" {
		// println("raw", strField)
		intField, _ = strconv.Atoi(strField)
		// println("int", intField)
		// println("convert to int", ok)
		return validateNumber(intField, fieldValidator)
		// } else if fieldValidator.fieldType == "date" {
		// 	field, ok := field.(string)
	} else {
		log.Logger.Errorf("unknown type: %s for field: %s", fieldValidator.FieldType, fieldValidator.FieldName)
		return false
	}

	if !ok {
		log.Logger.Error("type conversion failed")
		return false
	}

	return false
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
		// println(err)
		return false
	}

	pattern := intPatterns[0]

	// println("min", pattern.Min)
	// println("field", field)
	// println("max", pattern.Max)
	// println("1", pattern.Min < field)
	// println("2", field < pattern.Max)

	return pattern.Min <= field && field <= pattern.Max
}
