package main

import (
	"log"
	"validation_service/pkg/apiserver"
)

func main() {
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
