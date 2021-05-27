package views

import (
	"encoding/json"
	"regexp"
	validator_module "validation_service/internal/validator"
	"validation_service/pkg/log"
)

type Pattern struct {
	Chars string `json:"chars"`
	Min   string `json:"min"`
	min   int
	Max   int `json:"max"`
}

type StringPattern struct {
	Name               string     `json:"name"`
	Allow_white_spaces bool       `json:"allow_white_spaces"`
	Patterns           []*Pattern `json:"patterns"`
}

type FieldValidator struct {
	FieldName      string           `json:"field"`
	FieldType      string           `json:"type"`
	StringPatterns []*StringPattern `json:"patterns"`
}

type validatorClass struct {
	Schema             string            `json:"$schema"`
	FieldValidators    []*FieldValidator `json:"validators"`
	FieldValidatorsMap map[string]*FieldValidator
}

func ValidateCar(data map[string]interface{}) bool {
	rawValidator, err := validator_module.Validator.GetRaw("car")
	if err != nil {
		return false
	}

	validator := getValidator(rawValidator)
	if validator == nil {
		return false
	}

	for k, v := range data {
		log.Logger.Infof("validating %s: %s", k, v)
		fieldValidator, ok := validator.FieldValidatorsMap[k]
		if !ok {
			log.Logger.Info("TODO: что-то сделать")
		}

		ok = validate(v, fieldValidator)
		log.Logger.Info(ok)
	}

	return true
}

func getValidator(data []byte) *validatorClass {
	v := &validatorClass{}

	err := json.Unmarshal(data, v)
	if err != nil {
		log.Logger.Error(err)
	}

	v.FieldValidatorsMap = map[string]*FieldValidator{}
	for _, field := range v.FieldValidators {
		println("field ", field.FieldName)
		v.FieldValidatorsMap[field.FieldName] = field
	}

	return v
}

func validate(field interface{}, fieldValidator *FieldValidator) bool {
	var ok bool

	log.Logger.Info(field)

	if fieldValidator.FieldType == "string" {
		field, ok = field.(string)
	} else if fieldValidator.FieldType == "number" {
		field, ok = field.(int)
		// } else if fieldValidator.fieldType == "date" {
		// 	field, ok := field.(string)
	} else {
		log.Logger.Errorf("unknown type: %s for field: %s", fieldValidator.FieldType, fieldValidator.FieldName)
		return false
	}

	if !ok {
		log.Logger.Info("type conversion failed")
		return false
	}

	leftBody := field.(string)
	for _, pattern := range fieldValidator.StringPatterns[0].Patterns {
		println(leftBody, pattern.Chars)

		if pattern.Min == "" {
			pattern.min = pattern.Max
		}

		matched, err := regexp.Match(pattern.Chars, []byte(leftBody[:pattern.Max]))
		println(matched, err)

		leftBody = leftBody[pattern.Max:]
	}

	return true
}
