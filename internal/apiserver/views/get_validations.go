package views

import (
	"validation_service/internal/validations"
)

func GetCar() (interface{}, error) {
	return validations.GetValidation("car")
}