package service

import (
	"gwi_api/internal/domain"
	"gwi_api/internal/repository"
)

type FavoritesService interface {
	AddFavorites(userId uint64, assetId uint64) error
	GetFavorites(userId uint64, params domain.PaginationParams) ([]domain.Asset, error)
}

type favoritesService struct {
	dbHandler repository.DBHandler
}

func NewFavoritesService(dbHandler repository.DBHandler) FavoritesService {
	return &favoritesService{dbHandler}
}

func (f *favoritesService) AddFavorites(userId uint64, assetId uint64) error {
	return f.dbHandler.NewFavoritesQuery().AddFavorites(userId, assetId)
}

func (f *favoritesService) GetFavorites(userId uint64, params domain.PaginationParams) ([]domain.Asset, error) {
	return f.dbHandler.NewFavoritesQuery().GetFavorites(userId, params)
}
