
package handlers_test

import (
	"assetsApp/internal/handlers"
	"assetsApp/internal/models"
	favouriteServices "assetsApp/internal/services/favourite"
	"assetsApp/tests/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func TestFavouriteHandler_GetFavourites(t *testing.T) {
	userID := uuid.New()
	mockStore := &mocks.MockAssetStore{
		GetFavouritesFunc: func(uid uuid.UUID) []models.Favourite {
			if uid == userID {
				return []models.Favourite{{Asset: &models.Chart{ID: "test-chart"}}}
			}
			return nil
		},
	}
	service := favouriteServices.NewFavouriteService(mockStore)
	handler := handlers.NewFavouriteHandler(service)

	req, err := http.NewRequest("GET", "/users/"+userID.String()+"/favourites", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/favourites", handler.GetFavourites).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var favs []map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &favs); err != nil {
		t.Fatalf("could not parse response body: %v", err)
	}

	if len(favs) != 1 {
		t.Errorf("expected 1 favourite, got %d", len(favs))
	}
}

func TestFavouriteHandler_AddFavourite(t *testing.T) {
	userID := uuid.New()
	assetID := uuid.New().String()
	mockStore := &mocks.MockAssetStore{
		AddFavouriteFunc: func(uid uuid.UUID, aid string, assetType string) bool {
			return uid == userID && aid == assetID
		},
	}
	service := favouriteServices.NewFavouriteService(mockStore)
	handler := handlers.NewFavouriteHandler(service)

	body := map[string]string{"asset_type": "chart"}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "/users/"+userID.String()+"/favourites/"+assetID, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/favourites/{assetId}", handler.AddFavourite).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestFavouriteHandler_RemoveFavourite(t *testing.T) {
	userID := uuid.New()
	assetID := uuid.New().String()
	mockStore := &mocks.MockAssetStore{
		RemoveFavouriteFunc: func(uid uuid.UUID, aid string) bool {
			return uid == userID && aid == assetID
		},
	}
	service := favouriteServices.NewFavouriteService(mockStore)
	handler := handlers.NewFavouriteHandler(service)

	req, err := http.NewRequest("DELETE", "/users/"+userID.String()+"/favourites/"+assetID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/favourites/{assetId}", handler.RemoveFavourite).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestFavouriteHandler_GetFavourites_InvalidUserID(t *testing.T) {
	handler := handlers.NewFavouriteHandler(nil)

	req, err := http.NewRequest("GET", "/users/invalid-uuid/favourites", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/favourites", handler.GetFavourites).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}

func TestFavouriteHandler_AddFavourite_NotFound(t *testing.T) {
	userID := uuid.New()
	assetID := uuid.New().String()
	mockStore := &mocks.MockAssetStore{
		AddFavouriteFunc: func(uid uuid.UUID, aid string, assetType string) bool {
			return false // Simulate not found
		},
	}
	service := favouriteServices.NewFavouriteService(mockStore)
	handler := handlers.NewFavouriteHandler(service)

	body := map[string]string{"asset_type": "chart"}
	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "/users/"+userID.String()+"/favourites/"+assetID, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/favourites/{assetId}", handler.AddFavourite).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestFavouriteHandler_RemoveFavourite_NotFound(t *testing.T) {
	userID := uuid.New()
	assetID := uuid.New().String()
	mockStore := &mocks.MockAssetStore{
		RemoveFavouriteFunc: func(uid uuid.UUID, aid string) bool {
			return false // Simulate not found
		},
	}
	service := favouriteServices.NewFavouriteService(mockStore)
	handler := handlers.NewFavouriteHandler(service)

	req, err := http.NewRequest("DELETE", "/users/"+userID.String()+"/favourites/"+assetID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{userId}/favourites/{assetId}", handler.RemoveFavourite).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

