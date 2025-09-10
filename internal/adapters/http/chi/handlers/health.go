package handlers

import (
	"net/http"
)

// HealthHandler responds with 200 OK (simple readiness).
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler { return &HealthHandler{} }

// ServeHTTP godoc
// @Summary      Health check
// @Description  Simple readiness probe.
// @Tags         health
// @Produce      json
// @Success      200  {object} handlers.HealthResponse
// @Router       /healthz [get]
func (handler *HealthHandler) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write([]byte(`{"ok":true}`))
}
