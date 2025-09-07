package repository

import (
	"platform-go-challenge/internal/domain/entities"
)

type Repository interface {
	GetFavourites(userID string) ([]entities.Asset, error)
	AddFavourite(userID string, assetID string) error
	RemoveFavourite(userID string, assetID string) error
	UpdateFavouriteDescription(userID string, assetID string, newDescription string) error

	GetUser(userID string) *entities.User
	SaveUser(user *entities.User)

	GetAsset(assetID string) *entities.Asset
	SaveAsset(asset entities.Asset) error
}
