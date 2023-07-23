package repository

import (
	//"fmt"

	"gwi_api/internal/domain"
)

type FavoritesQuery interface {
	AddFavorites(userId uint64, assetId uint64) error
	GetFavorites(userId uint64, params domain.PaginationParams) ([]domain.Asset, error)
}

type favoritesQuery struct{}

func (f *favoritesQuery) AddFavorites(userId uint64, assetId uint64) error {
	return DB.AddFavorites(userId, assetId)
}

func (f *favoritesQuery) GetFavorites(userId uint64, params domain.PaginationParams) ([]domain.Asset, error) {

	return DB.GetAllFavorites(userId, params)
}
