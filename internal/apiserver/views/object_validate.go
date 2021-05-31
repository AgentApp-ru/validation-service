package views

import (
	"encoding/json"
	validator_module "validation_service/internal/validator"
	"validation_service/internal/validator/fields"
	date_validation "validation_service/internal/validator/fields/date"
	num_validation "validation_service/internal/validator/fields/number"
	str_validation "validation_service/internal/validator/fields/string"
	"validation_service/pkg/log"
)

type validatorClass struct {
	Schema             string                   `json:"$schema"`
	FieldValidators    []*fields.FieldValidator `json:"validators"`
	FieldValidatorsMap map[string]*fields.FieldValidator
}

func ValidateCar(data map[string]interface{}) (bool, []string) {
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

	v.FieldValidatorsMap = map[string]*fields.FieldValidator{}
	for _, field := range v.FieldValidators {
		v.FieldValidatorsMap[field.FieldName] = field
	}

	return v
}

func validate(field interface{}, fieldValidator *fields.FieldValidator) bool {
	var (
		ok       bool
		strField string
	)

	strField, ok = field.(string)
	if !ok {
		log.Logger.Error("type conversion failed")
		return false
	}

	switch fieldValidator.FieldType {
	case "string":
		return str_validation.Validate(strField, fieldValidator)
	case "number":
		return num_validation.Validate(strField, fieldValidator)
	case "date":
		return date_validation.Validate(strField, fieldValidator)
	default:
		log.Logger.Errorf("unknown type: %s for field: %s", fieldValidator.FieldType, fieldValidator.FieldName)
		return false
	}
}
