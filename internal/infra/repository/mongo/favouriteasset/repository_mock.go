package favouriteasset

import (
	"context"
	assetsrv "platform-go-challenge/internal/app/favouriteasset"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RepoMock describes a mock struct.
type RepoMock struct {
	mock.Mock
}

// GetFavouriteAssets mock.
func (rm *RepoMock) GetFavouriteAssets(ctx context.Context, userID primitive.ObjectID) (*[]assetsrv.FavouriteAsset, error) {
	args := rm.MethodCalled("GetFavouriteAssets", ctx, userID)

	return args.Get(0).(*[]assetsrv.FavouriteAsset), args.Error(1)
}

// AddToFavourites mock.
func (rm *RepoMock) AddToFavourites(ctx context.Context, userID, assetID primitive.ObjectID) (*primitive.ObjectID, error) {
	args := rm.MethodCalled("AddToFavourites", ctx, userID, assetID)

	return args.Get(0).(*primitive.ObjectID), args.Error(1)
}

// UpdateFavouriteAsset mock.
func (rm *RepoMock) UpdateFavouriteAsset(ctx context.Context, userID, fAssetID primitive.ObjectID, fAsset assetsrv.EditFavouriteAsset) error {
	args := rm.MethodCalled("UpdateFavouriteAsset", ctx, userID, fAssetID, fAsset)

	return args.Error(0)
}

// RemoveAssetFromFavourites mock.
func (rm *RepoMock) RemoveAssetFromFavourites(ctx context.Context, fAssetID, userID primitive.ObjectID) error {
	args := rm.MethodCalled("RemoveAssetFromFavourites", ctx, fAssetID, userID)

	return args.Error(0)
}
