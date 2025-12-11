package httpx

import (
	"encoding/json"
	"net/http"
)

// ProblemDetail represents an RFC 9457 problem details response body.
type ProblemDetail struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

// WriteJSON writes a JSON response with the given status code.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteProblem writes a Problem Details JSON response with the given status,
// type, title, detail, and instance.
func WriteProblem(w http.ResponseWriter, status int, problemType, title, detail, instance string) {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ProblemDetail{
		Type:     problemType,
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	})
}
