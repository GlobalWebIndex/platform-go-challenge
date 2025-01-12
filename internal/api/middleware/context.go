package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var contextTimeOut time.Duration = 15 * time.Second

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
	StartTimeKey contextKey = "start_time"
	JWTClaimsKey contextKey = "jwt_claims_key"
)

// ContextMiddleware injects standard context values into every request
func ContextMiddleware(ctx context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := GenerateRequestID()

			rCtx, cancel := context.WithTimeout(ctx, contextTimeOut)
			defer cancel()

			go func() {
				defer cancel()

				select {
				case <-ctx.Done():
				case <-rCtx.Done():
				}
			}()

			ctx = context.WithValue(ctx, RequestIDKey, requestID)
			ctx = context.WithValue(ctx, StartTimeKey, time.Now())

			w.Header().Set("X-Request-ID", requestID)

			r = r.WithContext(rCtx)
			next.ServeHTTP(w, r)
		})
	}
}

// GetRequestID retrieves request ID from context
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// GetClaims retrieves JWT claims from context
func GetClaims(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(JWTClaimsKey).(*Claims)
	return claims, ok
}

// GetUserID retrieves user ID from JWT claims
func GetUserID(ctx context.Context) string {
	if claims, ok := GetClaims(ctx); ok {
		return claims.UserID
	}
	return ""
}

// GenerateRequestID creates a unique request ID
func GenerateRequestID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), uuid.New().String())
}
