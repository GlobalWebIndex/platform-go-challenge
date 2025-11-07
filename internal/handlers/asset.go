package handlers

import (
	"assetsApp/internal/models"
	assetServices "assetsApp/internal/services/asset"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AssetHandler struct {
	service *assetServices.AssetService
}

func NewAssetHandler(service *assetServices.AssetService) *AssetHandler {
	return &AssetHandler{service: service}
}

func (h *AssetHandler) GetAssets(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["userId"]
	userID, err := uuid.Parse(userIDStr)

	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	log.Printf("GetAssets called for user %v", userID)
	favs := h.service.GetAssets(userID)
	json.NewEncoder(w).Encode(favs)
	log.Printf("GetAssets completed for user %v", userID)

}

func (h *AssetHandler) AddAsset(w http.ResponseWriter, r *http.Request) {
	userIDStr := mux.Vars(r)["userId"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	// log.Printf("AddAsset called for user %v", userID)

	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	assetType, ok := body["type"].(string)
	if !ok {
		http.Error(w, "Asset type required", http.StatusBadRequest)
		return
	}
	// log.Printf("Adding asset of type %s for user %v", assetType, userID)

	asset, err := models.CreateAsset(assetType, body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.service.AddAsset(userID, asset)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(asset)
}

func (h *AssetHandler) RemoveAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	assetID := vars["assetId"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	log.Printf("RemoveAsset called for user %v, asset %s", userID, assetID)
	if !h.service.RemoveAsset(userID, assetID) {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("RemoveAsset completed for user %v, asset %s", userID, assetID)
}

func (h *AssetHandler) EditAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	assetID := vars["assetId"]
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	log.Printf("EditAsset called for user %v, asset %s", userID, assetID)
	var body struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.service.EditDescription(userID, assetID, body.Description) {
		http.Error(w, "Asset not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.Printf("EditAsset completed for user %v, asset %s", userID, assetID)
}
