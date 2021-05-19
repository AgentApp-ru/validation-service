package validator

import (
	"encoding/json"
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
