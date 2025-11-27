package commands

import (
	"fmt"

	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/google/uuid"
)

// DeleteFavoriteRequest represents the command to delete a favorite
type DeleteFavoriteRequest struct {
	UserID     uuid.UUID
	FavoriteID uuid.UUID
}

// DeleteFavoriteRequestHandler interface
type DeleteFavoriteRequestHandler interface {
	Handle(command DeleteFavoriteRequest) error
}

type deleteFavoriteRequestHandler struct {
	repo favourite.Repository
}

// NewDeleteFavoriteRequestHandler constructor
func NewDeleteFavoriteRequestHandler(repo favourite.Repository) DeleteFavoriteRequestHandler {
	return deleteFavoriteRequestHandler{repo: repo}
}

// Handle deletes a favorite for a specific user
func (h deleteFavoriteRequestHandler) Handle(command DeleteFavoriteRequest) error {
	// Check if the favorite exists for this user
	fav, err := h.repo.GetByID(command.UserID, command.FavoriteID)
	if err != nil {
		return fmt.Errorf("failed to check favorite existence: %w", err)
	}
	if fav == nil {
		return fmt.Errorf("favorite with ID %s does not exist for user %s", command.FavoriteID, command.UserID)
	}

	// Delete the favorite
	if err := h.repo.Delete(command.UserID, command.FavoriteID); err != nil {
		return fmt.Errorf("failed to delete favorite: %w", err)
	}

	return nil
}
