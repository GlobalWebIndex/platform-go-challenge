package middleware

import (
	"log"
	"net/http"

	"gwi/utils"
)

// RecoverMiddleware is a middleware that recovers from panics in the HTTP handler chain.
// It logs the panic and sends an internal server error response to the client.
//
// Parameters:
//   - next: The next http.Handler in the chain to be executed.
//
// Returns:
//
//	An http.Handler that wraps the provided handler with recovery functionality.
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				utils.SendError(w, utils.ErrInternalServer("Internal server error"))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
