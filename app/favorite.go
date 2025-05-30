package app

import (
	"context"

	"pgc/domain"
)

// FavoriteService is a struct that contains the favorite repository which is used to manage user favorite assets.
type FavoriteService struct {
	favoriteRepository domain.FavoriteRepository
}

// GetUserFavorites retrieves all favorites for a specific user.
func (s FavoriteService) GetUserFavorites(ctx context.Context, userID int) ([]domain.Favorite, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidUserID
	}

	return s.favoriteRepository.SearchByUserID(ctx, userID)
}

// Add adds a new favorite asset for a user.
func (s FavoriteService) Add(ctx context.Context, favorite domain.Favorite) (domain.Favorite, error) {
	switch {
	case favorite.UserID <= 0:
		return domain.Favorite{}, domain.ErrInvalidUserID
	case favorite.AssetID <= 0:
		return domain.Favorite{}, domain.ErrInvalidAssetID
	}

	favorite.ID = 0 // Ensure ID is zero (default value) for insertion

	return s.favoriteRepository.Insert(ctx, favorite)
}

// Edit updates the description of an existing favorite asset.
func (s FavoriteService) Edit(ctx context.Context, favoriteID, userID int, description string) (domain.Favorite, error) {
	switch {
	case favoriteID <= 0:
		return domain.Favorite{}, domain.ErrInvalidFavoriteID
	case userID <= 0:
		return domain.Favorite{}, domain.ErrInvalidUserID
	}

	return s.favoriteRepository.UpdateDescription(ctx, favoriteID, userID, description)
}

// Delete removes a favorite asset for a user.
func (s FavoriteService) Delete(ctx context.Context, favoriteID, userID int) error {
	switch {
	case favoriteID <= 0:
		return domain.ErrInvalidFavoriteID
	case userID <= 0:
		return domain.ErrInvalidUserID
	}

	return s.favoriteRepository.Remove(ctx, favoriteID, userID)
}

// NewFavoriteService creates a new instance of FavoriteService with the provided favorite repository.
func NewFavoriteService(favoriteRepository domain.FavoriteRepository) *FavoriteService {
	return &FavoriteService{favoriteRepository: favoriteRepository}
}
