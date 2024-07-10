package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	e := echo.New()

	// Define a handler that will be protected by the middleware
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	}

	// Apply the AuthMiddleware to the handler
	protectedRoute := AuthMiddleware(handler)

	// Test with a valid token
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, "gwi-token-12345")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute the request
	if assert.NoError(t, protectedRoute(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "success", rec.Body.String())
	}

	// Test with an invalid token
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, "invalid-token")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	// Execute the request
	err := protectedRoute(c)
	if assert.Error(t, err) {
		he, ok := err.(*echo.HTTPError)
		if assert.True(t, ok) {
			assert.Equal(t, http.StatusUnauthorized, he.Code)
			assert.Equal(t, "invalid token", he.Message)
		}
	}

	// Test with no token
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)

	// Execute the request
	err = protectedRoute(c)
	if assert.Error(t, err) {
		he, ok := err.(*echo.HTTPError)
		if assert.True(t, ok) {
			assert.Equal(t, http.StatusUnauthorized, he.Code)
			assert.Equal(t, "missing token", he.Message)
		}
	}
}
