package fields

import "encoding/json"


type FieldValidator struct {
	FieldName string          `json:"field"`
	FieldType string          `json:"type"`
	Patterns  json.RawMessage `json:"patterns"`
}
