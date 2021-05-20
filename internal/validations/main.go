package validations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"validation_service/pkg/config"
)

var validationsPath string

func init() {
	validationsPath = filepath.Join(config.Settings.BasePath, "validations")
}

func GetValidation(object string) (interface{}, error) {
	var data interface{}

	content, err := ioutil.ReadFile(filepath.Join(validationsPath, fmt.Sprintf("%s.json", object)))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &data)

	return data, err
}
