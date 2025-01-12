package domain

import (
	"context"

	"platform-go-challenge/internal/service/types"
)

type AssetService interface {
	AddFavorite(ctx context.Context, userID, assetID string) error
	RemoveFavorite(ctx context.Context, userID, assetID string) error
	ListFavorites(ctx context.Context, userID string, page, limit int) ([]types.Asset, error)
	UpdateDescription(ctx context.Context, userID, id string, description string) error
}
