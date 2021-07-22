package validator

import (
	"encoding/json"
	"fmt"
	"sync"
	"validation_service/internal/validator/fields"
)

type (
	fieldValidators struct {
		Validators []*fields.FieldValidator `json:"validators"`
	}

	validator struct {
		fieldValidators *fieldValidators
		fields          map[string]*fields.FieldValidator
		errors          chan string
		absentFields    chan string
		object          string
	}
)

func (v *validator) New(data []byte) error {
	var fieldValidators *fieldValidators

	err := json.Unmarshal(data, &fieldValidators)
	if err != nil {
		return err
	}

	v.fieldValidators = fieldValidators

	return nil
}

func (v *validator) Init(object string, agreementFields *sync.Map, errors, absentFields chan string) {
	v.errors = errors
	v.absentFields = absentFields
	v.object = object

	v.fields = make(map[string]*fields.FieldValidator)
	for _, field := range v.fieldValidators.Validators {
		field.Init(object, agreementFields, errors, absentFields)
		v.fields[field.FieldName] = field
	}
}

func (v *validator) Validate(data map[string]interface{}) {
	waiter := new(sync.WaitGroup)

	for fieldName, fieldValidator := range v.fields {
		waiter.Add(1)
		go fieldValidator.CheckRequirementField(data[fieldName], waiter)
	}

	waiter.Wait()

	for fieldName, fieldContent := range data {
		fieldValidator, ok := v.fields[fieldName]
		if !ok {
			v.errors <- fmt.Sprintf("%s/%s: %v", v.object, fieldName, fieldContent)
			continue
		}

		waiter.Add(1)
		go fieldValidator.ValidateField(fieldContent, waiter)
	}

	waiter.Wait()
}
