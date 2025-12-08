// Package app provides application setup, routing, and dependency wiring.
package app

import (
	"net/http"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/api/handler"
	"github.com/go-chi/chi/v5"
)

// NewRouter creates and configures the HTTP router with all routes and middleware.
func NewRouter() http.Handler {
	r := chi.NewRouter()

	// TODO: add middlewares (logging, metrics, auth)
	// r.Use(loggingMiddleware)
	// r.Use(metricsMiddleware)
	// r.Use(authMiddleware)

	r.Get("/health", handler.Health)

	return r
}
