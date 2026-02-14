package application

import (
	"errors"

	"github.com/gwi/platform-go-challenge/favourites/domain"
)

// ErrAssetNotFound is returned when adding a favourite with a source_asset_id that does not exist in the asset catalog.
var ErrAssetNotFound = errors.New("asset not found in catalog")

// AssetCatalogClient is the outbound port for the asset catalog. Implemented by adapters.
// Favourites uses it to resolve source_asset_id to full asset data (add-by-reference); lookup is user-scoped.
type AssetCatalogClient interface {
	// GetAsset returns the catalog asset for the given user, source ID and type, or ErrAssetNotFound when not found.
	GetAsset(userID string, sourceAssetID string, assetType string) (*domain.Asset, error)
}

// FavouritesRepository is the outbound port for persistence. Implemented by adapters (e.g. memory, PostgreSQL).
type FavouritesRepository interface {
	ListPaginated(userID string, limit, offset int) (items []domain.Asset, total int)
	Get(userID string, assetID string) (*domain.Asset, bool)
	Add(userID string, asset domain.Asset) domain.Asset
	Remove(userID string, assetID string) bool
	UpdateDescription(userID string, assetID string, description string) (*domain.Asset, bool)
	RemoveFavouritesBySourceAssetID(sourceAssetID string) int
}

// FavouritesUseCase is the inbound port. Consumed by HTTP handlers; implemented by FavouritesService.
type FavouritesUseCase interface {
	ListFavouritesPaginated(userID string, limit, offset int) (items []domain.Asset, total int)
	GetFavourite(userID string, assetID string) (*domain.Asset, bool)
	AddFavourite(userID string, asset domain.Asset) (domain.Asset, error)
	RemoveFavourite(userID string, assetID string) bool
	UpdateFavouriteDescription(userID string, assetID string, description string) (*domain.Asset, bool)
	OnAssetDeleted(sourceAssetID string) (removedCount int)
}
