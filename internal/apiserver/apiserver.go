package apiserver

import (
    "net/http"
    "validation_service/pkg/config"
)

func Start() error {
    return http.ListenAndServe(config.Settings.BindAddr, newServer())
}
