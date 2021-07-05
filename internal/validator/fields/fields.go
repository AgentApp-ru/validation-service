package fields

import "encoding/json"

type FieldValidator struct {
	FieldName    string           `json:"field"`
	FieldType    string           `json:"type"`
	Transformers *json.RawMessage `json:"enabled_transformers"`
	Patterns     json.RawMessage  `json:"patterns"`
}
