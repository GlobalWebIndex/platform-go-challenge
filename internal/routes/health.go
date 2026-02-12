package routes

import (
	"net/http"

	"github.com/giannis84/platform-go-challenge/internal/config"
	"github.com/giannis84/platform-go-challenge/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

// RegisterHealthRoutes creates the health check endpoints.
func RegisterHealthRoutes(rateCfg config.RateLimitConfig) func(r chi.Router) {
	return func(r chi.Router) {
		// Apply IP-based rate limiting if configured
		if rateCfg.Requests > 0 && rateCfg.Window > 0 {
			r.Use(httprate.Limit(
				rateCfg.Requests,
				rateCfg.Window,
				httprate.WithKeyFuncs(httprate.KeyByIP),
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusTooManyRequests)
					w.Write([]byte("rate limit exceeded"))
				}),
			))
		}

		r.Get("/health/live", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})

		r.Get("/health/ready", func(w http.ResponseWriter, r *http.Request) {
			if err := database.PingDB(r.Context()); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("database not ready"))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ready"))
		})
	}
}
