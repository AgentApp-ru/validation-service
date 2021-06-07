package views

import (
	validator_module "validation_service/internal/validator"
)

func Validate(object string, data map[string]interface{}) ([]string, error) {
	fieldsWithErrors := []string{}

	rawValidator, err := validator_module.Validator.GetRaw(object)
	if err != nil {
		return []string{}, err
	}

	validatorClass := validator_module.Validator.GetValidatorClass(rawValidator)
	if validatorClass == nil {
		return []string{}, err
	}

	for k, v := range data {
		fieldValidator, ok := validatorClass.FieldValidatorsMap[k]
		if !ok || !validatorClass.Validate(v, fieldValidator, object) {
			fieldsWithErrors = append(fieldsWithErrors, k)
		}
	}

	return fieldsWithErrors, nil
}
