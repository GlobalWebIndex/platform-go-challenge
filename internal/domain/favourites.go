package domain

import (
	"time"

	"github.com/google/uuid"
)

// Favourite represents "user favourites an asset".
type Favourite struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	AssetID   uuid.UUID
	CreatedAt time.Time
}

// NewFavourite creates a new favourite with a timestamp.
func NewFavourite(userID, assetID uuid.UUID) Favourite {
	return Favourite{
		UserID:    userID,
		AssetID:   assetID,
		CreatedAt: time.Now().UTC(),
	}
}
