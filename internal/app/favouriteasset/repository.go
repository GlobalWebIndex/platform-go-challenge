package favouriteasset

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Repository should be able to manage the favourite assets.
type Repository interface {
	GetFavouriteAssets(ctx context.Context, userID primitive.ObjectID) (*[]FavouriteAsset, error)
	AddToFavourites(ctx context.Context, userID, assetID primitive.ObjectID) (*primitive.ObjectID, error)
	UpdateFavouriteAsset(ctx context.Context, userID, fAssetID primitive.ObjectID, fAsset EditFavouriteAsset) error
	RemoveAssetFromFavourites(ctx context.Context, assetID, userID primitive.ObjectID) error
}
