package queries

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/google/uuid"
)

// GetFavoriteRequest represents the query model
type GetFavoriteRequest struct {
	UserID     uuid.UUID
	FavoriteID uuid.UUID
}

// GetFavoriteResult is the return model of Favorite Query Handlers
type GetFavoriteResult struct {
	ID          uuid.UUID           `json:"id"`
	Type        favourite.AssetType `json:"type"`
	Description string              `json:"description"`
	Data        json.RawMessage     `json:"data"`
	CreatedAt   time.Time           `json:"createdAt"`
}

// GetFavoriteRequestHandler interface for handling the query
type GetFavoriteRequestHandler interface {
	Handle(query GetFavoriteRequest) (*GetFavoriteResult, error)
}

type getFavoriteRequestHandler struct {
	repo favourite.Repository
}

// NewGetFavoriteRequestHandler constructor
func NewGetFavoriteRequestHandler(repo favourite.Repository) GetFavoriteRequestHandler {
	return getFavoriteRequestHandler{repo: repo}
}

// Handle fetches a specific favorite for a given user
func (h getFavoriteRequestHandler) Handle(query GetFavoriteRequest) (*GetFavoriteResult, error) {
	fav, err := h.repo.GetByID(query.UserID, query.FavoriteID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch favorite %s for user %s: %w", query.FavoriteID, query.UserID, err)
	}
	if fav == nil {
		return nil, fmt.Errorf("favorite %s not found for user %s", query.FavoriteID, query.UserID)
	}

	return &GetFavoriteResult{
		ID:          fav.ID,
		Type:        fav.Type,
		Description: fav.Description,
		Data:        fav.Data,
		CreatedAt:   fav.CreatedAt,
	}, nil
}
