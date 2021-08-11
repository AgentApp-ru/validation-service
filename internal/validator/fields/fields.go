package fields

import (
	"encoding/json"
	"fmt"
	"sync"
	date_validation "validation_service/internal/validator/fields/date"
	num_validation "validation_service/internal/validator/fields/number"
	str_validation "validation_service/internal/validator/fields/string"
	"validation_service/internal/validator/requirements"
	"validation_service/pkg/log"
)

type (
	FieldValidator struct {
		FieldName          string                    `json:"field"`
		FieldType          string                    `json:"type"`
		Requirements       requirements.Requirements `json:"requirements"`
		Transformers       *json.RawMessage          `json:"enabled_transformers"`
		AllowWhiteSpaces   bool                      `json:"allow_white_spaces"`
		Patterns           json.RawMessage           `json:"patterns"`
		object             string
		objectMap          *sync.Map
		agreementFieldsMap *sync.Map
		errors             chan string
		absentFields       chan string
	}

	fieldvalidatorImpl interface {
		Init(string, *sync.Map, *sync.Map, chan string, *json.RawMessage, json.RawMessage, bool)
		Validate(interface{}) bool
	}
)

func (fv *FieldValidator) Init(object string, agreementFields *sync.Map, errors, absentFields chan string) {
	fv.object = object
	fv.agreementFieldsMap = agreementFields
	fv.errors = errors
	fv.absentFields = absentFields

	objectMapLoaded, _ := fv.agreementFieldsMap.Load(fv.object)
	fv.objectMap = objectMapLoaded.(*sync.Map)
}

func (fv *FieldValidator) CheckRequirementField(field interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()

	isfieldPresent := field != nil
	if isfieldPresent {
		fv.objectMap.Store(fv.FieldName, field)
	}

	if !fv.Requirements.CheckRequired(isfieldPresent, fv.objectMap) {
		fv.absentFields <- fmt.Sprintf("%s/%s", fv.object, fv.FieldName)
		return
	}
}

func (fv *FieldValidator) ValidateField(field interface{}, waiter *sync.WaitGroup) {
	defer waiter.Done()

	var (
		ok                 bool
		fieldvalidatorImpl fieldvalidatorImpl
	)

	switch fv.FieldType {
	case "string":
		fieldvalidatorImpl = str_validation.New()
	case "number":
		fieldvalidatorImpl = num_validation.New()
	case "date":
		fieldvalidatorImpl = date_validation.New()
	default:
		log.Logger.Errorf("unknown type: %s for field: %s", fv.FieldType, fv.FieldName)
		fv.errors <- fmt.Sprintf("%s/%s: %v", fv.object, fv.FieldName, field)
		return
	}

	fieldvalidatorImpl.Init(
		fmt.Sprintf("%s/%s", fv.object, fv.FieldName),
		fv.objectMap,
		fv.agreementFieldsMap,
		fv.errors,
		fv.Transformers,
		fv.Patterns,
		fv.AllowWhiteSpaces,
	)
	ok = fieldvalidatorImpl.Validate(field)

	if !ok {
		fv.errors <- fmt.Sprintf("%s/%s: %v", fv.object, fv.FieldName, field)
	}
}
