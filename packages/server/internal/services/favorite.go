package services

import (
	"fmt"

	"github.com/alexandros-georgantas/platform-go-challenge/internal/helpers"
	"github.com/alexandros-georgantas/platform-go-challenge/internal/models"
	"github.com/alexandros-georgantas/platform-go-challenge/internal/serializers"
	"gorm.io/gorm"
)

type FavoriteService interface {
	GetFavorite(uId uint, aId uint) (*models.AssetResponse, error)
	GetFavorites(uId uint) (serializers.FavoritesResponse, error)
	AddToFavorites(uId uint, aId uint) (*models.Favorite, error)
	RemoveFromFavorites(fId uint) (*uint, error)
}

type favoriteService struct {
	db gorm.DB
}

func NewFavoriteService(db gorm.DB) (FavoriteService, error) {
	return &favoriteService{db: db}, nil
}

func (fs *favoriteService) GetFavorite(uId uint, aId uint) (*models.AssetResponse, error) {
	var favorite models.Favorite

	if err := fs.db.Joins("Asset").Where("user_id = ?", uId).First(&favorite, aId).Error; err != nil {
		return nil, err
	}

	result, aErr := helpers.AggregateAssetType(favorite.Asset, fs.db)

	if aErr != nil {
		return nil, aErr
	}

	return result, nil
}

func (fs *favoriteService) GetFavorites(uId uint) (serializers.FavoritesResponse, error) {
	var favorites []models.Favorite
	fr := models.FavoritesResponse{}
	pr := serializers.FavoritesResponse{}
	var count int64
	pr.Favorites = nil
	pr.TotalCount = 0

	// TODO: Here first check if we have items in cache and return them, if not proceed with hitting the DB
	// (take into consideration the pagination params)

	fs.db.Model(&models.Favorite{}).Where("user_id = ?", uId).Count(&count)

	if err := fs.db.Joins("Asset").Where("user_id = ?", uId).Order("asset_id asc").Order("updated_at desc").Find(&favorites).Error; err != nil {
		return pr, err
	}

	// TODO: hydrate the cache

	var favoriteAssets []models.FavoritesResponse
	for _, favorite := range favorites {

		fr.ID = favorite.ID
		fr.AssetID = favorite.AssetID
		fr.UserID = favorite.UserID

		result, err := helpers.AggregateAssetType(favorite.Asset, fs.db)
		if err != nil {
			return pr, err
		}
		fr.Asset = *result

		favoriteAssets = append(favoriteAssets, fr)
	}

	pr.Favorites = &favoriteAssets
	pr.TotalCount = int(count)

	return pr, nil
}

func (fs *favoriteService) AddToFavorites(uId uint, aId uint) (*models.Favorite, error) {
	favorite := models.Favorite{UserID: uId, AssetID: aId}

	result := fs.db.Preload("Asset").Create(&favorite).First(&favorite)

	// TODO: Here invalidation of user's favorites cache should take place

	if result.Error != nil {

		return nil, result.Error
		// for some reason the below snippet is not working
		// if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		// 	return nil, fmt.Errorf("already favorite")
		// } else {
		// 	return nil, fmt.Errorf("something went wrong while adding asset with id %v to favorites", aId)
		// }
	}
	return &favorite, nil
}

func (fs *favoriteService) RemoveFromFavorites(fId uint) (*uint, error) {
	if err := fs.db.Delete(&models.Favorite{}, fId).Error; err != nil {
		return nil, fmt.Errorf("something went wrong while deleting favorite %v", fId)
	}

	// TODO: Here invalidation of user's favorites cache should take place

	return &fId, nil
}
