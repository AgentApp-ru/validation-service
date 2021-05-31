package number

import (
	"encoding/json"
	"strconv"
	"validation_service/internal/validator/fields"
)

type IntPattern struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

func Validate(field string, fieldValidator *fields.FieldValidator) bool {
	var (
		intField    int
		err         error
		intPatterns []*IntPattern
	)

	intField, err = strconv.Atoi(field)
	if err != nil {
		return false
	}

	err = json.Unmarshal([]byte(fieldValidator.Patterns), &intPatterns)
	if err != nil {
		return false
	}

	pattern := intPatterns[0]

	return pattern.Min <= intField && intField <= pattern.Max
}
