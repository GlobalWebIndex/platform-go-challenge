package handlers_test

import (
	"assetsApp/internal/handlers"
	"assetsApp/internal/models"
	assetServices "assetsApp/internal/services/asset"
	"assetsApp/tests/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestAssetHandler_GetAssets(t *testing.T) {
	userID := uuid.New()
	mockStore := &mocks.MockAssetStore{
		GetFunc: func(uid uuid.UUID) []models.Asset {
			if uid == userID {
				return []models.Asset{&models.Chart{ID: "test-chart"}}
			}
			return nil
		},
	}
	service := assetServices.NewAssetService(mockStore)
	handler := handlers.NewAssetHandler(service)

	req, err := http.NewRequest("GET", "/users/"+userID.String()+"/assets", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/assets", handler.GetAssets).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var assets []map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &assets); err != nil {
		t.Fatalf("could not parse response body: %v", err)
	}

	if len(assets) != 1 {
		t.Errorf("expected 1 asset, got %d", len(assets))
	}
}

func TestAssetHandler_AddAsset(t *testing.T) {
	userID := uuid.New()
	assetID := uuid.New().String()
	mockStore := &mocks.MockAssetStore{
		AddFunc: func(uid uuid.UUID, asset models.Asset) {
			// No-op for this test
		},
	}
	service := assetServices.NewAssetService(mockStore)
	handler := handlers.NewAssetHandler(service)

	asset := map[string]interface{}{
		"id":           assetID,
		"type":         "chart",
		"title":        "Test Chart",
		"description":  "A test chart",
		"x_axis_title": "X",
		"y_axis_title": "Y",
		"data":         []map[string]interface{}{},
	}
	body, _ := json.Marshal(asset)

	req, err := http.NewRequest("POST", "/users/"+userID.String()+"/assets", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/assets", handler.AddAsset).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestAssetHandler_EditAsset(t *testing.T) {
	userID := uuid.New()
	assetID := uuid.New().String()
	mockStore := &mocks.MockAssetStore{
		EditDescriptionFunc: func(uid uuid.UUID, aid string, description string) bool {
			return uid == userID && aid == assetID
		},
	}
	service := assetServices.NewAssetService(mockStore)
	handler := handlers.NewAssetHandler(service)

	editBody := map[string]string{"description": "New Description"}
	body, _ := json.Marshal(editBody)

	req, err := http.NewRequest("PUT", "/users/"+userID.String()+"/assets/"+assetID, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/assets/{assetId}", handler.EditAsset).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAssetHandler_RemoveAsset(t *testing.T) {
	userID := uuid.New()
	assetID := uuid.New().String()
	mockStore := &mocks.MockAssetStore{
		RemoveFunc: func(uid uuid.UUID, aid string) bool {
			return uid == userID && aid == assetID
		},
	}
	service := assetServices.NewAssetService(mockStore)
	handler := handlers.NewAssetHandler(service)

	req, err := http.NewRequest("DELETE", "/users/"+userID.String()+"/assets/"+assetID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/assets/{assetId}", handler.RemoveAsset).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestAssetHandler_GetAssets_InvalidUserID(t *testing.T) {
	handler := handlers.NewAssetHandler(nil)

	req, err := http.NewRequest("GET", "/users/invalid-uuid/assets", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/assets", handler.GetAssets).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestAssetHandler_AddAsset_InvalidBody(t *testing.T) {
	userID := uuid.New()
	handler := handlers.NewAssetHandler(nil)

	req, err := http.NewRequest("POST", "/users/"+userID.String()+"/assets", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/assets", handler.AddAsset).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestAssetHandler_EditAsset_NotFound(t *testing.T) {
	userID := uuid.New()
	assetID := uuid.New().String()
	mockStore := &mocks.MockAssetStore{
		EditDescriptionFunc: func(uid uuid.UUID, aid string, description string) bool {
			return false // Simulate not found
		},
	}
	service := assetServices.NewAssetService(mockStore)
	handler := handlers.NewAssetHandler(service)

	editBody := map[string]string{"description": "New Description"}
	body, _ := json.Marshal(editBody)

	req, err := http.NewRequest("PUT", "/users/"+userID.String()+"/assets/"+assetID, bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/assets/{assetId}", handler.EditAsset).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestAssetHandler_RemoveAsset_NotFound(t *testing.T) {
	userID := uuid.New()
	assetID := uuid.New().String()
	mockStore := &mocks.MockAssetStore{
		RemoveFunc: func(uid uuid.UUID, aid string) bool {
			return false // Simulate not found
		},
	}
	service := assetServices.NewAssetService(mockStore)
	handler := handlers.NewAssetHandler(service)

	req, err := http.NewRequest("DELETE", "/users/"+userID.String()+"/assets/"+assetID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/assets/{assetId}", handler.RemoveAsset).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

