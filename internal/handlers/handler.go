package handlers

import (
	"github.com/gorilla/mux"
)

type Handler struct {
}

func Register(router *mux.Router) {
	// router.(mux.Router).HandleFunc("/validations/car", handleCarValidation()).Methods("GET")
}
