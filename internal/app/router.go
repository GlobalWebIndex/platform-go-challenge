// Package app provides application setup, routing, and dependency wiring.
package app

import (
	"net/http"

	roothandler "github.com/GlobalWebIndex/platform-go-challenge/internal/api/handler"
	v1handler "github.com/GlobalWebIndex/platform-go-challenge/internal/api/v1/handler"
	v1middleware "github.com/GlobalWebIndex/platform-go-challenge/internal/api/v1/middleware"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/auth"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/domain/favourites"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/infrastructure/logger"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/repository/memory"
	"github.com/go-chi/chi/v5"
)

// NewRouter creates and configures the HTTP router with all routes and middleware.
func NewRouter() http.Handler {
	r := chi.NewRouter()

	log := logger.New()

	// TODO: add middlewares (logging, metrics, auth)
	r.Use(logger.RequestLoggingMiddleware(log))
	// r.Use(metricsMiddleware)

	// Unversioned endpoints (health, metrics, etc).
	r.Get("/health", roothandler.Health)

	// --- Dependency wiring for favourites API (v1) ---

	// Repository: in-memory implementation for the challenge.
	favRepo := memory.NewFavouritesRepository()

	// Domain service on top of the repository.
	favSvc := favourites.NewService(favRepo)

	// HTTP handler for the v1 favourites API.
	favHandler := v1handler.NewFavouritesHandler(favSvc)

	// Authenticator: simple header-based for this exercise, pluggable later.
	headerAuth := auth.NewHeaderAuthenticator("X-User-ID")

	// Versioned API routes.
	r.Route("/v1", func(r chi.Router) {
		// Apply auth middleware to all v1 routes.
		r.Use(v1middleware.NewAuthUserMiddleware(headerAuth))

		// Registers:
		//   GET    /users/{userId}/favourites
		//   POST   /users/{userId}/favourites
		//   DELETE /users/{userId}/favourites/{assetId}
		//   PATCH  /users/{userId}/favourites/{assetId}
		favHandler.RegisterRoutes(r)
	})

	return r
}
