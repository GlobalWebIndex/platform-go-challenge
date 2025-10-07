package asset_service_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"gwi/configuration"
	"gwi/constants"
	asset_service "gwi/services/asset"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func GenerateJwtToken(userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
	})
	tokenString, _ := token.SignedString(configuration.JwtSecret)
	return tokenString
}

func TestEditAsset_ValidRequest(t *testing.T) {
	userID := 1
	authorizationHeader := "Bearer " + GenerateJwtToken(userID)

	editBody := `{
        "data": {
          "new_data_field": "new_value",
          "new_data_field_2": 64000
        },
        "new_field": "example"
      }`
	body := bytes.NewBufferString(editBody)

	req := httptest.NewRequest(http.MethodPut, "/"+constants.ENDPOINT_ASSETS+"/1", body)
	req.Header.Set("Authorization", authorizationHeader)
	// Inject path variable manually
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	asset_service.EditAsset(w, req)
	resp := w.Result()
	defer resp.Body.Close()

	// Assert
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 Created, got %d", resp.StatusCode)
	}

}

func TestEditAsset_InvalidToken(t *testing.T) {
	editBody := `{
        "data": {
          "new_data_field": "new_value",
          "new_data_field_2": 64000
        },
        "new_field": "example"
      }`
	body := bytes.NewBufferString(editBody)
	req := httptest.NewRequest(http.MethodPut, "/"+constants.ENDPOINT_ASSETS+"/1", body)
	req.Header.Set("Authorization", "Bearer invalid-token")

	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	asset_service.EditAsset(w, req)

	// Assert
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestEditAsset_AssetNotFound(t *testing.T) {
	userID := 1
	authorizationHeader := "Bearer " + GenerateJwtToken(userID)

	editBody := `{
        "data": {
          "new_data_field": "new_value",
          "new_data_field_2": 64000
        },
        "new_field": "example"
      }`
	body := bytes.NewBufferString(editBody)
	req := httptest.NewRequest(http.MethodPut, "/"+constants.ENDPOINT_ASSETS+"/999", body)
	req.Header.Set("Authorization", authorizationHeader)
	req = mux.SetURLVars(req, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	asset_service.EditAsset(w, req)

	// Assert
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
