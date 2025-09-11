package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var jwtSecret = []byte("supersecretkey") // Change this in production

// GenerateJWT creates a JWT token for a given user ID
func GenerateJWT(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // 24h expiry
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// TokenHandler issues a JWT for a given user_id (POST /token)
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user_id", http.StatusBadRequest)
		return
	}
	token, err := GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}
	resp := struct {
		Token string `json:"token"`
	}{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func extractUserIDFromToken(r *http.Request) (uuid.UUID, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return uuid.Nil, http.ErrNoCookie
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return uuid.Nil, http.ErrNoCookie
	}
	tokenStr := parts[1]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return uuid.Nil, http.ErrNoCookie
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, http.ErrNoCookie
	}
	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, http.ErrNoCookie
	}
	return uuid.Parse(userIDStr)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := extractUserIDFromToken(r)
		if err != nil || userID == uuid.Nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Attach userID to context for handlers
		ctx := r.Context()
		ctx = contextWithUserID(ctx, userID)
		next(w, r.WithContext(ctx))
	}
}

type contextKey string

const userIDKey contextKey = "userID"

func contextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func getUserIDFromContext(r *http.Request) uuid.UUID {
	val := r.Context().Value(userIDKey)
	if id, ok := val.(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}
