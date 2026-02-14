package http

import (
	"net/http"
	"strings"
	"time"
)

const (
	DefaultTokenExpiry = time.Hour
	DefaultListLimit   = 20
	MaxListLimit       = 100
)

func AuthorizeUser(w http.ResponseWriter, r *http.Request, verifier Verifier, pathUserID string) bool {
	tokenStr := extractBearerToken(r)
	if tokenStr == "" {
		WriteError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "missing or invalid Authorization header")
		return false
	}
	userID, err := verifier.Verify(tokenStr)
	if err != nil {
		WriteErrorWithLog(w, r, http.StatusUnauthorized, ErrCodeUnauthorized, "invalid or expired token", err)
		return false
	}
	if userID != pathUserID {
		WriteError(w, http.StatusForbidden, ErrCodeForbidden, "cannot access another user's favourites")
		return false
	}
	return true
}

func extractBearerToken(r *http.Request) string {
	v := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if len(v) < len(prefix) || !strings.EqualFold(v[:len(prefix)], prefix) {
		return ""
	}
	return strings.TrimSpace(v[len(prefix):])
}
