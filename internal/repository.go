package gwi

import "context"

// Repository defines the contract that a entity used to store data must implement
type Repository interface {
	ExistAssetInFavs(ctx context.Context, assetid string) bool
	RetrieveFavs(ctx context.Context, userid string) (map[AssetType][]Asseter, error)
	AddAssetToFavs(ctx context.Context, userid string, asset *Asset) (bool, error)
	UpdateFav(ctx context.Context, userid string, asset *Asset) (bool, error)
	RemoveFav(ctx context.Context, userid string, asset *Asset) (bool, error)
}
