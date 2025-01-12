package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddleware(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret")

	tests := []struct {
		name           string
		token          string
		expectedStatus int
		expectedUserID string
	}{
		{
			name:           "valid token",
			token:          createTestToken(t, "user123"),
			expectedStatus: http.StatusOK,
			expectedUserID: "user123",
		},
		{
			name:           "missing token",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			token:          "invalid.token.here",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "expired token",
			token:          createExpiredTestToken(t, "user123"),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectedUserID != "" {
					userID := GetUserID(r.Context())
					if userID != tt.expectedUserID {
						t.Errorf("expected userID %s, got %s", tt.expectedUserID, userID)
					}
				}
				w.WriteHeader(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			rr := httptest.NewRecorder()

			handler := AuthMiddleware()(nextHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}

func createTestToken(t *testing.T, userID string) string {
	token, err := GenerateToken(userID)
	if err != nil {
		t.Fatalf("failed to create test token: %v", err)
	}
	return token
}

func createExpiredTestToken(t *testing.T, userID string) string {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
		"iat": time.Now().Add(-2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		t.Fatalf("failed to create expired test token: %v", err)
	}
	return signedToken
}
