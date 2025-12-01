package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/models"
)

// JWT middleware
func TokenMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		var claims Claims

		token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			// Enforce expected signing algorithm (HS256)
			if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add userID and role to context
		r = r.WithContext(context.WithValue(r.Context(), cfg.ContextKeys.UserID, claims.UserID))
		r = r.WithContext(context.WithValue(r.Context(), cfg.ContextKeys.Role, claims.Role))

		next.ServeHTTP(w, r)
	})
}

func AdminAPIKeyMiddleware(cfg *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Admin-API-Key") != cfg.AdminAPIKey {
			http.Error(w, "Invalid Admin API Key", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetRoleFromContext(r *http.Request, cfg *config.Config) models.Role {
	val := r.Context().Value(cfg.ContextKeys.Role)
	role, ok := val.(models.Role)
	if !ok {
		return ""
	}
	return role
}
