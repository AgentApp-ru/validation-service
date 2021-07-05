package validator

import (
	"sync"
	"testing"
	"validation_service/pkg/config"
	"validation_service/pkg/log"
	"validation_service/pkg/storage/file"
)

func TestCarValidVinField(t *testing.T) {
	config.Init()
	log.Init()
	file.Init()
	Init(file.Storage)

	rawValidator, _ := Validator.GetRaw("car")
	validatorClass := Validator.GetValidatorClass(rawValidator)

	data := map[string]string{
		"number_plate":       "Е2",
		"vin":                "TMBED45J2B3209311",
		"manufacturing_year": "1929-01-01",
	}

	fieldsWithErrors := []string{}

	lock := sync.Mutex{}
	fieldsMap := make(map[string]interface{})
	validationChannel := make(chan ValidatedObject)
	for k, v := range data {
		fieldValidator, ok := validatorClass.FieldValidatorsMap[k]
		if !ok {
			fieldsWithErrors = append(fieldsWithErrors, k)
			continue
		}
		go validatorClass.Validate(v, k, fieldValidator, "car", fieldsMap, validationChannel, &lock)

	}
	counter := 0
	length := len(data)
	for validated := range validationChannel {
		var validatedObject ValidatedObject
		validatedObject = validated
		if !validatedObject.Validated {
			fieldsWithErrors = append(fieldsWithErrors, validatedObject.Title)
		}
		counter = counter + 1
		if counter == length {
			close(validationChannel)
		}
	}

	if len(fieldsWithErrors) != 2 {
		t.Errorf("fields with errors should be 2. And they are: %v", fieldsWithErrors)
	}
}

func TestCarValidFields(t *testing.T) {
	config.Init()
	log.Init()
	file.Init()
	Init(file.Storage)

	rawValidator, _ := Validator.GetRaw("car")
	validatorClass := Validator.GetValidatorClass(rawValidator)

	data := map[string]string{
		"number_plate":       "Е271ХМ178",
		"vin":                "TMBED45J2B3209311",
		"manufacturing_year": "1931-01-01",
	}

	fieldsWithErrors := []string{}

	lock := sync.Mutex{}
	fieldsMap := make(map[string]interface{})
	validationChannel := make(chan ValidatedObject)
	for k, v := range data {
		fieldValidator, ok := validatorClass.FieldValidatorsMap[k]
		if !ok {
			fieldsWithErrors = append(fieldsWithErrors, k)
			continue
		}
		go validatorClass.Validate(v, k, fieldValidator, "car", fieldsMap, validationChannel, &lock)
	}
	counter := 0
	length := len(data)
	for validated := range validationChannel {
		var validatedObject ValidatedObject
		validatedObject = validated
		if !validatedObject.Validated {
			fieldsWithErrors = append(fieldsWithErrors, validatedObject.Title)
		}
		counter = counter + 1
		if counter == length {
			close(validationChannel)
		}
	}

	if len(fieldsWithErrors) > 0 {
		t.Errorf("fields with errors should be 2. And they are: %v", fieldsWithErrors)
	}
}
