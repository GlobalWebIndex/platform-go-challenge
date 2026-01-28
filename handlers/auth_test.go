package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware_DisabledAuth(t *testing.T) {
	// Save original state
	origAuthEnabled := authEnabled
	origSecret := jwtSecret
	defer func() {
		authEnabled = origAuthEnabled
		jwtSecret = origSecret
	}()

	// Disable auth
	authEnabled = false
	jwtSecret = []byte("test-secret")

	called := false
	next := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	middleware := AuthMiddleware(next)

	// Request without auth header should still pass
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if !called {
		t.Fatal("next handler not called when auth is disabled")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	// Save original state
	origAuthEnabled := authEnabled
	origSecret := jwtSecret
	defer func() {
		authEnabled = origAuthEnabled
		jwtSecret = origSecret
	}()

	// Enable auth
	authEnabled = true
	jwtSecret = []byte("test-secret")

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	middleware := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if !containsString(rec.Body.String(), "Missing authorization header") {
		t.Fatalf("expected 'Missing authorization header', got %q", rec.Body.String())
	}
}

func TestAuthMiddleware_InvalidBearerFormat(t *testing.T) {
	// Save original state
	origAuthEnabled := authEnabled
	origSecret := jwtSecret
	defer func() {
		authEnabled = origAuthEnabled
		jwtSecret = origSecret
	}()

	authEnabled = true
	jwtSecret = []byte("test-secret")

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	middleware := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if !containsString(rec.Body.String(), "Invalid authorization format") {
		t.Fatalf("expected 'Invalid authorization format', got %q", rec.Body.String())
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Save original state
	origAuthEnabled := authEnabled
	origSecret := jwtSecret
	defer func() {
		authEnabled = origAuthEnabled
		jwtSecret = origSecret
	}()

	authEnabled = true
	jwtSecret = []byte("test-secret")

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	middleware := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if !containsString(rec.Body.String(), "Invalid token") {
		t.Fatalf("expected 'Invalid token', got %q", rec.Body.String())
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Save original state
	origAuthEnabled := authEnabled
	origSecret := jwtSecret
	defer func() {
		authEnabled = origAuthEnabled
		jwtSecret = origSecret
	}()

	testSecret := []byte("test-secret")
	authEnabled = true
	jwtSecret = testSecret

	called := false
	next := func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	}

	// Create a valid JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "user123",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString(testSecret)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	middleware := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if !called {
		t.Fatal("next handler not called with valid token")
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// Save original state
	origAuthEnabled := authEnabled
	origSecret := jwtSecret
	defer func() {
		authEnabled = origAuthEnabled
		jwtSecret = origSecret
	}()

	testSecret := []byte("test-secret")
	authEnabled = true
	jwtSecret = testSecret

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Create an expired JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "user123",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
	})

	tokenString, err := token.SignedString(testSecret)
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	middleware := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if !containsString(rec.Body.String(), "token is expired") {
		t.Fatalf("expected 'token is expired', got %q", rec.Body.String())
	}
}

func TestAuthMiddleware_WrongSecret(t *testing.T) {
	// Save original state
	origAuthEnabled := authEnabled
	origSecret := jwtSecret
	defer func() {
		authEnabled = origAuthEnabled
		jwtSecret = origSecret
	}()

	authEnabled = true
	jwtSecret = []byte("correct-secret")

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Create token with different secret
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   "user123",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString([]byte("wrong-secret"))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	middleware := AuthMiddleware(next)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rec := httptest.NewRecorder()

	middleware.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if !containsString(rec.Body.String(), "Invalid token") {
		t.Fatalf("expected 'Invalid token', got %q", rec.Body.String())
	}
}

// Helper to check if substring is in a string
func containsString(haystack, needle string) bool {
	return strings.Contains(haystack, needle)
}
