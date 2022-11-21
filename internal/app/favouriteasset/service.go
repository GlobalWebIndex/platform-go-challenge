package favouriteasset

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service handles favourite asset information.
type Service interface {
	GetFavouriteAssets(ctx context.Context, userID string) (*GetFavouriteAssetsRes, error)
	AddAssetToFavourites(ctx context.Context, userID, assetID string) (*AddFavouriteAssetRes, error)
	EditFavouriteAsset(ctx context.Context, userID, fAssetID string, fAsset EditFavouriteAsset) error
	RemoveAssetFromFavourites(ctx context.Context, assetID, userID string) error
}

type assetService struct {
	ar Repository
}

// NewService constructor.
func NewService(assetRepo Repository) Service {
	return &assetService{
		ar: assetRepo,
	}
}

func (a assetService) GetFavouriteAssets(ctx context.Context, userID string) (*GetFavouriteAssetsRes, error) {
	pUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	favouriteAssets, err := a.ar.GetFavouriteAssets(ctx, pUserID)
	if err != nil {
		return nil, err
	}

	return &GetFavouriteAssetsRes{
		FavouriteAssets: favouriteAssets,
	}, err
}

func (a assetService) AddAssetToFavourites(ctx context.Context, userID, assetID string) (*AddFavouriteAssetRes, error) {
	pUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	pAssetID, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return nil, err
	}
	res, err := a.ar.AddToFavourites(ctx, pUserID, pAssetID)
	if err != nil {
		return nil, err
	}

	return &AddFavouriteAssetRes{ID: res.Hex()}, nil
}

func (a assetService) EditFavouriteAsset(ctx context.Context, userID, fAssetID string, fAsset EditFavouriteAsset) error {
	pUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	pAssetID, err := primitive.ObjectIDFromHex(fAssetID)
	if err != nil {
		return err
	}

	return a.ar.UpdateFavouriteAsset(ctx, pUserID, pAssetID, fAsset)
}

func (a assetService) RemoveAssetFromFavourites(ctx context.Context, fAssetID, userID string) error {
	pUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	pAssetID, err := primitive.ObjectIDFromHex(fAssetID)
	if err != nil {
		return err
	}

	return a.ar.RemoveAssetFromFavourites(ctx, pAssetID, pUserID)
}
