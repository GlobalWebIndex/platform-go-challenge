package favouriteasset

import (
	"context"
	"encoding/json"
	"errors"
	assetsrv "platform-go-challenge/internal/app/favouriteasset"
	"time"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CacheTTL contains the time for the result to get cached.
const CacheTTL = 30 * time.Second

// CRepo is the favouriteasset repo with cache functionality.
type CRepo struct {
	repo        Repo
	redisClient *redis.Client
}

// NewCacheRepository constructor.
func NewCacheRepository(repo Repo, redisClient *redis.Client) (*CRepo, error) {
	return &CRepo{
		repo:        repo,
		redisClient: redisClient,
	}, nil
}

// GetFavouriteAssets gets a list of user's favourite assets with cache.
func (r *CRepo) GetFavouriteAssets(ctx context.Context, userID primitive.ObjectID) (*[]assetsrv.FavouriteAsset, error) {
	favouriteAssets := &[]assetsrv.FavouriteAsset{}

	// checks cache for cached results
	results, err := r.redisClient.Get(ctx, "GetFavouriteAssets:"+userID.Hex()).Result()
	if err == redis.Nil {
		// cache key not found, should call repo & cache the result.
		// calls repo.
		favouriteAssets, err = r.repo.GetFavouriteAssets(ctx, userID)
		if err != nil {
			return nil, err

		}

		// converts repo result to json string.
		strFavouriteAssets, err := json.Marshal(favouriteAssets)
		if err != nil {
			return nil, errors.New("failed to marshal favouriteAssets" + err.Error())

		}

		// caches the result.
		r.redisClient.Set(ctx, "GetFavouriteAssets:"+userID.Hex(), strFavouriteAssets, CacheTTL)
		if err != nil {
			return nil, errors.New("Failed saving FavouriteAssets in cache: " + err.Error())
		}

		// returns the repo assets.
		return favouriteAssets, err
	} else if err != nil {

		// otherwise returns an error if err found
		return nil, err
	}

	err = json.Unmarshal([]byte(results), &favouriteAssets)
	if err != nil {
		return nil, err
	}

	// serves the cached result.
	return favouriteAssets, err
}

// AddToFavourites flushes cache & calls repo in order to add asset to user favourites.
func (r *CRepo) AddToFavourites(ctx context.Context, userID, assetID primitive.ObjectID) (*primitive.ObjectID, error) {
	go r.redisClient.FlushDBAsync(ctx)

	return r.repo.AddToFavourites(ctx, userID, assetID)
}

// UpdateFavouriteAsset flushes cache & calls repo in order to edit asset from user favourites.
func (r *CRepo) UpdateFavouriteAsset(ctx context.Context, userID, fAssetID primitive.ObjectID, fAsset assetsrv.EditFavouriteAsset) error {
	go r.redisClient.FlushDBAsync(ctx)

	return r.repo.UpdateFavouriteAsset(ctx, userID, fAssetID, fAsset)
}

// RemoveAssetFromFavourites flushes cache & calls repo in order to remove asset from user favourites.
func (r *CRepo) RemoveAssetFromFavourites(ctx context.Context, fAssetID, userID primitive.ObjectID) error {
	go r.redisClient.FlushDBAsync(ctx)

	return r.repo.RemoveAssetFromFavourites(ctx, fAssetID, userID)
}
