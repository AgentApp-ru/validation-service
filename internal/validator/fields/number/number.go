package number

import (
	"encoding/json"
	"sync"
	"validation_service/pkg/log"
)

type (
	Pattern struct {
		Min int `json:"min"`
		Max int `json:"max"`
	}
	Validator struct {
		fieldName        string
		objectMap        *sync.Map
		allFieldsMap     *sync.Map
		errors           chan string
		transformers     *json.RawMessage
		patterns         json.RawMessage
		allowWhiteSpaces bool
	}
)

func New() *Validator {
	return new(Validator)
}

func (nv *Validator) Init(
	fieldName string,
	objectMap,
	allFieldsMap *sync.Map,
	errors chan string,
	transformers *json.RawMessage,
	patterns json.RawMessage,
	allowWhiteSpaces bool,
	_ int,
) {
	nv.fieldName = fieldName
	nv.objectMap = objectMap
	nv.allFieldsMap = allFieldsMap
	nv.errors = errors
	nv.transformers = transformers
	nv.patterns = patterns
	nv.allowWhiteSpaces = allowWhiteSpaces
}

func (nv *Validator) Validate(field interface{}) bool {
	var (
		floatField float64
		intField   int
		ok         bool
		patterns   []*Pattern
	)

	if floatField, ok = field.(float64); !ok {
		log.Logger.Errorf("type conversion failed: %s -> %v", nv.fieldName, field)
		return false
	}

	if err := json.Unmarshal([]byte(nv.patterns), &patterns); err != nil {
		log.Logger.Error("Ошибка при ")
		return false
	}

	pattern := patterns[0]

	intField = int(floatField)
	if pattern.Min > intField || intField > pattern.Max {
		return false
	}
	return true
}
