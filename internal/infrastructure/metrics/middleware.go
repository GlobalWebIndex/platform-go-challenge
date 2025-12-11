package metrics

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// statusRecorder wraps http.ResponseWriter to capture the status code.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

// RequestMetricsMiddleware records HTTP request metrics (total and duration).
// Uses chi route pattern for endpoint label when available, otherwise falls back to URL path.
func RequestMetricsMiddleware(m *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{
				ResponseWriter: w,
				status:         http.StatusOK, // default if WriteHeader is not called
			}

			next.ServeHTTP(rec, r)

			duration := time.Since(start)

			// Prefer the route pattern, fall back to raw path if not available.
			pathLabel := r.URL.Path
			if rctx := chi.RouteContext(r.Context()); rctx != nil {
				if pattern := rctx.RoutePattern(); pattern != "" {
					pathLabel = pattern
				}
			}

			m.observeRequest(r.Method, pathLabel, rec.status, duration)
		})
	}
}
