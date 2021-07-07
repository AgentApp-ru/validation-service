package fields

import (
	"encoding/json"
	"fmt"
	"sync"
	date_validation "validation_service/internal/validator/fields/date"
	num_validation "validation_service/internal/validator/fields/number"
	str_validation "validation_service/internal/validator/fields/string"
	"validation_service/pkg/log"
)

type (
	FieldValidator struct {
		FieldName        string           `json:"field"`
		FieldType        string           `json:"type"`
		Transformers     *json.RawMessage `json:"enabled_transformers"`
		AllowWhiteSpaces bool             `json:"allow_white_spaces"`
		Patterns         json.RawMessage  `json:"patterns"`
	}

	Validator interface {
		Validate(
			interface{},
			*json.RawMessage,
			json.RawMessage,
			bool,
		) bool
	}
)

func (fv *FieldValidator) ValidateField(field interface{}, object string, allFieldsMap *sync.Map, errors chan string, waiter *sync.WaitGroup) {
	defer waiter.Done()

	objectMapLoaded, _ := allFieldsMap.Load(object)
	objectMap := objectMapLoaded.(*sync.Map)

	var (
		ok        bool
		validator Validator
	)

	switch fv.FieldType {
	case "string":
		validator = str_validation.New(objectMap, allFieldsMap, errors)
	case "number":
		validator = num_validation.New(objectMap, allFieldsMap, errors)
	case "date":
		validator = date_validation.New(objectMap, allFieldsMap, errors)
	default:
		log.Logger.Errorf("unknown type: %s for field: %s", fv.FieldType, fv.FieldName)
		ok = false
	}
	ok = validator.Validate(field, fv.Transformers, fv.Patterns, fv.AllowWhiteSpaces)

	if !ok {
		errors <- fmt.Sprintf("%s/%s: %v", object, fv.FieldName, field)
	} else {
		objectMap.Store(fv.FieldName, field)
	}
}
