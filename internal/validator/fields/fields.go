package fields

import "encoding/json"

type FieldValidator struct {
	FieldName          string           `json:"field"`
	FieldType          string           `json:"type"`
	Transformers       *json.RawMessage `json:"enabled_transformers"`
	Allow_white_spaces bool             `json:"allow_white_spaces"`
	Patterns           json.RawMessage  `json:"patterns"`
}
