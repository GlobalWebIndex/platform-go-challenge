package domain

import (
	"context"
	"errors"
	"platform-go-challenge/internal/service/types"
)

var (
	ErrNotFound         = errors.New("resource not found")
	ErrInvalidAssetType = errors.New("invalid asset type")
	ErrAlreadyExists    = errors.New("resource already exists")
	ErrInvalidInput     = errors.New("invalid input")
)

type Repository interface {
	GetAssets(ctx context.Context, ids []any) ([]types.Asset, error)
	UpdateAsset(ctx context.Context, asset types.Asset) error
	AddFavorite(ctx context.Context, fav types.Favorite) error
	RemoveFavorite(ctx context.Context, userID, assetID string) error
	ListFavorites(ctx context.Context, userID string, page, limit int) ([]types.Favorite, error)
}
