package application

import (
	"context"

	"github.com/manos/favourites/favourites/domain"
)

type FavouriteRepository interface {
	FindByUser(ctx context.Context, userID string) ([]domain.FavouriteEntity, error)
	FindByUserPage(ctx context.Context, userID string, offset int, limit int) ([]domain.FavouriteEntity, int, error)
	Save(ctx context.Context, fav domain.FavouriteEntity) error
	FindByID(ctx context.Context, favouriteID string) (domain.FavouriteEntity, error)
	Delete(ctx context.Context, favouriteID string) error
	UpdateDescription(ctx context.Context, favouriteID string, description string) error
}

// Outbound port: Favourites depends on assets through this interface.
// Return favourites-owned DTOs (anti-corruption layer).
type AssetClient interface {
	GetInsight(ctx context.Context, userId string, assetId string) (AssetDTO, error)
	GetAudience(ctx context.Context, userId string, assetId string) (AssetDTO, error)
	GetChart(ctx context.Context, userId string, assetId string) (AssetDTO, error)
	GetInsights(ctx context.Context, userId string, assetIds []string) (map[string]AssetDTO, error)
	GetAudiences(ctx context.Context, userId string, assetIds []string) (map[string]AssetDTO, error)
	GetCharts(ctx context.Context, userId string, assetIds []string) (map[string]AssetDTO, error)
}
