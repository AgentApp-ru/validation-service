package views

import (
	"validation_service/internal/validator"
)

func GetValidationPattern(object string) (interface{}, error) {
	return validator.Registry.Get(object)
}
