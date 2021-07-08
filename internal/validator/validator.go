package validator

import (
	"encoding/json"
	"fmt"
	"sync"
	"validation_service/internal/validator/fields"
	"validation_service/pkg/storage"
)

type ValidatedObject struct {
	Title     string
	Validated bool
}

type validator struct {
	storage storage.Storage
}

var Validator *validator

func Init(store storage.Storage) {
	Validator = &validator{
		storage: store,
	}
}

func (v *validator) GetRaw(object string) ([]byte, error) {
	var (
		rawData []byte
		err     error
	)

	rawData, err = v.storage.Get(object)
	return rawData, err
}

func (v *validator) Get(object string) (interface{}, error) {
	var (
		result  interface{}
		rawData []byte
		err     error
	)

	rawData, err = v.storage.Get(object)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawData, &result)

	return result, err
}

func (v *validator) GetValidatorClass(data []byte) (*validatorClass, error) {
	vc := &validatorClass{}

	err := json.Unmarshal(data, vc)
	if err != nil {
		return nil, err
	}

	vc.FieldValidatorsMap = make(map[string]*fields.FieldValidator)
	for _, field := range vc.FieldValidators {
		vc.FieldValidatorsMap[field.FieldName] = field
	}

	return vc, nil
}

type validatorClass struct {
	Schema             string                   `json:"$schema"`
	FieldValidators    []*fields.FieldValidator `json:"validators"`
	FieldValidatorsMap map[string]*fields.FieldValidator
}

// func (vc *validatorClass) Validate(field interface{}, fieldTitle string, fieldValidator *fields.FieldValidator, object string, fieldsMap map[string]interface{}, validationChannel chan ValidatedObject, lock *sync.Mutex) {
func (vc *validatorClass) Validate(object string, fields *sync.Map, data map[string]interface{}, errors chan string) {
	waiter := new(sync.WaitGroup)

	for k, v := range data {
		println("validate", k)
		fieldValidator, ok := vc.FieldValidatorsMap[k]
		if !ok {
			errors <- fmt.Sprintf("%s/%s: %v", object, k, v)
			continue
		}
		waiter.Add(1)
		go fieldValidator.ValidateField(v, object, fields, errors, waiter)
	}

	waiter.Wait()
}
