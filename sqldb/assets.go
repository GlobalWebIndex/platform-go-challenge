package sqldb

import (
	"context"
	"platform-go-challenge/domain"
)

func (d *DB) AddAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error) {
	return nil, nil
}
func (d *DB) DeleteAsset(ctx context.Context, assetID string) error {
	return nil
}
func (d *DB) UpdateAsset(ctx context.Context, asset domain.Asset) (*domain.Asset, error) {
	return nil, nil
}
func (d *DB) ListAssets(ctx context.Context, userID string, query domain.QueryAssets) error {
	return nil
}
func (d *DB) FavourAnAsset(ctx context.Context, userID, assetID string) (string, error) {
	return "", nil
}
