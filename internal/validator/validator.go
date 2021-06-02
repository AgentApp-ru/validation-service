package validator

import (
	"encoding/json"
	"validation_service/internal/validator/fields"
	date_validation "validation_service/internal/validator/fields/date"
	num_validation "validation_service/internal/validator/fields/number"
	str_validation "validation_service/internal/validator/fields/string"
	"validation_service/pkg/log"
	"validation_service/pkg/storage"
)

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

type validatorClass struct {
	Schema             string                   `json:"$schema"`
	FieldValidators    []*fields.FieldValidator `json:"validators"`
	FieldValidatorsMap map[string]*fields.FieldValidator
}

func (v *validator) GetValidatorClass(data []byte) *validatorClass {
	vc := &validatorClass{}

	err := json.Unmarshal(data, vc)
	if err != nil {
		log.Logger.Error(err)
	}

	vc.FieldValidatorsMap = map[string]*fields.FieldValidator{}
	for _, field := range vc.FieldValidators {
		vc.FieldValidatorsMap[field.FieldName] = field
	}

	return vc
}

func (vc *validatorClass) Validate(field interface{}, fieldValidator *fields.FieldValidator) bool {
	switch fieldValidator.FieldType {
	case "string":
		return str_validation.Validate(field, fieldValidator)
	case "number":
		return num_validation.Validate(field, fieldValidator)
	case "date":
		return date_validation.Validate(field, fieldValidator)
	default:
		log.Logger.Errorf("unknown type: %s for field: %s", fieldValidator.FieldType, fieldValidator.FieldName)
		return false
	}
}
