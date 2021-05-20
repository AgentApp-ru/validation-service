package apiserver

import (
	"fmt"
	"net/http"
	"time"
	"validation_service/internal/apiserver/views"
	"validation_service/pkg/config"
	"validation_service/pkg/http_response"

	"github.com/gorilla/mux"
)

type server struct {
	HttpServer http.Server
}

func NewServer() *server {
	r := mux.NewRouter()
	configureRouter(r)

	httpServer := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf("127.0.0.1%s", config.Settings.BindAddr),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	s := &server{
		HttpServer: *httpServer,
	}

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.HttpServer.Handler.ServeHTTP(w, r)
}

func configureRouter(router *mux.Router) {
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/ping", handlePing()).Methods("GET")
	v1Router := apiRouter.PathPrefix("/v1").Subrouter()
	v1Router.HandleFunc("/validations/car", handleCarValidation()).Methods("GET")
	v1Router.HandleFunc("/validations/person", handlePersonValidation()).Methods("GET")
	v1Router.HandleFunc("/validations/driver", handleDriverValidation()).Methods("GET")
	v1Router.HandleFunc("/validations/general-conditions", handleGeneralConditionsValidation()).Methods("GET")
}

func handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pong := views.Ping()
		http_response.HttpRespond(w, http.StatusOK, pong)
	}
}

func handleCarValidation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := views.GetCar()
		if err != nil {
			http_response.HttpError(w, err)
			return
		}
		http_response.HttpRespond(w, http.StatusOK, content)
	}
}

func handlePersonValidation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := views.GetInsurerOwner()
		if err != nil {
			http_response.HttpError(w, err)
			return
		}
		http_response.HttpRespond(w, http.StatusOK, content)
	}
}

func handleDriverValidation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := views.GetDriver()
		if err != nil {
			http_response.HttpError(w, err)
			return
		}
		http_response.HttpRespond(w, http.StatusOK, content)
	}
}

func handleGeneralConditionsValidation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := views.GetGeneralConditions()
		if err != nil {
			http_response.HttpError(w, err)
			return
		}
		http_response.HttpRespond(w, http.StatusOK, content)
	}
}
