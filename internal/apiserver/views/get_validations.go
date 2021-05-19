package views

import (
	"validation_service/pkg/consul"
)

func GetCar() (interface{}, error) {
	return consul.Get("car")
}

func GetInsurerOwner() (interface{}, error) {
	return consul.Get("person")
}

func GetDriver() (interface{}, error) {
	return consul.Get("driver")
}

func GetGeneralConditions() (interface{}, error) {
	return consul.Get("general")
}
