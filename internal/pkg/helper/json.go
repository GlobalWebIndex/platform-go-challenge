package helper

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error   string      `json:"error"`
	Details interface{} `json:"details,omitempty"`
}

func WriteJSONError(w http.ResponseWriter, status int, err error, details interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   err.Error(),
		Details: details,
	})
}
