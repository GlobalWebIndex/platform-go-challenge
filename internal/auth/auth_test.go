package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeaderAuthenticator_DefaultHeaderName(t *testing.T) {
	t.Parallel()

	auth := NewHeaderAuthenticator("")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-User-ID", "user-123")

	userID, err := auth.Authenticate(req)
	if err != nil {
		t.Fatalf("Authenticate returned error: %v", err)
	}
	if userID != "user-123" {
		t.Fatalf("expected userID %q, got %q", "user-123", userID)
	}
}

func TestHeaderAuthenticator_CustomHeaderName(t *testing.T) {
	t.Parallel()

	auth := NewHeaderAuthenticator("X-Custom-User")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Custom-User", "user-xyz")

	userID, err := auth.Authenticate(req)
	if err != nil {
		t.Fatalf("Authenticate returned error: %v", err)
	}
	if userID != "user-xyz" {
		t.Fatalf("expected userID %q, got %q", "user-xyz", userID)
	}
}

func TestHeaderAuthenticator_MissingHeader_ReturnsErrUnauthenticated(t *testing.T) {
	t.Parallel()

	auth := NewHeaderAuthenticator("X-User-ID")

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := auth.Authenticate(req)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != ErrUnauthenticated {
		t.Fatalf("expected ErrUnauthenticated, got %v", err)
	}
}
