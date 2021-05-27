package main

import (
	"validation_service/internal/validator"
	"validation_service/pkg/config"
	"validation_service/pkg/consul"
)

func main() {
	config.Init()

	consul.Init()

	validator.Init(consul.Storage)
}
