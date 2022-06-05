package domain

import "context"

type MockDB struct{}

func (d *MockDB) AddAsset(ctx context.Context, asset Asset) (*Asset, error) {
	return nil, nil
}
func (d *MockDB) DeleteAsset(ctx context.Context, assetID uint) error {
	return nil
}
func (d *MockDB) UpdateAsset(ctx context.Context, asset Asset) (*Asset, error) {
	return nil, nil
}
func (d *MockDB) ListAssets(ctx context.Context, userID uint, query QueryAssets) error {
	return nil
}
func (d *MockDB) FavourAnAsset(ctx context.Context, userID, assetID uint) (uint, error) {
	return 0, nil
}

func (d *MockDB) CreateUser(ctx context.Context, user User) (*User, error) {
	return nil, nil
}
func (d *MockDB) FindUser(ctx context.Context, cred LoginCredentials) (*User, error) {
	return nil, nil
}
