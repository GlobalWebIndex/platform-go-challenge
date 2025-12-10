package logger

import (
	"net/http"
	"time"
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

// RequestLoggingMiddleware returns a middleware that logs basic request
// information (method, path, status code, duration).
//
// In a real system this could be extended with request IDs, user IDs,
// correlation IDs, etc.
func RequestLoggingMiddleware(lg *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{
				ResponseWriter: w,
				status:         http.StatusOK, // default if WriteHeader not called
			}

			next.ServeHTTP(rec, r)

			duration := time.Since(start)
			lg.Info(
				"http_request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rec.status,
				"duration_ms", duration.Milliseconds(),
			)
		})
	}
}
