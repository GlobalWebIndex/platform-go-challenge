package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/meletios/gwi-engineering-challenge/models"
)

var UsersFavorites = make(map[string][]models.Asset)

func GetFavorites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	favorites, exists := UsersFavorites[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(favorites)
}

func AddFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]

	var asset models.Asset
	err := json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	UsersFavorites[userID] = append(UsersFavorites[userID], asset)
	w.WriteHeader(http.StatusCreated)
}

func RemoveFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	assetID := vars["assetID"]

	favorites, exists := UsersFavorites[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	for i, asset := range favorites {
		if (asset.Chart != nil && asset.Chart.ID == assetID) ||
			(asset.Insight != nil && asset.Insight.ID == assetID) ||
			(asset.Audience != nil && asset.Audience.ID == assetID) {
			UsersFavorites[userID] = append(favorites[:i], favorites[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Asset not found", http.StatusNotFound)
}

func EditFavorite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	assetID := vars["assetID"]

	var newAsset models.Asset
	err := json.NewDecoder(r.Body).Decode(&newAsset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	favorites, exists := UsersFavorites[userID]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	for i, asset := range favorites {
		if (asset.Chart != nil && asset.Chart.ID == assetID) ||
			(asset.Insight != nil && asset.Insight.ID == assetID) ||
			(asset.Audience != nil && asset.Audience.ID == assetID) {
			UsersFavorites[userID][i] = newAsset
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	http.Error(w, "Asset not found", http.StatusNotFound)
}
