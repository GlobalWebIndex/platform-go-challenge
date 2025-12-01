package middleware

import (
	"os"
    "context"
    "encoding/json"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v5"
    "github.com/dintzeler/platform-go-challenge/customerrors"
)

type contextKey string

const UserIDKey contextKey = "user_id"

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) 

type Claims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}

func JWTMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(&customerrors.ValidationError{
                Message:   "Missing authorization header",
                ErrorCode: "MISSING_AUTH_HEADER",
            })
            return
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(&customerrors.ValidationError{
                Message:   "Invalid authorization header format",
                ErrorCode: "INVALID_AUTH_FORMAT",
            })
            return
        }

        tokenString := tokenParts[1]

        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(&customerrors.ValidationError{
                Message:   "Invalid or expired token",
                ErrorCode: "INVALID_TOKEN",
            })
            return
        }

        // Extract claims
        claims, ok := token.Claims.(*Claims)
        if !ok {
            w.WriteHeader(http.StatusUnauthorized)
            json.NewEncoder(w).Encode(&customerrors.ValidationError{
                Message:   "Invalid token claims",
                ErrorCode: "INVALID_CLAIMS",
            })
            return
        }

        // Add user ID to context
        ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
        r = r.WithContext(ctx)

        // Call next handler
        next(w, r)
    }
}