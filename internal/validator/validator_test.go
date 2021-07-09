package validator

import (
	"sync"
	"testing"
	"validation_service/internal/validator/fields"
)

func TestNewValidator(t *testing.T) {
	data := []byte(`
	{
		"validators": [
			{
				"field": "field",
				"type": "type",
				"enabled_transformers": "enabled_transformers",
				"allow_white_spaces": false,
				"patterns": "patterns"
			}
		]
	}
	`)

	v := new(validator)

	err := v.New(data)
	if err != nil {
		t.Error(err.Error())
	}
	if v.fieldValidators == nil {
		t.Errorf("fieldValidators shouldn't be nil after New")
	}
}

func TestInitValidator(t *testing.T) {
	v := &validator{
		fieldValidators: &fieldValidators{
			Validators: []*fields.FieldValidator{},
		},
	}

	errors := make(chan string)
	v.Init("object", &sync.Map{}, errors)
	if v.object != "object" {
		t.Errorf("expected: 'object' actual '%s'", v.object)
	}
	if v.errors != errors {
		t.Errorf("errors chan should be equal to received")
	}
}
