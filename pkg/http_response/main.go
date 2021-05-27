package http_response

import (
	"encoding/json"
	"net/http"
)

func HttpError(w http.ResponseWriter, err error) {
	HttpRespond(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
}

func HttpRespond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
