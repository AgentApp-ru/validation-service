package number

import (
	"encoding/json"
	"sync"
	"validation_service/internal/validator/requirements"
	"validation_service/pkg/log"
)

type (
	Pattern struct {
		Min       int       `json:"min"`
		Max       int       `json:"max"`
		Condition Condition `json:"condition"`
	}
	Condition struct {
		Value      DependingValue `json:"value"`
		Expression string         `json:"expression"`
		Statement  []Statement    `json:"statement"`
	}
	DependingValue struct {
		Scope string `json:"scope"`
		Key   string `json:"key"`
	}
	Statement struct {
		Value string `json:"value"`
		Min   int    `json:"min"`
		Max   int    `json:"max"`
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
		log.Logger.Error("Ошибка при парсинге json.")
		return false
	}

	pattern := patterns[0]

	if pattern.Condition.Statement != nil {
		initialValue := requirements.WaitingForValue(pattern.Condition.Value.Key, nv.objectMap)
		if initialValue == nil {
			log.Logger.Error("Ошибка при извлечении значения параметра requirements.requirement.value.condition.value")
			return false
		}
		checkCondition(pattern, initialValue)
	}

	intField = int(floatField)
	if pattern.Min > intField || intField > pattern.Max {
		return false
	}
	return true
}

func checkCondition(pattern *Pattern, initialValue interface{}) {
	switch pattern.Condition.Expression {
	case "starts_with":
		value := initialValue.(string)
		for _, state := range pattern.Condition.Statement {
			if state.Value == value[:len(state.Value)] {
				pattern.Min = state.Min
				pattern.Max = state.Max
			}
		}
	}
}
