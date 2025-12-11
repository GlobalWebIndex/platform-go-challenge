package auth

import (
	"errors"
	"net/http"
)

// ErrUnauthenticated is returned when authentication fails.
var ErrUnauthenticated = errors.New("unauthenticated")

// Authenticator defines how we authenticate a request and obtain an
// opaque user identifier.
//
// The rest of the system only deals with the returned userID string and
// does not care about how authentication is implemented (header, JWT,
// OAuth2, etc.).
type Authenticator interface {
	Authenticate(r *http.Request) (userID string, err error)
}

// HeaderAuthenticator implements Authenticator by reading an opaque
// user identifier from a specific HTTP header (e.g. X-User-ID).
//
// This is intentionally simple for the exercise but keeps the shape of
// a real pluggable authenticator.
type HeaderAuthenticator struct {
	HeaderName string
}

// NewHeaderAuthenticator creates a HeaderAuthenticator for the given header.
// Defaults to "X-User-ID" if headerName is empty.
func NewHeaderAuthenticator(headerName string) HeaderAuthenticator {
	if headerName == "" {
		headerName = "X-User-ID"
	}
	return HeaderAuthenticator{HeaderName: headerName}
}

// Authenticate returns the header value as the user ID, or ErrUnauthenticated if missing.
func (a HeaderAuthenticator) Authenticate(r *http.Request) (string, error) {
	userID := r.Header.Get(a.HeaderName)
	if userID == "" {
		return "", ErrUnauthenticated
	}
	return userID, nil
}
