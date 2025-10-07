package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gwi/constants"
	"gwi/routes"

	"github.com/gorilla/mux"
)

func TestAssetIntegration_AddAndList(t *testing.T) {
	// Setup test server
	router := routes.RegisterRoutes()
	userID := 1
	authorizationHeader := "Bearer " + GenerateJwtToken(userID)

	// Edit asset
	req := httptest.NewRequest(http.MethodPut, "/"+constants.ENDPOINT_ASSETS+"/1", bytes.NewBufferString(`{"asset_id": 3}`))
	req.Header.Set("Authorization", authorizationHeader)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

}
