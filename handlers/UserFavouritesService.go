package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"platform-go-challenge/data"
	"platform-go-challenge/models"
)

// UserFavouritesService implements the UserFavouritesHandler interface.
type UserFavouritesService struct {
	dataService data.UserDataService
}

// NewUserFavouritesService creates a new UserFavouritesService.
func NewUserFavouritesService(dataService data.UserDataService) *UserFavouritesService {
	return &UserFavouritesService{
		dataService: dataService,
	}
}

// GetUserFavourites godoc
// @Summary Get user's favourites
// @Description Get the list of user's favourite assets
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {array} models.Asset
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/favourites [get]
func (h *UserFavouritesService) GetUserFavourites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	user, err := h.dataService.GetUser(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(user.Favourites)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// AddAssetToFavourites godoc
// @Summary Add asset to user's favourites
// @Description Add an asset to the user's list of favourite assets
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param asset body models.Asset true "Asset"
// @Success 200 {array} models.Asset
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/favourites [post]
func (h *UserFavouritesService) AddAssetToFavourites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var asset models.Asset
	err := json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = h.dataService.AddAssetToFavourites(userID, asset)
	if err != nil {
		if err.Error() == "duplicate asset ID" {
			http.Error(w, "Asset already exists in favourites", http.StatusConflict)
		} else {
			http.Error(w, "Failed to add asset", http.StatusInternalServerError)
		}
		return
	}

	user, err := h.dataService.GetUser(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(user.Favourites)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// RemoveAssetFromFavourites godoc
// @Summary Remove asset from user's favourites
// @Description Remove an asset from the user's list of favourite assets
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param asset_id path string true "Asset ID"
// @Success 200 {array} models.Asset
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/favourites/{asset_id} [delete]
func (h *UserFavouritesService) RemoveAssetFromFavourites(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	assetID := vars["asset_id"]

	err := h.dataService.RemoveAssetFromFavourites(userID, assetID)
	if err != nil {
		http.Error(w, "Failed to remove asset", http.StatusInternalServerError)
		return
	}

	user, _ := h.dataService.GetUser(userID)
	err = json.NewEncoder(w).Encode(user.Favourites)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// EditAssetDescription godoc
// @Summary Edit asset description
// @Description Edit the description of an asset in the user's list of favourite assets
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Param asset_id path string true "Asset ID"
// @Param asset body models.Asset true "Asset"
// @Success 200 {array} models.Asset
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id}/favourites/{asset_id} [put]
func (h *UserFavouritesService) EditAssetDescription(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	assetID := vars["asset_id"]

	var updatedAsset models.Asset
	err := json.NewDecoder(r.Body).Decode(&updatedAsset)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = h.dataService.EditAssetDescription(userID, assetID, updatedAsset.Description)
	if err != nil {
		http.Error(w, "Failed to edit asset description", http.StatusInternalServerError)
		return
	}

	user, _ := h.dataService.GetUser(userID)
	err = json.NewEncoder(w).Encode(user.Favourites)
	if err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
