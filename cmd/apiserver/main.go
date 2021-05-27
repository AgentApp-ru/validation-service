package main

import (
	"errors"
	"net/http"
	"os"
	"syscall"
	"validation_service/internal/apiserver"
	"validation_service/internal/validator"
	"validation_service/pkg/config"
	"validation_service/pkg/consul"
	"validation_service/pkg/log"
	"validation_service/pkg/shutdown"
)

func main() {
	config.Init()

	log.Init()
	consul.Init()

	validator.Init(consul.Storage)

	server := apiserver.NewServer()

	go shutdown.Graceful([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		&server.HttpServer)

	if err := server.ServeHTTP(); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			log.Logger.Warn("server shutdown")
		default:
			log.Logger.Fatal(err)
		}
	}
}
