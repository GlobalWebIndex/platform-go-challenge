package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gwi/configuration"
	"gwi/constants"
	"gwi/routes"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJwtToken(userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
	})
	tokenString, _ := token.SignedString(configuration.JwtSecret)
	return tokenString
}

func TestFavoriteIntegration_AddAndList(t *testing.T) {
	// Setup test server
	router := routes.RegisterRoutes()
	userID := 1
	authorizationHeader := "Bearer " + GenerateJwtToken(userID)

	// Add favorite
	req := httptest.NewRequest(http.MethodPost, "/"+constants.ENDPOINT_FAVORITES, bytes.NewBufferString(`{"asset_id": 3}`))
	req.Header.Set("Authorization", authorizationHeader)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	// List favorites
	req2 := httptest.NewRequest(http.MethodGet, "/"+constants.ENDPOINT_FAVORITES, nil)
	req2.Header.Set("Authorization", authorizationHeader)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}

}
