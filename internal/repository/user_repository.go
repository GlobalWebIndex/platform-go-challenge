package repository

import (
	"github.com/ceciivanov/platform-go-challenge/internal/models"
)

// UserRepository defines the methods that any type of user repository must implement
type UserRepository interface {
	GetUserFavorites(userID int) (map[int]models.Asset, error)
	AddUserFavorite(userID int, asset models.Asset) error
	DeleteUserFavorite(userID, assetID int) error
	EditUserFavorite(userID int, assetID int, asset models.Asset) error
}
