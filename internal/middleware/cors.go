package middleware

import (
	"net/http"
	"strings"
)

// CORS returns a middleware that allows requests from the given origins.
// Example: CORS([]string{"http://localhost:8081"})
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
	// normalize to lower-case for comparison
	lower := make([]string, len(allowedOrigins))
	for i, o := range allowedOrigins {
		lower[i] = strings.ToLower(o)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := strings.ToLower(r.Header.Get("Origin"))

			// allow if origin is in the allowlist (or if empty, skip CORS)
			allow := false
			if origin != "" {
				for _, o := range lower {
					if origin == o || o == "*" {
						allow = true
						break
					}
				}
			}

			if allow {
				w.Header().Set("Access-Control-Allow-Origin", origin) // or "*" if you used wildcard
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
				// allow headers we use: Content-Type, X-API-Key, and common fetch headers
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-API-Key, Accept, Authorization")
				// If you need cookies, also set: w.Header().Set("Access-Control-Allow-Credentials","true")
			}

			// Handle preflight
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
