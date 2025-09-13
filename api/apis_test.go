package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"platform-go-challenge/db"
	"platform-go-challenge/models"
)

func resetData() {
	database = db.NewDB()
}

func createTestAsset(userId string, assetId string, description string) models.Asset {
	asset := models.Asset{
		ID:          assetId,
		Description: description,
	}

	err := database.AddFavorite(userId, asset)
	if err != nil {
		panic("failed to add favorite: " + err.Error())
	}
	return asset
}

func TestPostandGetFavorites(test *testing.T) {

	resetData()

	asset := models.Asset{ID: "12", Description: "Just an asset"}
	body, _ := json.Marshal(asset)

	request, _ := http.NewRequest("POST", "/users/user123/12", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	FavoritesAPIs(response, request)

	if response.Result().StatusCode != http.StatusCreated {
		test.Errorf("Expected status code 201, got %d", response.Result().StatusCode)
	}

	request, _ = http.NewRequest("GET", "/users/user123/12", nil)
	response = httptest.NewRecorder()
	FavoritesAPIs(response, request)

	if response.Result().StatusCode != http.StatusOK {
		test.Errorf("Expected status code 200, got %d", response.Result().StatusCode)
	}

	var favorites []models.Asset
	err := json.NewDecoder(response.Result().Body).Decode(&favorites)
	if err != nil {
		test.Errorf("failed to decode favorites %v", err)
	}

	if len(favorites) != 1 || favorites[0].ID != "12" {
		test.Errorf("Expected 1 favorite, got %d", len(favorites))
	}
}

func TestUpdateFavorite(test *testing.T) {
	resetData()
	asset := createTestAsset("user123", "12", "Existing old asset")

	body, _ := json.Marshal(asset)
	request := httptest.NewRequest("PUT", "/users/user123/12", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	FavoritesAPIs(response, request)

	updatedAsset := models.Asset{Description: "Now asset is updated"}
	body, _ = json.Marshal(updatedAsset)
	request = httptest.NewRequest("PUT", "/users/user123/12", bytes.NewBuffer(body))
	response = httptest.NewRecorder()
	FavoritesAPIs(response, request)

	if response.Result().StatusCode != http.StatusOK {
		test.Errorf("Expected status code 200, got %d", response.Result().StatusCode)
	}

	var returnedAsset models.Asset
	err := json.NewDecoder(response.Body).Decode(&returnedAsset)
	if err != nil {
		test.Errorf("failed to decode favorites %v", err)
	}
	if returnedAsset.Description != "Now asset is updated" {
		test.Errorf("Expected 'Now asset is updated', got %s", returnedAsset.Description)
	}
}
