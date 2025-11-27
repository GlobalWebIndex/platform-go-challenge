package commands

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akazantzidis/gwi-ass/internal/app/notification"
	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/google/uuid"
)

// AddFavoriteRequest model for creating a favorite
type AddFavoriteRequest struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Type        favourite.AssetType
	Description string
	Data        json.RawMessage
}

// CreateFavoriteRequestHandler interface for handling add favorite
type CreateFavoriteRequestHandler interface {
	Handle(command AddFavoriteRequest) error
}

type addFavoriteRequestHandler struct {
	repo                favourite.Repository
	notificationService notification.Service
}

// NewAddFavoriteRequestHandler constructor
func NewAddFavoriteRequestHandler(
	repo favourite.Repository,
	notificationService notification.Service,
) CreateFavoriteRequestHandler {
	return addFavoriteRequestHandler{
		repo:                repo,
		notificationService: notificationService,
	}
}

// Handle adds a new favorite
func (h addFavoriteRequestHandler) Handle(req AddFavoriteRequest) error {
	fav := favourite.Favorite{
		ID:          req.ID,
		Type:        req.Type,
		Description: req.Description,
		Data:        req.Data,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// Store under the correct user
	if err := h.repo.Add(req.UserID, fav); err != nil {
		return fmt.Errorf("failed to add favorite: %w", err)
	}

	// Send notification
	n := notification.Notification{
		Subject: "New Favorite added",
		Message: fmt.Sprintf(
			"A new favorite with description '%s' was added for user %s",
			fav.Description,
			req.UserID.String(),
		),
	}

	if err := h.notificationService.Notify(n); err != nil {
		return fmt.Errorf("favorite added but failed to send notification: %w", err)
	}

	return nil
}
