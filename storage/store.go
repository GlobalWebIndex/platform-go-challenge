package storage

import (
	"context"

	"github.com/arhsxro/platform-go-challenge/models"
)

// Signatures of the operations that can be perfomred on the db
type Store interface {
	GetUserFavorites(ctx context.Context, userID, filterType string, page, pageSize int) ([]models.Asset, error)
	AddFavorite(ctx context.Context, userID string, asset models.Asset) error
	RemoveFavorite(ctx context.Context, userID, assetID string) error
	UpdateDescription(ctx context.Context, userID, assetID, newDescription string) error
	Close() error
}
