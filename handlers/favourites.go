package handlers

import (
	"encoding/json"
	"go-challenge/repository"
	"net/http"
	"strconv"
)

type FavouriteHandler struct {
	repo repository.FavouriteRepository
}

func NewFavouriteHandler(repo repository.FavouriteRepository) *FavouriteHandler {
	return &FavouriteHandler{repo: repo}
}

func (handler *FavouriteHandler) GetUserFavourites(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	if userID == "" {
		http.Error(w, "user ID is required", http.StatusBadRequest)
		return
	}

	query := r.URL.Query()
	limitStr := query.Get("limit")
	pageStr := query.Get("page")

	// Set defaults and convert
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // Default to page 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 20 {
		limit = 10 // Default to 10 items per page
	}

	// Calculate the Offset
	offset := (page - 1) * limit

	favourites, totalCount, err := handler.repo.GetUserFavourites(userID, limit, offset)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return JSON res
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"user_id":     userID,
		"favorites":   favourites,
		"page":        page,
		"limit":       limit,
		"total_count": totalCount,
	})
}

func (handler *FavouriteHandler) AddFavourite(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")

	var req struct {
		AssetID string `json:"asset_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	favourite, err := handler.repo.AddFavourite(userID, req.AssetID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(favourite)
}

func (handler *FavouriteHandler) RemoveFavourite(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userID")
	assetID := r.PathValue("assetID")

	// remove the favourite
	if err := handler.repo.RemoveFavourite(userID, assetID); err != nil {
		// If fails, return error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Return 200
	w.WriteHeader(http.StatusOK)

	// w.WriteHeader(http.StatusNoContent)

	// Create and send success response
	json.NewEncoder(w).Encode(map[string]any{
		"message":  "Favourite removed successfully",
		"user_id":  userID,
		"asset_id": assetID,
	})

}

func (h *FavouriteHandler) UpdateAssetDescription(w http.ResponseWriter, r *http.Request) {
	assetID := r.PathValue("assetID")

	if assetID == "" {
		http.Error(w, "asset ID is required", http.StatusBadRequest)
		return
	}

	var req struct {
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Description) > 100 {
		http.Error(w, "description too long (max 100 characters)", http.StatusBadRequest)
		return
	}

	// Update the asset's description
	updatedAsset, err := h.repo.UpdateAssetDescription(assetID, req.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Asset description updated successfully",
		"asset":   updatedAsset,
	})
}
