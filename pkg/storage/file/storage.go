package file

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"validation_service/pkg/config"
)

type fileStorage struct {
	validationsPath string
}

var Storage *fileStorage

func Init() {
	Storage = &fileStorage{
		validationsPath: filepath.Join(config.Settings.BasePath, "validations"),
	}
}

func (s *fileStorage) Get(object string) ([]byte, error) {
	rawData, err := ioutil.ReadFile(filepath.Join(s.validationsPath, fmt.Sprintf("%s.json", object)))
	return rawData, err
}
