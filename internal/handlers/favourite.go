package handlers

import (
	"assetsApp/internal/models"
	favouriteServices "assetsApp/internal/services/favourite"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type FavouriteHandler struct {
	service *favouriteServices.FavouriteService
}

func NewFavouriteHandler(service *favouriteServices.FavouriteService) *FavouriteHandler {
	return &FavouriteHandler{service: service}
}

func (h *FavouriteHandler) GetFavourites(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["userId"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	favs := h.service.GetFavourites(userID)
	if favs == nil {
		json.NewEncoder(w).Encode([]models.Favourite{})
		return
	}
	json.NewEncoder(w).Encode(favs)
}

func (h *FavouriteHandler) AddFavourite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	assetID := vars["assetId"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var body struct {
		AssetType string `json:"asset_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.service.AddFavourite(userID, assetID, body.AssetType) {
		http.Error(w, "Could not add favourite", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *FavouriteHandler) RemoveFavourite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	assetID := vars["assetId"]

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	if !h.service.RemoveFavourite(userID, assetID) {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
