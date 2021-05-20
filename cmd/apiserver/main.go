package main

import (
	"errors"
	"net/http"
	"os"
	"syscall"
	"validation_service/internal/apiserver"
	"validation_service/internal/validations"
	"validation_service/pkg/config"
	"validation_service/pkg/log"
	"validation_service/pkg/shutdown"
)

func main() {
	config.Init()
	log.Init(config.Settings)

	validations.Init()

	server := apiserver.NewServer()

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		&server.HttpServer)

	if err := server.HttpServer.ListenAndServe(); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Logger.Warn("server shutdown")
		default:
			log.Logger.Fatal(err)
		}
	}
}
