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
	fieldsMap := make(map[string]interface{})
	// ValidatedObject хранит в себе имя поля и валидность поля
	validationChannel := make(chan validator_module.ValidatedObject)
	validationLength := 0
	for k, v := range data {
		fieldValidator, ok := validatorClass.FieldValidatorsMap[k]
		// если есть валидатор для поля, то запускаем гоурутину
		if !ok {
			fieldsWithErrors = append(fieldsWithErrors, k)
			continue
		}
		go validatorClass.Validate(v, k, fieldValidator, object, fieldsMap, validationChannel)
		// записываем количество валидных полей
		validationLength = validationLength + 1
	}
	// counter для закрытия канала после обработки всех гоурутин
	// как я понял, можно сделать без него, но не понял как
	counter := 0
	for validated := range validationChannel {
		var validatedObject validator_module.ValidatedObject
		validatedObject = validated
		if !validatedObject.Validated {
			fieldsWithErrors = append(fieldsWithErrors, validatedObject.Title)
		}
		counter = counter + 1
		if counter == validationLength {
			close(validationChannel)
		}
	}
	return fieldsWithErrors, nil
}
