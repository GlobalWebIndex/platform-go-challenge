package service

import (
	"github.com/ceciivanov/platform-go-challenge/internal/models"
	"github.com/ceciivanov/platform-go-challenge/internal/repository"
)

// UserService struct defines methods related to user operations
type UserService struct {
	UserRepository repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		UserRepository: repo,
	}
}

// GetUserFavorites returns a map of user's favorite assets
func (s *UserService) GetUserFavorites(userID int) (map[int]models.Asset, error) {
	return s.UserRepository.GetUserFavorites(userID)
}

// AddUserFavorite adds an asset to the user's favorites
func (s *UserService) AddUserFavorite(userID int, asset models.Asset) error {
	return s.UserRepository.AddUserFavorite(userID, asset)
}

// DeleteUserFavorite deletes an asset from the user's favorites
func (s *UserService) DeleteUserFavorite(userID, assetID int) error {
	return s.UserRepository.DeleteUserFavorite(userID, assetID)
}

// EditUserFavorite edits an asset in the user's favorites
func (s *UserService) EditUserFavorite(userID int, assetID int, asset models.Asset) error {
	return s.UserRepository.EditUserFavorite(userID, assetID, asset)
}
