package validations

import (
    "encoding/json"
    "io/ioutil"
    "log"
    "path/filepath"
    "validation_service/pkg/config"
)

var CarValidationJSON []byte

func init() {
    var err error

    CarValidationJSON, err = ioutil.ReadFile(filepath.Join(config.Settings.BasePath, "validations", "car.json"))
    if err != nil {
        log.Fatal(err)
    }
}

func GetValidation(object string) (interface{}, error) {
    var data interface{}

    err := json.Unmarshal(CarValidationJSON, &data)

    return data, err
}
