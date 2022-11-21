package favouriteasset

import (
	"context"

	"github.com/stretchr/testify/mock"
)

// SvcMock describes a mock struct.
type SvcMock struct {
	mock.Mock
}

// GetFavouriteAssets mock.
func (m *SvcMock) GetFavouriteAssets(ctx context.Context, userID string) (*GetFavouriteAssetsRes, error) {
	args := m.MethodCalled("GetFavouriteAssets", ctx, userID)

	return args.Get(0).(*GetFavouriteAssetsRes), args.Error(1)
}

// AddAssetToFavourites mock.
func (m *SvcMock) AddAssetToFavourites(ctx context.Context, userID, assetID string) (*AddFavouriteAssetRes, error) {
	args := m.MethodCalled("AddAssetToFavourites", ctx, userID, assetID)

	return args.Get(0).(*AddFavouriteAssetRes), args.Error(1)
}

// EditFavouriteAsset mock.
func (m *SvcMock) EditFavouriteAsset(ctx context.Context, userID, fAssetID string, fAsset EditFavouriteAsset) error {
	args := m.MethodCalled("EditFavouriteAsset", ctx, userID, fAssetID, fAsset)

	return args.Error(0)
}

// RemoveAssetFromFavourites mock.
func (m *SvcMock) RemoveAssetFromFavourites(ctx context.Context, assetID, userID string) error {
	args := m.MethodCalled("RemoveAssetFromFavourites", ctx, assetID, userID)

	return args.Error(0)
}
