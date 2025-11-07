package storage

import (
	"assetsApp/internal/models"
	"github.com/google/uuid"
)


type AssetStore interface {
	Get(userID uuid.UUID) []models.Asset
	Add(userID uuid.UUID, asset models.Asset)
	Remove(userID uuid.UUID, assetID string) bool
	EditDescription(userID uuid.UUID, assetID, newDesc string) bool

	GetFavourites(userID uuid.UUID) []models.Favourite
	AddFavourite(userID uuid.UUID, assetID, assetType string) bool
	RemoveFavourite(userID uuid.UUID, assetID string) bool
}

