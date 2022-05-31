package domain

import "context"

func NewDomain(db IDBRepository) *Domain {
	return &Domain{
		repo: db,
	}
}

func (d *Domain) AddAsset(ctx context.Context, asset Asset) error {
	return nil
}

func (d *Domain) DeleteAsset(ctx context.Context, assetID string) error {
	return nil
}

func (d *Domain) UpdateAsset(ctx context.Context, assetID string, asset Asset) error {
	return nil
}
func (d *Domain) ListAssets(ctx context.Context, userID string, query QueryAssets) error {
	return nil
}
func (d *Domain) FavourAnAsset(ctx context.Context, userID, assetID string) error {
	return nil
}
func (d *Domain) CreateUser(ctx context.Context, user User) error {
	return nil
}
func (d *Domain) LoginUser(ctx context.Context, cred LoginCredentials) error {
	return nil
}
