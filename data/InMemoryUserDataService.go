package data

import (
	"errors"
	"platform-go-challenge/models"
)

// InMemoryUserDataService is an in-memory implementation of UserDataService.
type InMemoryUserDataService struct {
	users map[string]*models.User
}

// NewInMemoryUserDataService creates a new InMemoryUserDataService.
func NewInMemoryUserDataService() *InMemoryUserDataService {
	return &InMemoryUserDataService{
		users: make(map[string]*models.User),
	}
}

// GetUser retrieves a user by ID.
func (s *InMemoryUserDataService) GetUser(userID string) (*models.User, error) {
	user, ok := s.users[userID]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

// AddUser adds a new user.
func (s *InMemoryUserDataService) AddUser(user *models.User) error {
	s.users[user.ID] = user
	return nil
}

// AddAssetToFavourites adds an asset to the user's list of favourites.
func (s *InMemoryUserDataService) AddAssetToFavourites(userID string, asset models.Asset) error {
	user, ok := s.users[userID]
	if !ok {
		user = &models.User{ID: userID}
		err := s.AddUser(user)
		if err != nil {
			return err
		}
	}

	// Helper function to check for duplicate asset ID
	exists := func(favourites []models.Asset, assetID string) bool {
		for _, favAsset := range favourites {
			if favAsset.ID == assetID {
				return true
			}
		}
		return false
	}

	// Use the helper function to check for duplicate asset ID
	if exists(user.Favourites, asset.ID) {
		return errors.New("duplicate asset ID")
	}

	user.Favourites = append(user.Favourites, asset)
	return nil
}

// RemoveAssetFromFavourites removes an asset from the user's list of favourites.
func (s *InMemoryUserDataService) RemoveAssetFromFavourites(userID string, assetID string) error {
	user, ok := s.users[userID]
	if !ok {
		return errors.New("user not found")
	}
	for i, asset := range user.Favourites {
		if asset.ID == assetID {
			user.Favourites = append(user.Favourites[:i], user.Favourites[i+1:]...)
			return nil
		}
	}
	return errors.New("asset not found")
}

// EditAssetDescription edits the description of an asset in the user's list of favourites.
func (s *InMemoryUserDataService) EditAssetDescription(userID string, assetID string, description string) error {
	user, ok := s.users[userID]
	if !ok {
		return errors.New("user not found")
	}
	for i, asset := range user.Favourites {
		if asset.ID == assetID {
			user.Favourites[i].Description = description
			return nil
		}
	}
	return errors.New("asset not found")
}
