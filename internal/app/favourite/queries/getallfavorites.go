package queries

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/google/uuid"
)

// GetAllFavoritesRequest represents a query to fetch all favorites for a user
type GetAllFavoritesRequest struct {
	UserID uuid.UUID
}

// GetAllFavoritesResult represents the data returned for each favorite
type GetAllFavoritesResult struct {
	ID          uuid.UUID           `json:"id"`
	Type        favourite.AssetType `json:"type"`
	Description string              `json:"description"`
	Data        json.RawMessage     `json:"data"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

// GetAllFavoritesRequestHandler interface
type GetAllFavoritesRequestHandler interface {
	Handle(query GetAllFavoritesRequest) ([]GetAllFavoritesResult, error)
}

type getAllFavoritesRequestHandler struct {
	repo favourite.Repository
}

// NewGetAllFavoritesRequestHandler constructor
func NewGetAllFavoritesRequestHandler(repo favourite.Repository) GetAllFavoritesRequestHandler {
	return getAllFavoritesRequestHandler{repo: repo}
}

// Handle fetches all favorites for a specific user
func (h getAllFavoritesRequestHandler) Handle(query GetAllFavoritesRequest) ([]GetAllFavoritesResult, error) {
	items, err := h.repo.GetAll(query.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch favorites for user %s: %w", query.UserID, err)
	}

	var result []GetAllFavoritesResult
	for _, fav := range items {
		result = append(result, GetAllFavoritesResult{
			ID:          fav.ID,
			Type:        fav.Type,
			Description: fav.Description,
			Data:        fav.Data,
			CreatedAt:   fav.CreatedAt,
			UpdatedAt:   fav.UpdatedAt,
		})
	}

	return result, nil
}
