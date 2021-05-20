package views

import (
	"validation_service/internal/validations"
)

func GetCar() (interface{}, error) {
	return validations.GetValidation("car")
}

func GetInsurerOwner() (interface{}, error) {
	return validations.GetValidation("person")
}

func GetDriver() (interface{}, error) {
	return validations.GetValidation("driver")
}
