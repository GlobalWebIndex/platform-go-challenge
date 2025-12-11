// Package handler provides HTTP handlers for generic API endpoints (health, metrics, etc.).
package handler

import (
	"encoding/json"
	"net/http"
)

// healthResponse represents the health check response structure.
type healthResponse struct {
	Status string `json:"status"`
}

// Health returns the service health status as a JSON response.
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	resp := healthResponse{Status: "healthy"}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// TODO: Use a structured logger.
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
