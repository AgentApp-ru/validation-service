package main

import (
	"log"
	"validation_service/internal/apiserver"
)

func main() {
	if err := apiserver.Start(); err != nil {
		log.Fatal(err)
	}
}
