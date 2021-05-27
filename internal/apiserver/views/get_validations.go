package views

import (
	"validation_service/internal/validator"
)

func GetCar() (interface{}, error) {
	return validator.Validator.Get("car")
}

func GetInsurerOwner() (interface{}, error) {
	return validator.Validator.Get("person")
}

func GetDriver() (interface{}, error) {
	return validator.Validator.Get("driver")
}

func GetGeneralConditions() (interface{}, error) {
	return validator.Validator.Get("general")
}