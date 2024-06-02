package middleware

import (
	"context"
	"net/http"
	"strings"
	"user-favorites-api/store"

	"github.com/dgrijalva/jwt-go"
)

func extractAssetIDFromPath(r *http.Request) string {
    // Get the path from the request URL
    path := r.URL.Path

    // Split the path by '/' to get the individual segments
    segments := strings.Split(path, "/")

    // The asset ID is typically the last segment in the path
    assetID := segments[len(segments)-1]

    return assetID
}

func JWTAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "authorization header required", http.StatusUnauthorized)
            return
        }

        bearerToken := strings.Split(authHeader, " ")
        if len(bearerToken) != 2 {
            http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
            return
        }

        tokenStr := bearerToken[1]
        claims := &store.Claims{}

        token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            return store.JwtKey, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "username", claims.Username)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
