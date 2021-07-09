package validator

import (
	"testing"
)

type TestStorage struct {
}

func (ts *TestStorage) Get(a string) ([]byte, error) {
	return []byte("test bytes"), nil
}

func TestGetValidationPatternFromRegistry(t *testing.T) {
	r := &registry{
		storage: &TestStorage{},
	}

	result, _ := r.GetValidationPattern("test")
	if string(result) != "test bytes" {
		t.Errorf("expected: 'test bytes' actual '%s'", string(result))
	}
}

func TestGetValidator(t *testing.T) {
	r := &registry{
		storage: &TestStorage{},
	}

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

	validator, err := r.GetValidator(data)
	if err != nil {
		t.Error(err.Error())
	}
	if validator.fieldValidators == nil {
		t.Errorf("fieldValidators shouldn't be nil after GetValidator")
	}
}
