package date

import (
	"encoding/json"
	"sync"
	"time"
	"validation_service/pkg/log"
)

type (
	DateDateValue  string
	DependingValue struct {
		Scope string `json:"scope"`
		Key   string `json:"key"`
	}

	DateValidator struct {
		fieldName        string
		objectMap        *sync.Map
		allFieldsMap     *sync.Map
		errors           chan string
		transformers     *json.RawMessage
		patterns         json.RawMessage
		allowWhiteSpaces bool
	}

	DatePattern struct {
		Min []*DateMinMaxPattern `json:"min"`
		Max []*DateMinMaxPattern `json:"max"`
	}

	DateMinMaxPattern struct {
		PatternType string          `json:"type"`
		Value       json.RawMessage `json:"value"`
	}
)

func New() *DateValidator {
	return new(DateValidator)
}

func (dv *DateValidator) Init(
	fieldName string,
	objectMap,
	allFieldsMap *sync.Map,
	errors chan string,
	transformers *json.RawMessage,
	patterns json.RawMessage,
	allowWhiteSpaces bool,
	_ int,
	_ int,
) {
	dv.fieldName = fieldName
	dv.objectMap = objectMap
	dv.allFieldsMap = allFieldsMap
	dv.errors = errors
	dv.transformers = transformers
	dv.patterns = patterns
	dv.allowWhiteSpaces = allowWhiteSpaces
}

func (dv *DateValidator) Validate(field interface{}) bool {
	var (
		fieldDate    time.Time
		datePatterns []*DatePattern
		err          error
		ok           bool
		strField     string
	)

	if strField, ok = field.(string); !ok {
		log.Logger.Errorf("type conversion failed: %s -> %v", dv.fieldName, field)
		return false
	}

	if err = json.Unmarshal([]byte(dv.patterns), &datePatterns); err != nil {
		// TODO: log error
		return false
	}

	if fieldDate, err = time.Parse("2006-01-02", strField); err != nil {
		// TODO: log error
		return false
	}

	pattern := datePatterns[0]
	if !dv.validateMinDate(fieldDate, pattern) {
		return false
	}
	if !dv.validateMaxDate(fieldDate, pattern) {
		return false
	}

	return true
}
