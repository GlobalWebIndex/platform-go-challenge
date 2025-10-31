package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/pkg/auth"
)

// Add the missing cxtKey type
type cxtKey string

const (
	UserClaimsKey cxtKey = "user_claims"
	UserRolesKey  cxtKey = "user_roles"
)

// AuthMiddleware verifies JWT tokens
func AuthMiddleware(keycloak *auth.KeycloakClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
				return
			}

			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, `{"error": "Invalid authorization header format"}`, http.StatusUnauthorized)
				return
			}

			token := parts[1]

			// Verify token
			claims, err := keycloak.VerifyToken(token)
			if err != nil {
				http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			// Get user roles
			roles := keycloak.GetUserRoles(claims)

			// Add claims and roles to context
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			ctx = context.WithValue(ctx, UserRolesKey, roles)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles, ok := r.Context().Value(UserRolesKey).([]string)
			if !ok {
				http.Error(w, `{"error": "No roles found in context"}`, http.StatusForbidden)
				return
			}

			// Check if user has the required role
			hasRole := false
			for _, role := range roles {
				if role == requiredRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, `{"error": "Insufficient permissions"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole checks if user has at least one of the required roles
func RequireAnyRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles, ok := r.Context().Value(UserRolesKey).([]string)
			if !ok {
				http.Error(w, `{"error": "No roles found in context"}`, http.StatusForbidden)
				return
			}

			// Check if user has at least one of the required roles
			hasRequiredRole := false
			for _, userRole := range roles {
				for _, requiredRole := range requiredRoles {
					if userRole == requiredRole {
						hasRequiredRole = true
						break
					}
				}
				if hasRequiredRole {
					break
				}
			}

			if !hasRequiredRole {
				http.Error(w, `{"error": "Insufficient permissions"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetClaimsFromContext retrieves claims from context
func GetClaimsFromContext(ctx context.Context) (*auth.CustomClaims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*auth.CustomClaims)
	return claims, ok
}

// GetRolesFromContext retrieves roles from context
func GetRolesFromContext(ctx context.Context) ([]string, bool) {
	roles, ok := ctx.Value(UserRolesKey).([]string)
	return roles, ok
}
