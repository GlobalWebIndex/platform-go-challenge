package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"platform-go-challenge/internal/common"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/models"
	"platform-go-challenge/internal/storage"
)

type AssetHandler struct {
	cfg       *config.Config
	favorites *storage.FavoritesStore
}

// AssetListResponse wraps a list of assets in a JSON
type AssetListResponse struct {
	Assets []AssetResponse `json:"assets"`
}

type AssetResponse struct {
	Type  string       `json:"type"`
	Asset models.Asset `json:"asset"`
}

func NewAssetHandler(cfg *config.Config, favorites *storage.FavoritesStore) *AssetHandler {
	return &AssetHandler{cfg: cfg, favorites: favorites}
}

func (ah *AssetHandler) list(w http.ResponseWriter, r *http.Request) {
	// Extract authenticated user info from context
	userIDVal := r.Context().Value(ah.cfg.ContextKeys.UserID)
	roleVal := r.Context().Value(ah.cfg.ContextKeys.Role)

	if userIDVal == nil || userIDVal == uuid.Nil || roleVal == nil {
		errorResponse := common.ErrorResponse{
			Error: common.ErrUnauthorized,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized) // 401 Unauthorized
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	authenticatedUserID := userIDVal.(uuid.UUID)
	role := roleVal.(models.Role)

	// Extract user id from query params (if present) - default is the authenticated user
	targetUserID := authenticatedUserID
	userIDStr := r.URL.Query().Get("user_id")

	if userIDStr != "" {
		parsedID, err := uuid.Parse(userIDStr)
		if err != nil {
			errorResponse := common.ErrorResponse{
				Error: common.ErrInvalidUserID,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
			return
		}

		// Normal users may not view other users’ favorites
		if role != models.RoleAdmin && parsedID != authenticatedUserID {
			errorResponse := common.ErrorResponse{
				Error: common.ErrForbiddenAccess,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden) // 403 Forbidden
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
			return
		}

		targetUserID = parsedID
	}

	// Fetch favorites from storage
	assets, err := ah.favorites.GetUserFavorites(targetUserID)
	if err != nil {
		errorResponse := common.ErrorResponse{
			Error: common.ErrFetchingFavorites + ": " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Wrap the assets in AssetListResponse, including the type field dynamically
	assetResponses := make([]AssetResponse, 0, len(assets))
	for _, asset := range assets {
		// Determine the type of asset and create a response object with type information
		var assetType string
		switch asset.(type) {
		case *models.Chart:
			assetType = "chart"
		case *models.Insight:
			assetType = "insight"
		case *models.Audience:
			assetType = "audience"
		}

		assetResponses = append(assetResponses, AssetResponse{
			Type:  assetType,
			Asset: asset, // Asset data
		})
	}

	// Wrap the assets in a list response
	assetListResponse := AssetListResponse{
		Assets: assetResponses,
	}
	// Respond with JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(assetListResponse); err != nil {
		// Create an error response for failure to encode the response
		errorResponse := common.ErrorResponse{
			Error: "failed to encode response: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
	}
}

func (ah *AssetHandler) add(w http.ResponseWriter, r *http.Request) {
	// Extract authenticated user info from context
	userIDVal := r.Context().Value(ah.cfg.ContextKeys.UserID)
	roleVal := r.Context().Value(ah.cfg.ContextKeys.Role)

	if userIDVal == nil || userIDVal == uuid.Nil || roleVal == nil {
		// Create an ErrorResponse and encode it as JSON
		errorResponse := common.ErrorResponse{
			Error: common.ErrUnauthorized, // You can still use the same string from common.ErrUnauthorized
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized) // Set the status code first
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	authenticatedUserID := userIDVal.(uuid.UUID)
	role := roleVal.(models.Role)

	// Parse the request body
	var payload struct {
		UserID    *uuid.UUID      `json:"user_id,omitempty"` // optional
		AssetData json.RawMessage `json:"asset"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		// Invalid JSON in payload
		errorResponse := common.ErrorResponse{
			Error: common.ErrInvalidJSON + ": " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Determine the target user whose favorites we modify
	targetUserID := authenticatedUserID

	if payload.UserID != nil {
		// Normal users can only modify their own favorites
		if role != models.RoleAdmin && *payload.UserID != authenticatedUserID {
			errorResponse := common.ErrorResponse{
				Error: common.ErrForbiddenAccess,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
			return
		}

		targetUserID = *payload.UserID
	}

	// Create the asset from raw json
	asset, err := models.CreateAsset(payload.AssetData)
	if err != nil {
		if err.Error() == common.ErrInvalidAssetType {
			http.Error(w, common.ErrInvalidAssetType, http.StatusBadRequest)
		} else {
			http.Error(w, "unexpected error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Add to favorites store
	if err := ah.favorites.AddFavorite(targetUserID, asset); err != nil {
		errorResponse := common.ErrorResponse{
			Error: common.ErrAddFavorite + ": " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Success
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	successResponse := map[string]string{
		"status": "favorite added",
	}
	if err := json.NewEncoder(w).Encode(successResponse); err != nil {
		http.Error(w, "Failed to encode success response", http.StatusInternalServerError)
	}
}

func (ah *AssetHandler) remove(w http.ResponseWriter, r *http.Request) {
	// Extract authenticated user info from context
	userIDVal := r.Context().Value(ah.cfg.ContextKeys.UserID)
	roleVal := r.Context().Value(ah.cfg.ContextKeys.Role)

	// Check for missing or invalid user or role
	if userIDVal == nil || userIDVal == uuid.Nil || roleVal == nil {
		// Create a structured error response
		errorResponse := common.ErrorResponse{
			Error: common.ErrUnauthorized,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	authenticatedUserID := userIDVal.(uuid.UUID)
	role := roleVal.(models.Role)

	// Determine which user's favorites to modify
	targetUserID := authenticatedUserID
	userIDStr := r.URL.Query().Get("user_id")

	if userIDStr != "" {
		parsedID, err := uuid.Parse(userIDStr)
		if err != nil {
			errorResponse := common.ErrorResponse{
				Error: common.ErrInvalidUserID,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
			return
		}

		// Normal users can only modify their own favorites
		if role != models.RoleAdmin && parsedID != authenticatedUserID {
			errorResponse := common.ErrorResponse{
				Error: common.ErrForbiddenAccess,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
			return
		}

		targetUserID = parsedID
	}

	// Extract asset_id
	assetID := r.URL.Query().Get("asset_id")
	if assetID == "" {
		errorResponse := common.ErrorResponse{
			Error: common.ErrMissingAssetID,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Remove the asset from favorites
	if err := ah.favorites.RemoveFavorite(targetUserID, assetID); err != nil {
		if err.Error() == common.ErrAssetNotFound {
			errorResponse := common.ErrorResponse{
				Error: common.ErrAssetNotFound,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound) // 404 Not Found
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
		} else {

			errorResponse := common.ErrorResponse{
				Error: common.ErrRemoveFavorite + ": " + err.Error(),
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
		}
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK) // 200 OK (or 204 No Content if you prefer)
	successResponse := map[string]string{
		"status": "favorite removed",
	}
	if err := json.NewEncoder(w).Encode(successResponse); err != nil {
		http.Error(w, "Failed to encode success response", http.StatusInternalServerError)
	}
}

func (ah *AssetHandler) update(w http.ResponseWriter, r *http.Request) {
	// Extract authenticated user info from context
	userIDVal := r.Context().Value(ah.cfg.ContextKeys.UserID)
	roleVal := r.Context().Value(ah.cfg.ContextKeys.Role)

	// Check for missing or invalid user or role
	if userIDVal == nil || userIDVal == uuid.Nil || roleVal == nil {
		// Create a structured error response
		errorResponse := common.ErrorResponse{
			Error: common.ErrUnauthorized,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	authenticatedUserID := userIDVal.(uuid.UUID)
	role := roleVal.(models.Role)

	// Parse payload
	var payload struct {
		UserID    *uuid.UUID      `json:"user_id,omitempty"` // optional
		AssetData json.RawMessage `json:"asset"`             // must include asset ID
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		// Invalid JSON in payload
		errorResponse := common.ErrorResponse{
			Error: common.ErrInvalidJSON + ": " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Determine the target user whose favorites we modify
	targetUserID := authenticatedUserID

	if payload.UserID != nil {
		// Normal users can only modify their own favorites
		if role != models.RoleAdmin && *payload.UserID != authenticatedUserID {
			errorResponse := common.ErrorResponse{
				Error: common.ErrForbiddenAccess,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
				http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
			}
			return
		}
		targetUserID = *payload.UserID
	}

	// Create partial asset from raw JSON
	asset, err := models.CreateAsset(payload.AssetData)
	if err != nil {
		errorResponse := common.ErrorResponse{
			Error: common.ErrInvalidAssetType + ": " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	if asset.GetID() == "" {
		errorResponse := common.ErrorResponse{
			Error: common.ErrAssetIDRequired,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Update the favorite
	if err := ah.favorites.UpdateFavorite(targetUserID, asset); err != nil {
		errorResponse := common.ErrorResponse{
			Error: common.ErrUpdateFavorite + ": " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "favorite updated",
	})
}
