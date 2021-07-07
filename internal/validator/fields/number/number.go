package number

import (
	"encoding/json"
	"sync"
	"validation_service/pkg/log"
)

type (
	IntPattern struct {
		Min int `json:"min"`
		Max int `json:"max"`
	}
	Validator struct {
		objectMap    *sync.Map
		allFieldsMap *sync.Map
		errors       chan string
	}
)

func New(objectMap, allFieldsMap *sync.Map, errors chan string) *Validator {
	return &Validator{
		objectMap:    objectMap,
		allFieldsMap: allFieldsMap,
		errors:       errors,
	}
}

func (v *Validator) Validate(
	field interface{},
	transformers *json.RawMessage,
	patterns json.RawMessage,
	allowWhiteSpaces bool,
) bool {
	var (
		floatField  float64
		intField    int
		ok          bool
		intPatterns []*IntPattern
	)

	if floatField, ok = field.(float64); !ok {
		log.Logger.Error("type conversion failed")
		return false
	}

	if err := json.Unmarshal([]byte(patterns), &intPatterns); err != nil {
		log.Logger.Error("json parsing error")
		return false
	}

	pattern := intPatterns[0]

	intField = int(floatField)
	if pattern.Min > intField || intField > pattern.Max {
		return false
	}
	return true
}
