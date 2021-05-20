package validations

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"validation_service/pkg/config"
)

func GetValidation(object string) (interface{}, error) {
	var data interface{}

	content, err := ioutil.ReadFile(filepath.Join(config.Settings.BasePath, "validations", fmt.Sprintf("%s.json", object)))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &data)

	return data, err
}
