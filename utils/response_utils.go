package utils

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// SendSuccess writes a JSON response with a success status and the provided data.
// It sets the Content-Type header to application/json and writes an HTTP status code of 200 (OK).
//
// Parameters:
//   - w: The http.ResponseWriter used to construct the response.
//   - data: The data to be included in the response body.
func SendSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{Success: true, Data: data})
}

// SendError writes an error response to the HTTP response writer.
// It checks if the provided error is of type APIError and sets the appropriate
// HTTP status code and message. If the error is not of type APIError, it
// defaults to a 500 Internal Server Error response with a generic message.
//
// Parameters:
//   - w: The http.ResponseWriter to write the response to.
//   - err: The error to be processed and sent in the response.
func SendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	if apiErr, ok := err.(*APIError); ok {
		w.WriteHeader(apiErr.Status)
		json.NewEncoder(w).Encode(APIResponse{Success: false, Message: apiErr.Message})
		return
	}

	// fallback for unknown errors
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(APIResponse{Success: false, Message: "Internal server error"})
}
