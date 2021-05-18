package apiserver

import (
    "encoding/json"
    "net/http"
    "validation_service/internal/apiserver/views"

    "github.com/gorilla/mux"
)

type server struct {
    router *mux.Router
}

func newServer() *server {
    s := &server{
        router: mux.NewRouter(),
    }

    s.configureRouter()
    return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
    apiRouter := s.router.PathPrefix("/api").Subrouter()
    apiRouter.HandleFunc("/ping", s.handlePing()).Methods("GET")
    v1Router := apiRouter.PathPrefix("/v1").Subrouter()
    v1Router.HandleFunc("/car-validation", s.handleCarValidation()).Methods("GET")
}

func (s *server) handlePing() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        pong := views.Ping()
        s.respond(w, http.StatusOK, pong)
    }
}

func (s *server) handleCarValidation() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        content, err := views.GetCar()
        if err != nil {
            s.error(w, err)
            return
        }
        s.respond(w, http.StatusOK, content)
    }
}

func (s *server) error(w http.ResponseWriter, err error) {
    s.respond(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, code int, data interface{}) {
    w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(code)
    if data != nil {
        json.NewEncoder(w).Encode(data)
    }
}
