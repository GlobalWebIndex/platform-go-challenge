package commands

import (
	"encoding/json"
	"fmt"

	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/google/uuid"
)

// PatchFavoriteRequest represents fields for partial updates
type PatchFavoriteRequest struct {
	Type        *favourite.AssetType `json:"type,omitempty"`
	Description *string              `json:"description,omitempty"`
	Data        *json.RawMessage     `json:"data,omitempty"`
}

// UpdatePartialFavoriteRequestHandler interface for PATCH
type UpdatePartialFavoriteRequestHandler interface {
	HandlePartial(userID uuid.UUID, favoriteID uuid.UUID, req PatchFavoriteRequest) (*favourite.Favorite, error)
}

type updatePartialFavoriteRequestHandler struct {
	repo favourite.Repository
}

// NewUpdatePartialFavoriteRequestHandler constructor
func NewUpdatePartialFavoriteRequestHandler(repo favourite.Repository) UpdatePartialFavoriteRequestHandler {
	return &updatePartialFavoriteRequestHandler{repo: repo}
}

// HandlePartial applies only the provided fields to an existing favorite
func (h *updatePartialFavoriteRequestHandler) HandlePartial(userID uuid.UUID, favoriteID uuid.UUID, req PatchFavoriteRequest) (*favourite.Favorite, error) {
	// Fetch favorite for the user
	fav, err := h.repo.GetByID(userID, favoriteID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch favorite: %w", err)
	}
	if fav == nil {
		return nil, fmt.Errorf("favorite with ID %s not found for user %s", favoriteID, userID)
	}

	// Apply only provided updates
	if req.Type != nil {
		fav.Type = *req.Type
	}
	if req.Description != nil {
		fav.Description = *req.Description
	}
	if req.Data != nil {
		fav.Data = *req.Data
	}

	// Persist the update
	if err := h.repo.Update(userID, *fav); err != nil {
		return nil, fmt.Errorf("failed to update favorite: %w", err)
	}

	return fav, nil
}
