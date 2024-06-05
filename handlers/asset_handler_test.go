package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meletios/gwi-engineering-challenge/models"
)

func TestGetFavorites(t *testing.T) {
	// Setup
	UsersFavorites["user1"] = []models.Asset{
		{Chart: &models.Chart{ID: "chart1", Title: "Sales Data", AxesTitles: "Month vs Sales", Data: "Jan:100,Feb:150", Description: "Monthly sales data"}},
	}

	req, err := http.NewRequest("GET", "/favorites/user1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/favorites/{userID}", GetFavorites)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var assets []models.Asset
	err = json.NewDecoder(rr.Body).Decode(&assets)
	if err != nil {
		t.Fatal(err)
	}

	if len(assets) != 1 || assets[0].Chart.ID != "chart1" {
		t.Errorf("handler returned unexpected body: got %v", rr.Body.String())
	}
}

func TestAddFavorite(t *testing.T) {
	asset := models.Asset{Chart: &models.Chart{ID: "chart2", Title: "Revenue Data", AxesTitles: "Month vs Revenue", Data: "Jan:200,Feb:250", Description: "Monthly revenue data"}}
	body, err := json.Marshal(asset)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/favorites/user2", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/favorites/{userID}", AddFavorite)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if len(UsersFavorites["user2"]) != 1 || UsersFavorites["user2"][0].Chart.ID != "chart2" {
		t.Errorf("handler did not add favorite correctly: got %v", UsersFavorites["user2"])
	}
}

func TestRemoveFavorite(t *testing.T) {
	UsersFavorites["user3"] = []models.Asset{
		{Chart: &models.Chart{ID: "chart3", Title: "Profit Data", AxesTitles: "Month vs Profit", Data: "Jan:300,Feb:350", Description: "Monthly profit data"}},
	}

	req, err := http.NewRequest("DELETE", "/favorites/user3/chart3", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/favorites/{userID}/{assetID}", RemoveFavorite)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	if len(UsersFavorites["user3"]) != 0 {
		t.Errorf("handler did not remove favorite correctly: got %v", UsersFavorites["user3"])
	}
}

func TestEditFavorite(t *testing.T) {
	UsersFavorites["user4"] = []models.Asset{
		{Chart: &models.Chart{ID: "chart4", Title: "Expense Data", AxesTitles: "Month vs Expense", Data: "Jan:400,Feb:450", Description: "Monthly expense data"}},
	}

	newAsset := models.Asset{Chart: &models.Chart{ID: "chart4", Title: "Updated Expense Data", AxesTitles: "Month vs Expense", Data: "Jan:500,Feb:550", Description: "Updated monthly expense data"}}
	body, err := json.Marshal(newAsset)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/favorites/user4/chart4", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/favorites/{userID}/{assetID}", EditFavorite)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if UsersFavorites["user4"][0].Chart.Title != "Updated Expense Data" {
		t.Errorf("handler did not edit favorite correctly: got %v", UsersFavorites["user4"][0])
	}
}
