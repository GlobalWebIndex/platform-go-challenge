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

func createTestAsset(userId string, assetId string, typeId string, description string, data map[string]interface{}) models.Asset {
	asset := models.Asset{
		ID:          assetId,
		TypeID:      typeId,
		Description: description,
		Data:        data,
	}

	err := database.AddFavorite(userId, asset)
	if err != nil {
		panic("failed to add favorite: " + err.Error())
	}
	return asset
}

func TestPostAndGetFavorites(test *testing.T) {

	resetData()

	//POST an asset
	asset := models.Asset{ID: "12", TypeID: "Chart", Description: "Just an asset chart"}
	body, _ := json.Marshal(asset)

	request, _ := http.NewRequest("POST", "/users/user123/12", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	FavoritesAPIs(response, request)

	if response.Result().StatusCode != http.StatusCreated {
		test.Errorf("Expected status code 201, got %d", response.Result().StatusCode)
	}

	// Call GET on all favorites
	request, _ = http.NewRequest("GET", "/users/user123", nil)
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

func TestGetSingleAsset(test *testing.T) {
	resetData()

	chartData := map[string]interface{}{"title": "Sales", "x_axis": "Month", "y_axis": "Revenue", "data": []int{100, 200, 300}}
	createTestAsset("user123", "12", "Chart", "Chart", chartData)

	insightData := map[string]interface{}{"text": "40% of millennials spend more than 3 hours on social media daily"}
	createTestAsset("user123", "13", "Insight", "Insight", insightData)

	request, _ := http.NewRequest("GET", "/users/user123/13", nil)
	response := httptest.NewRecorder()
	FavoritesAPIs(response, request)

	if response.Result().StatusCode != http.StatusOK {
		test.Errorf("Expected status code 200, got %d", response.Result().StatusCode)
	}

	var filteredAsset models.Asset
	err := json.NewDecoder(response.Result().Body).Decode(&filteredAsset)
	if err != nil {
		test.Errorf("failed to decode favorites %v", err)
	}
	if filteredAsset.ID != "13" || filteredAsset.Description != "Insight" {
		test.Errorf("Expected 1 Insight, got %s", filteredAsset.Description)
	}

}

func TestUpdateFavorite(test *testing.T) {
	resetData()
	chartData := map[string]interface{}{"title": "Sales", "x_axis": "Month", "y_axis": "Revenue", "data": []int{100, 200, 300}}
	asset := createTestAsset("user123", "12", "Chart", "Existing old chart", chartData)

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

func TestDeleteAsset(test *testing.T) {
	resetData()
	chartData := map[string]interface{}{"title": "Sales", "x_axis": "Month", "y_axis": "Revenue", "data": []int{100, 200, 300}}
	createTestAsset("user123", "12", "Chart", "Users favorite chart", chartData)

	insightData := map[string]interface{}{"text": "40% of millennials spend more than 3 hours on social media daily"}
	createTestAsset("user123", "13", "Insight", "Users favorite insight", insightData)

	audienceData := map[string]interface{}{"gender": "Male", "birth_country": "USA", "age_group": "24-35", "social_hours": "3+", "purchases_last": "5"}
	createTestAsset("user123", "14", "Audience", "Users favorite audience data", audienceData)

	request, _ := http.NewRequest("DELETE", "/users/user123/12", nil)
	response := httptest.NewRecorder()
	FavoritesAPIs(response, request)

	if response.Result().StatusCode != http.StatusNoContent {
		test.Errorf("Expected status code 204, got %d", response.Result().StatusCode)
	}

	// check that only 2 assets should exist
	request, _ = http.NewRequest("GET", "/users/user123", nil)
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

	if len(favorites) != 2 {
		test.Errorf("Expected 2 favorite, got %d", len(favorites))
	}

}
