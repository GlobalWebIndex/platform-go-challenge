package middleware

import (
	"net/http"
	"strings"

	"platform-go-challenge/internal/auth"
	"platform-go-challenge/internal/util"
)

// AuthMiddleware validates JWT tokens in incoming requests
func AuthMiddleware(jwtManager *auth.JWTManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				util.RespondError(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			// Extract token from "Bearer <token>" format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				util.RespondError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			tokenString := parts[1]

			// Validate token
			_, err := jwtManager.ValidateToken(tokenString)
			if err != nil {
				util.RespondError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// You can store the claims in the request context if needed
			// For now, just call the next handler
			next.ServeHTTP(w, r)
		})
	}
}