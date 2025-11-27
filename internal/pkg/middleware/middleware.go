package middleware

import (
	"context"
	"github.com/akazantzidis/gwi-ass/internal/pkg/helper"
	"net/http"
	"strings"
)

type contextKey string

const ContextUserKey contextKey = "auth_user"

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "invalid Authorization format", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := helper.ParseAndValidateToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUserKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClaimsFromContext(ctx context.Context) *helper.CustomClaims {
	if v := ctx.Value(ContextUserKey); v != nil {
		if c, ok := v.(*helper.CustomClaims); ok {
			return c
		}
	}
	return nil
}

func HasRole(c *helper.CustomClaims, role string) bool {
	if c == nil {
		return false
	}
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := ClaimsFromContext(r.Context())
			if claims == nil || !HasRole(claims, role) {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
