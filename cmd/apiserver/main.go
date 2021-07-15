package main

import (
	"errors"
	"net/http"
	"os"
	"syscall"
	"validation_service/internal/apiserver"
	"validation_service/internal/validator"
	"validation_service/pkg/config"
	"validation_service/pkg/log"
	"validation_service/pkg/shutdown"
	"validation_service/pkg/storage/consul"
	"validation_service/pkg/storage/file"
)

func main() {

	config.Init()
	log.Init()

	switch config.Settings.StrorageInfo.Backend {
	case "consul":
		consul.Init()
		validator.Init(consul.Storage)
	case "file":
		file.Init()
		validator.Init(file.Storage)
	default:
		log.Logger.Fatalf("unknown storage backend: %s", config.Settings.StrorageInfo.Backend)
	}

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
