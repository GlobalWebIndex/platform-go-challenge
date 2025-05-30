package presentation

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"pgc/domain"
)

const (
	// pathParamAssetID is the URL path parameter key for the asset ID.
	pathParamAssetID = "aid"
	// pathParamID is the URL path parameter key for the favorite ID.
	pathParamID = "id"
)

// FavoriteHandler represents an HTTP handler for managing favorite-related requests.
// It provides methods to handle CRUD operations for user favorites.
type FavoriteHandler struct {
	favoriteService domain.FavoriteService // Service for managing favorite assets.
}

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Error string `json:"error"` // Error message describing the issue.
}

// FavoriteRequest represents the request body for adding or updating a favorite.
// It contains an optional description of the favorite asset.
type FavoriteRequest struct {
	Description string `json:"description,omitempty"` // Optional description of the favorite asset.
}

// GetUserFavorites handles the retrieval of a user's favorite assets.
// It validates the user ID from the request header and fetches the favorites using the service.
func (h FavoriteHandler) GetUserFavorites(writer http.ResponseWriter, request *http.Request) {
	userIDStr := request.Header.Get("X-User-ID")
	if userIDStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "X-User-ID header is required",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "invalid user ID format",
		})
		return
	}

	favorites, err := h.favoriteService.GetUserFavorites(request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidUserID):
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		}
		return
	}

	if len(favorites) == 0 {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(favorites)
}

// AddAssetToUserFavorites handles the addition of an asset to a user's favorites.
// It validates the user ID and asset ID, parses the request body, and adds the favorite using the service.
func (h FavoriteHandler) AddAssetToUserFavorites(writer http.ResponseWriter, request *http.Request) {
	userIDStr := request.Header.Get("X-User-ID")
	if userIDStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "X-User-ID header is required",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "invalid user ID format",
		})
		return
	}

	assetIDStr := request.PathValue(pathParamAssetID)
	if assetIDStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "asset ID is required in URL path",
		})
		return
	}

	assetID, err := strconv.Atoi(assetIDStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "invalid asset ID format",
		})
		return
	}

	var req FavoriteRequest
	if request.Body != nil {
		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: "invalid json body",
			})
			return
		}
	}

	favorite := domain.Favorite{
		UserID:      userID,
		AssetID:     assetID,
		Description: req.Description,
	}

	createdFavorite, err := h.favoriteService.Add(request.Context(), favorite)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidUserID):
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		case errors.Is(err, domain.ErrInvalidAssetID):
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		case errors.Is(err, domain.ErrDuplicateFavorite):
			writer.WriteHeader(http.StatusConflict)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		}

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(createdFavorite)
	return
}

// RemoveAssetFromUserFavorites handles the removal of an asset from a user's favorites.
// It validates the user ID and favorite ID, and removes the favorite using the service.
func (h FavoriteHandler) RemoveAssetFromUserFavorites(writer http.ResponseWriter, request *http.Request) {
	userIDStr := request.Header.Get("X-User-ID")
	if userIDStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "X-User-ID header is required",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "invalid user ID format",
		})
		return
	}

	IDstr := request.PathValue(pathParamID)
	if IDstr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "asset ID is required in URL path",
		})
		return
	}

	ID, err := strconv.Atoi(IDstr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "invalid asset ID format",
		})
		return
	}

	err = h.favoriteService.Delete(context.Background(), ID, userID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidUserID):
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		case errors.Is(err, domain.ErrInvalidAssetID):
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		case errors.Is(err, domain.ErrFavoriteNotFound):
			writer.WriteHeader(http.StatusNotFound)
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		}
		return
	}

	writer.WriteHeader(http.StatusNoContent)
	return

}

// UpdateUserFavoriteAsset handles the update of a user's favorite asset.
// It validates the user ID and favorite ID, parses the request body, and updates the favorite using the service.
func (h FavoriteHandler) UpdateUserFavoriteAsset(writer http.ResponseWriter, request *http.Request) {
	userIDStr := request.Header.Get("X-User-ID")
	if userIDStr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "X-User-ID header is required",
		})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "invalid user ID format",
		})
		return
	}

	IDstr := request.PathValue(pathParamID)
	if IDstr == "" {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "asset ID is required in URL path",
		})
		return
	}

	ID, err := strconv.Atoi(IDstr)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(ErrorResponse{
			Error: "invalid asset ID format",
		})
		return
	}

	var req FavoriteRequest
	if request.Body != nil {
		if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: "invalid json body",
			})
			return
		}
	}

	editedFavorite, err := h.favoriteService.Edit(context.Background(), ID, userID, req.Description)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidUserID):
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		case errors.Is(err, domain.ErrInvalidFavoriteID):
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		case errors.Is(err, domain.ErrFavoriteNotFound):
			writer.WriteHeader(http.StatusNotFound)
		default:
			writer.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(writer).Encode(ErrorResponse{
				Error: err.Error(),
			})
		}
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(editedFavorite)
}

// NewFavoriteHandler creates and returns a new instance of the FavoriteHandler.
// This function is intended to initialize and configure the favorite-related HTTP handlers.
func NewFavoriteHandler(favoriteService domain.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{
		favoriteService: favoriteService,
	}
}
