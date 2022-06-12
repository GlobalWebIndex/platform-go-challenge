package domain

import "context"

type MockDB struct{}

func (d *MockDB) AddAsset(ctx context.Context, asset Asset) (*Asset, error) {
	return nil, nil
}
func (d *MockDB) DeleteAsset(ctx context.Context, at AssetType, assetID uint) error {
	return nil
}
func (d *MockDB) UpdateAsset(ctx context.Context, asset Asset) (*Asset, error) {
	return nil, nil
}
func (d *MockDB) GetAsset(ctx context.Context, at AssetType, assetID uint) (*Asset, error) {
	return nil, nil
}
func (d *MockDB) ListAssets(ctx context.Context, query QueryAssets) (*ListedAssets, error) {
	return nil, nil
}
func (d *MockDB) FavouriteAsset(ctx context.Context, userID, assetID uint, isFavourite bool) (uint, error) {
	return 0, nil
}

func (d *MockDB) AddUser(ctx context.Context, user User) (*User, error) {
	return nil, nil
}
func (d *MockDB) FindUser(ctx context.Context, cred LoginCredentials) (*User, error) {
	return nil, nil
}
