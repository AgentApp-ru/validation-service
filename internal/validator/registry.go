package validator

import (
	"encoding/json"
	"validation_service/pkg/storage"
)

type registry struct {
	storage storage.Storage
}

var Registry *registry

func Init(store storage.Storage) {
	Registry = &registry{
		storage: store,
	}
}

func (v *registry) GetValidationPattern(object string) ([]byte, error) {
	var (
		rawData []byte
		err     error
	)

	rawData, err = v.storage.Get(object)
	return rawData, err
}

func (v *registry) Get(object string) (interface{}, error) {
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

func (v *registry) GetValidator(data []byte) (*validator, error) {
	vc := new(validator)

	if err := vc.New(data); err != nil {
		return nil, err
	}

	return vc, nil
}
