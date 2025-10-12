package handler

import (
	"encoding/json"
	"net/http"

	"platform-go-challenge/internal/api_errors"
	"platform-go-challenge/internal/model"
	"platform-go-challenge/internal/storage"
	"platform-go-challenge/internal/util"
)

// GetFavorites returns all favorites for a user
func GetFavorites(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("id")

		favorites, err := store.GetUserFavorites(userID)
		if err != nil {
			util.RespondError(w, http.StatusInternalServerError, "Failed to retrieve favorites")
			return
		}

		util.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"user_id":   userID,
			"favorites": favorites,
			"count":     len(favorites),
		})
	}
}

// AddFavorite adds a new asset to user's favorites
func AddFavorite(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("id")

		var req struct {
			Asset model.Asset `json:"asset"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.RespondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		baseAsset := req.Asset.GetBaseAsset()
		if baseAsset == nil {
			util.RespondError(w, http.StatusBadRequest, "Asset must have one of: chart, insight, or audience")
			return
		}

		if baseAsset.ID == "" {
			util.RespondError(w, http.StatusBadRequest, "Asset ID is required")
			return
		}

		if err := store.AddFavorite(userID, &req.Asset); err != nil {
			util.RespondError(w, http.StatusInternalServerError, "Failed to add favorite")
			return
		}

		util.RespondJSON(w, http.StatusCreated, map[string]interface{}{
			"message":  "Favorite added successfully",
			"asset_id": baseAsset.ID,
		})
	}
}

// RemoveFavorite removes an asset from user's favorites
func RemoveFavorite(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("id")
		assetID := r.PathValue("assetId")

		if err := store.RemoveFavorite(userID, assetID); err != nil {
			if err == api_errors.ErrUserNotFound || err == api_errors.ErrFavoriteNotFound {
				util.RespondError(w, http.StatusNotFound, "Favorite not found")
				return
			}
			util.RespondError(w, http.StatusInternalServerError, "Failed to remove favorite")
			return
		}

		util.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"message":  "Favorite removed successfully",
			"asset_id": assetID,
		})
	}
}

// EditFavorite updates the description of a favorite
func EditFavorite(store storage.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.PathValue("id")
		assetID := r.PathValue("assetId")

		var req struct {
			Description string `json:"description"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.RespondError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		if err := store.UpdateFavoriteDescription(userID, assetID, req.Description); err != nil {
			if err == api_errors.ErrUserNotFound || err == api_errors.ErrFavoriteNotFound {
				util.RespondError(w, http.StatusNotFound, "Favorite not found")
				return
			}
			util.RespondError(w, http.StatusInternalServerError, "Failed to update favorite")
			return
		}

		util.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"message":  "Favorite updated successfully",
			"asset_id": assetID,
		})
	}
}