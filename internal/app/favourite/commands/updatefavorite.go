package commands

import (
	"encoding/json"
	"fmt"

	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/google/uuid"
)

// UpdateFavoriteRequest represents the full update command for a favorite
type UpdateFavoriteRequest struct {
	UserID      uuid.UUID
	ID          uuid.UUID
	Type        favourite.AssetType
	Description string
	Data        json.RawMessage
}

// UpdateFavoriteRequestHandler interface
type UpdateFavoriteRequestHandler interface {
	Handle(command UpdateFavoriteRequest) error
}

type updateFavoriteRequestHandler struct {
	repo favourite.Repository
}

// NewUpdateFavoriteRequestHandler constructor
func NewUpdateFavoriteRequestHandler(repo favourite.Repository) UpdateFavoriteRequestHandler {
	return updateFavoriteRequestHandler{repo: repo}
}

// Handle updates a favorite for a specific user
func (h updateFavoriteRequestHandler) Handle(command UpdateFavoriteRequest) error {
	// Fetch existing favorite for the user
	favorite, err := h.repo.GetByID(command.UserID, command.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch favorite: %w", err)
	}
	if favorite == nil {
		return fmt.Errorf("favorite with ID %s does not exist for user %s", command.ID, command.UserID)
	}

	// Update fields
	favorite.Type = command.Type
	favorite.Description = command.Description
	favorite.Data = command.Data

	// Persist the update
	if err := h.repo.Update(command.UserID, *favorite); err != nil {
		return fmt.Errorf("failed to update favorite: %w", err)
	}

	return nil
}
