package main

import (
	"validation_service/pkg/config"
	"validation_service/pkg/consul"
)

func main() {
	config.Init()
	config.Get()

	consul.Init()
	consul.Get("car")
}
