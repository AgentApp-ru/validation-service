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
	NumberValidator struct {
		fieldName        string
		objectMap        *sync.Map
		allFieldsMap     *sync.Map
		errors           chan string
		transformers     *json.RawMessage
		patterns         json.RawMessage
		allowWhiteSpaces bool
	}
)

func New() *NumberValidator {
	return new(NumberValidator)
}

func (nv *NumberValidator) Init(
	fieldName        string,
	objectMap,
	allFieldsMap *sync.Map,
	errors chan string,
	transformers *json.RawMessage,
	patterns json.RawMessage,
	allowWhiteSpaces bool,
) {
	nv.fieldName = fieldName
	nv.objectMap = objectMap
	nv.allFieldsMap = allFieldsMap
	nv.errors = errors
	nv.transformers = transformers
	nv.patterns = patterns
	nv.allowWhiteSpaces = allowWhiteSpaces
}

func (nv *NumberValidator) Validate(field interface{}) bool {
	var (
		floatField  float64
		intField    int
		ok          bool
		intPatterns []*IntPattern
	)

	if floatField, ok = field.(float64); !ok {
		log.Logger.Errorf("type conversion failed: %s -> %v", nv.fieldName, field)
		return false
	}

	if err := json.Unmarshal([]byte(nv.patterns), &intPatterns); err != nil {
		log.Logger.Error("Ошибка при ")
		return false
	}

	pattern := intPatterns[0]

	intField = int(floatField)
	if pattern.Min > intField || intField > pattern.Max {
		return false
	}
	return true
}
