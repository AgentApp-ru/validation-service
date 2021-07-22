package apiserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"validation_service/internal/apiserver/views"
	"validation_service/pkg/config"
	"validation_service/pkg/http_response"
	"validation_service/pkg/log"

	"github.com/gorilla/mux"
)

type server struct {
	HttpServer http.Server
}

func NewServer() *server {
	r := mux.NewRouter()
	configureRouter(r)

	return &server{
		HttpServer: http.Server{
			Handler:      r,
			Addr:         fmt.Sprintf("0.0.0.0%s", config.Settings.BindAddr),
			WriteTimeout: 30 * time.Second,
			ReadTimeout:  30 * time.Second,
		},
	}
}

func (s *server) ServeHTTP() error {
	return s.HttpServer.ListenAndServe()
}

func LogRequestMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		log.Logger.Infof("%s %s", r.Method, r.URL)
	})
}

func configureRouter(router *mux.Router) {
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.Use(LogRequestMiddleware)
	apiRouter.HandleFunc("/ping", handlePing()).Methods("GET")

	v1Router := apiRouter.PathPrefix("/v1").Subrouter()
	v1Router.HandleFunc("/validations/car", handleGetValidationPattern("car")).Methods("GET")
	v1Router.HandleFunc("/validations/person", handleGetValidationPattern("person")).Methods("GET")
	v1Router.HandleFunc("/validations/driver", handleGetValidationPattern("driver")).Methods("GET")
	v1Router.HandleFunc("/validations/agreement", handleGetValidationPattern("agreement")).Methods("GET")

	v1Router.HandleFunc("/validations", handleValidateAll()).Methods("POST")
	// v1Router.HandleFunc("/validations/car", handleValidate("car")).Methods("POST")
	// v1Router.HandleFunc("/validations/person", handleValidate("person")).Methods("POST")
	// v1Router.HandleFunc("/validations/driver", handleValidate("driver")).Methods("POST")
}

func handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pong := views.Ping()
		http_response.HttpRespond(w, http.StatusOK, pong)
	}
}

func handleGetValidationPattern(object string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		content, err := views.GetValidationPattern(object)
		if err != nil {
			http_response.HttpError(w, err)
			return
		}
		http_response.HttpRespond(w, http.StatusOK, content)
	}
}

func handleValidateAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var (
			bodyRaw []byte
			err     error
		)
		defer r.Body.Close()

		bodyRaw, err = ioutil.ReadAll(r.Body)
		if err != nil {
			http_response.HttpError(w, fmt.Errorf("error read body: %s", err))
			return
		}

		absentFields, errorFields, err := views.ValidateAgreement(bodyRaw)
		if err != nil {
			http_response.HttpError(w, err)
			log.Logger.Infof("error duration %v", time.Since(start).Round(time.Second))
			return
		}

		response := map[string][]string{
			"absent fields": absentFields,
			"error fields": errorFields,
		}
		http_response.HttpRespond(w, http.StatusOK, response)
		log.Logger.Infof("duration %v", time.Since(start).Round(time.Second))
	}
}

// func handleValidate(object string) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var (
// 			body []byte
// 			err  error
// 		)
// 		defer r.Body.Close()

// 		body, err = ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			http_response.HttpError(w, fmt.Errorf("error read body: %s", err))
// 			return
// 		}

// 		var b map[string]interface{}
// 		if err = json.Unmarshal(body, &b); err != nil {
// 			http_response.HttpError(w, fmt.Errorf("error convert to json body: %s", err))
// 			return
// 		}

// 		fieldsWithErrors, err := views.Validate(object, b)
// 		if err != nil {
// 			http_response.HttpError(w, fmt.Errorf("errors while validating: %s", err))
// 			return
// 		}
// 		if len(fieldsWithErrors) != 0 {
// 			errors := map[string][]string{
// 				"error fields": fieldsWithErrors,
// 			}
// 			http_response.HttpRespond(w, http.StatusOK, errors)
// 			return
// 		}

// 		http_response.HttpRespond(w, http.StatusOK, nil)
// 	}
// }
