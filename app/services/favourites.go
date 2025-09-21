package services

import (
	"platform-go-challenge/domain"
)

type FavouriteService struct {
	fr domain.FavouriteRepository
	ur domain.UserRepository
	ar domain.AssetRepository
}

func NewFavouriteService(fr domain.FavouriteRepository, ur domain.UserRepository, ar domain.AssetRepository) *FavouriteService {
	return &FavouriteService{fr: fr, ur: ur, ar: ar}
}

func (fs FavouriteService) AddToFavourites(userID, assetID string) (domain.Favourite, error) {
	userArangoID, err := fs.ur.GetUser(userID)
	if err != nil {
		return domain.Favourite{}, err
	}

	assetArangoID, err := fs.ar.GetAsset(assetID)
	if err != nil {
		return domain.Favourite{}, err
	}

	favourite, err := fs.fr.AddToFavourites(userArangoID, assetArangoID)
	if err != nil {
		return domain.Favourite{}, err
	}
	return favourite, nil
}

func (fs FavouriteService) GetFavourites(userID string) ([]domain.Asset, error) {
	userArangoID, err := fs.ur.GetUser(userID)
	if err != nil {
		return []domain.Asset{}, err
	}

	assets, err := fs.fr.GetFavourites(userArangoID)
	if err != nil {
		return []domain.Asset{}, err
	}
	return assets, nil
}

func (fs FavouriteService) EditFavourites(userID, assetID, description string) (domain.Asset, error) {
	userArangoID, err := fs.ur.GetUser(userID)
	if err != nil {
		return domain.Asset{}, err
	}

	assetArangoID, err := fs.ar.GetAsset(assetID)
	if err != nil {
		return domain.Asset{}, err
	}

	asset, err := fs.fr.EditFavourites(userArangoID, assetArangoID, description)
	if err != nil {
		return domain.Asset{}, err
	}
	return asset, nil
}

func (fs FavouriteService) DeleteFavourite(userID, assetID string) (domain.Favourite, error) {
	userArangoID, err := fs.ur.GetUser(userID)
	if err != nil {
		return domain.Favourite{}, err
	}

	assetArangoID, err := fs.ar.GetAsset(assetID)
	if err != nil {
		return domain.Favourite{}, err
	}

	favourite, err := fs.fr.DeleteFavourite(userArangoID, assetArangoID)
	if err != nil {
		return domain.Favourite{}, err
	}
	return favourite, nil
}
