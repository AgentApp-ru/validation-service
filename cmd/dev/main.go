package main

import (
	"time"
	"validation_service/pkg/config"
	"validation_service/pkg/log"
)

func main() {
	config.Init()
	log.Init()

	i := 1
	for {
		log.Logger.Infof("test %d", i)
		i++

		time.Sleep(5 * time.Second)
	}

}
