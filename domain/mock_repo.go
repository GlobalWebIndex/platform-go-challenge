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
func (d *MockDB) FavouriteAsset(ctx context.Context, userID, assetID uint, at AssetType, isFavourite bool) (uint, error) {
	return 0, nil
}
func (d *MockDB) ListFavouriteAssets(ctx context.Context, userID uint, onlyFav bool, query QueryAssets) (*ListedAssets, error) {
	return nil, nil
}
func (d *MockDB) AddUser(ctx context.Context, user User) (*User, error) {
	return nil, nil
}
func (d *MockDB) FindUser(ctx context.Context, username string) (*User, error) {
	return nil, nil
}

func (d *MockDB) GetUser(ctx context.Context, userID uint) (*User, error) {
	return nil, nil
}
