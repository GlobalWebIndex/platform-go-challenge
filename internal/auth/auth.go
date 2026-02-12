package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userIDKey contextKey = "userID"

// AuthConfig holds JWT authentication configuration.
type AuthConfig struct {
	// Secret is the JWT signing secret. When empty, unsigned tokens may be accepted
	// if AllowUnsignedTokens is true.
	Secret string

	// AllowUnsignedTokens permits unsigned JWT tokens (alg=none) when true.
	// This should ONLY be enabled for local development and testing.
	AllowUnsignedTokens bool
}

// JWTMiddleware returns HTTP middleware that validates a JWT from the
// Authorization header and places the "sub" claim into the request context.
//
// When Secret is empty AND AllowUnsignedTokens is true, unsigned tokens (alg=none)
// are accepted — this is intended for local development and testing only.
// When Secret is non-empty, only HS256-signed tokens are accepted.
// When Secret is empty AND AllowUnsignedTokens is false, all requests are rejected.
func JWTMiddleware(cfg AuthConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString, ok := extractBearerToken(r)
			if !ok {
				http.Error(w, `{"error":"missing or malformed Authorization header"}`, http.StatusUnauthorized)
				return
			}

			// Reject unsigned tokens unless explicitly allowed
			if cfg.Secret == "" && !cfg.AllowUnsignedTokens {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			claims, err := parseToken(tokenString, cfg.Secret)
			if err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"%s"}`, err.Error()), http.StatusUnauthorized)
				return
			}

			sub, err := claims.GetSubject()
			if err != nil || sub == "" {
				http.Error(w, `{"error":"token missing sub claim"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, sub)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext returns the user ID stored by JWTMiddleware.
// Returns an empty string if no user ID is present.
func UserIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(userIDKey).(string)
	return v
}

// extractBearerToken pulls the token from "Authorization: Bearer <token>".
func extractBearerToken(r *http.Request) (string, bool) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return "", false
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", false
	}
	return parts[1], true
}

// parseToken validates the JWT string. If secret is empty, only alg=none is
// accepted (for dev/test). Otherwise HS256 is required.
func parseToken(tokenString, secret string) (jwt.MapClaims, error) {
	if secret == "" {
		// Development mode: accept unsigned tokens only.
		token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})
		if err != nil {
			return nil, fmt.Errorf("invalid token: %w", err)
		}
		if token.Method.Alg() != "none" {
			return nil, fmt.Errorf("no jwt secret configured; only unsigned tokens (alg=none) are accepted")
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("invalid token claims")
		}
		return claims, nil
	}

	// Production mode: require HS256.
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}
