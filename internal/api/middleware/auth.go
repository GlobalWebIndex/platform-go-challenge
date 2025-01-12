package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"platform-go-challenge/internal/api/types"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    string
	ExpiresAt int64
	IssuedAt  int64
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// AuthMiddleware checks for valid JWT tokens
func AuthMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/swagger") {
				next.ServeHTTP(w, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				sendAuthError(w, types.ErrCodeUnauthorized, "No token provided")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims := &jwt.RegisteredClaims{}

			token, err := jwt.ParseWithClaims(
				tokenString,
				claims, func(token *jwt.Token) (interface{}, error) {
					// Verify signing method
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return jwtSecret, nil
				},
			)

			if err != nil || !token.Valid {
				sendAuthError(w, types.ErrCodeInvalidToken, "Invalid token")
				return
			}

			// Check expiry
			if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
				sendAuthError(w, types.ErrCodeInvalidToken, "Token expired")
				return
			}

			customClaims := &Claims{
				UserID:    claims.Subject,
				ExpiresAt: claims.ExpiresAt.Unix(),
				IssuedAt:  claims.IssuedAt.Unix(),
			}

			ctx := context.WithValue(r.Context(), JWTClaimsKey, customClaims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GenerateToken creates a new JWT token
func GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24 * 365).Unix(),
		"iat": time.Now().Unix(),
	})

	return token.SignedString(jwtSecret)
}

func sendAuthError(w http.ResponseWriter, code string, message string) {
	response := types.NewErrorResponse(code, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(response)
}
