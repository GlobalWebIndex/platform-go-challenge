package favouriteasset

import (
	"context"
	assetsrv "platform-go-challenge/internal/app/favouriteasset"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository definition.
type Repository struct {
	DB *mongo.Database
}

// Repo contains the interfaces of the repo.
type Repo interface {
	GetFavouriteAssets(ctx context.Context, userID primitive.ObjectID) (*[]assetsrv.FavouriteAsset, error)
	AddToFavourites(ctx context.Context, userID, assetID primitive.ObjectID) (*primitive.ObjectID, error)
	UpdateFavouriteAsset(ctx context.Context, userID, fAssetID primitive.ObjectID, fAsset assetsrv.EditFavouriteAsset) error
	RemoveAssetFromFavourites(ctx context.Context, fAssetID, userID primitive.ObjectID) error
}

// NewRepository constructor.
func NewRepository(db *mongo.Database) (*Repository, error) {
	return &Repository{
		DB: db,
	}, nil
}

// GetFavouriteAssets gets a list of user's favourite assets.
func (r *Repository) GetFavouriteAssets(ctx context.Context, userID primitive.ObjectID) (*[]assetsrv.FavouriteAsset, error) {
	filter := bson.A{
		bson.M{
			"$lookup": bson.M{
				"from":         "assets",
				"localField":   "asset_id",
				"foreignField": "_id",
				"as":           "user_favourite_assets",
			}},
		bson.M{
			"$match": bson.M{
				"user_id": userID,
			}},
	}
	cursor, err := r.DB.Collection("favourite_assets").Aggregate(ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []repoFavouriteAsset
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	var favouriteAssets []assetsrv.FavouriteAsset
	for _, result := range results {
		favouriteAsset, err := result.adaptToModel()
		if err != nil {
			return nil, err
		}
		favouriteAssets = append(favouriteAssets, *favouriteAsset)
	}

	return &favouriteAssets, err
}

// AddToFavourites adds the asset to user's favourites.
func (r *Repository) AddToFavourites(ctx context.Context, userID, assetID primitive.ObjectID) (*primitive.ObjectID, error) {
	favouriteAsset := bson.D{
		{"user_id", userID},
		{"asset_id", assetID},
		{"description", ""},
	}

	c, err := r.DB.Collection("favourite_assets").InsertOne(ctx, favouriteAsset)
	if err != nil {
		return nil, err
	}

	lastInsertID := c.InsertedID.(primitive.ObjectID)

	return &lastInsertID, nil
}

// UpdateFavouriteAsset updates a user favourite asset.
func (r *Repository) UpdateFavouriteAsset(ctx context.Context, userID, fAssetID primitive.ObjectID, fAsset assetsrv.EditFavouriteAsset) error {
	filter := bson.M{
		"_id":     fAssetID,
		"user_id": userID,
	}
	favouriteAsset := bson.D{
		{"$set", bson.D{{"description", fAsset.Description}}},
	}

	_, err := r.DB.Collection("favourite_assets").UpdateOne(ctx, filter, favouriteAsset)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAssetFromFavourites removes a user asset from favourites.
func (r *Repository) RemoveAssetFromFavourites(ctx context.Context, fAssetID, userID primitive.ObjectID) error {

	filter := bson.M{
		"_id":     fAssetID,
		"user_id": userID,
	}

	_, err := r.DB.Collection("favourite_assets").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
