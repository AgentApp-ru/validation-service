package validator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"validation_service/pkg/config"
	"validation_service/pkg/storage"
)

type validator struct {
	storage storage.Storage
}

var Validator *validator
var validationsPath string

func Init(store storage.Storage) {
	Validator = &validator{
		storage: store,
	}
	validationsPath = filepath.Join(config.Settings.BasePath, "validations")
}

func (v *validator) GetRaw(object string) ([]byte, error) {
	var (
		rawData []byte
		err     error
	)

	rawData, err = ioutil.ReadFile(filepath.Join(validationsPath, fmt.Sprintf("%s.json", object)))

	// rawData, err = v.storage.Get(object)
	// if err != nil {
	// 	return nil, err
	// }

	return rawData, err
}

func (v *validator) Get(object string) (interface{}, error) {
	var (
		result  interface{}
		rawData []byte
		err     error
	)

	rawData, err = v.GetRaw(object)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawData, &result)

	return result, err
}
