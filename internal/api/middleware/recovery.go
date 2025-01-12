package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"platform-go-challenge/internal/api/types"
	"runtime/debug"
)

// RecoveryMiddleware catches panics and returns 500 error response
func RecoveryMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Capture stack trace
					stack := debug.Stack()

					// Log error details
					slog.Error("panic recovered",
						slog.Any("error", err),
						slog.String("stack", string(stack)),
						slog.String("request_id", GetRequestID(r.Context())),
					)

					// Return 500 response
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusInternalServerError)

					response := types.NewErrorResponse(
						types.ErrCodeInternal,
						"An unexpected error occurred",
					)
					json.NewEncoder(w).Encode(response)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
