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
	metricsinfra "github.com/GlobalWebIndex/platform-go-challenge/internal/infrastructure/metrics"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/repository/memory"
	"github.com/go-chi/chi/v5"
)

// NewRouter creates and configures the HTTP router with all routes and middleware.
func NewRouter() http.Handler {
	r := chi.NewRouter()

	// Infrastructure: logger and metrics.
	log := logger.New()
	metrics := metricsinfra.New()

	// Global middlewares (applied to all routes).
	r.Use(logger.RequestLoggingMiddleware(log))
	r.Use(metricsinfra.RequestMetricsMiddleware(metrics))

	// Unversioned endpoints (health, metrics, docs).
	r.Get("/health", roothandler.Health)

	// NOTE: /metrics is intentionally unauthenticated for this exercise so
	// that Prometheus/Grafana in docker-compose can scrape it easily.
	// TODO: In a real environment, protect /metrics with internal auth
	//       (mTLS, IP allowlist, separate admin gateway, etc.).
	r.Handle("/metrics", metrics.Handler())

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
		// Router that owns the {userId} path parameter.
		r.Route("/users/{userId}", func(r chi.Router) {
			// Apply auth middleware to all user-scoped routes so it can
			// see {userId} via chi.URLParam and enforce the match against
			// the authenticated user.
			r.Use(v1middleware.NewAuthUserMiddleware(headerAuth))

			// Registers:
			//   GET    /v1/users/{userId}/favourites
			//   POST   /v1/users/{userId}/favourites
			//   DELETE /v1/users/{userId}/favourites/{assetId}
			//   PATCH  /v1/users/{userId}/favourites/{assetId}
			r.Get("/favourites", favHandler.ListUserFavourites)
			r.Post("/favourites", favHandler.AddFavourite)
			r.Delete("/favourites/{assetId}", favHandler.RemoveFavourite)
			r.Patch("/favourites/{assetId}", favHandler.EditFavouriteDescription)
		})
	})

	return r
}
