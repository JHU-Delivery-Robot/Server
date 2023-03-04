package restserver

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func Error(w http.ResponseWriter, httpCode int) {
	http.Error(w, http.StatusText(httpCode), httpCode)
}

func Respond(w http.ResponseWriter, data interface{}, httpCode int) {
	var resp interface{}
	if v, ok := data.(error); ok {
		resp = errorResponse{Error: v.Error()}
	} else {
		resp = data
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	if resp != nil {
		_ = json.NewEncoder(w).Encode(resp)
	}
}
