package number

import (
	"encoding/json"
	"validation_service/internal/validator/fields"
	"validation_service/pkg/log"
)

type IntPattern struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func Validate(field interface{}, fieldValidator *fields.FieldValidator) interface{} {
	var (
		floatField  float64
		intField    int
		ok          bool
		intPatterns []*IntPattern
	)

	if floatField, ok = field.(float64); !ok {
		log.Logger.Error("type conversion failed")
		return nil
	}

	if err := json.Unmarshal([]byte(fieldValidator.Patterns), &intPatterns); err != nil {
		log.Logger.Error("json parsing error")
		return nil
	}

	pattern := intPatterns[0]

	intField = int(floatField)
	if pattern.Min <= intField && intField <= pattern.Max {
		return intField
	}
	return nil
}
