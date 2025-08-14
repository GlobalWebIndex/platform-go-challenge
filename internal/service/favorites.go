package service

import (
	"context"

	"gwi-favorites-service/internal/domain"
	"gwi-favorites-service/internal/repository"

	"github.com/sirupsen/logrus"
)

// FavoritesService handles business logic for favorites
type FavoritesService struct {
	repo   repository.FavoritesRepository
	logger *logrus.Logger
}

// NewFavoritesService creates a new favorites service
func NewFavoritesService(repo repository.FavoritesRepository, logger *logrus.Logger) *FavoritesService {
	return &FavoritesService{
		repo:   repo,
		logger: logger,
	}
}

// GetUserFavorites retrieves all favorites for a user
func (s *FavoritesService) GetUserFavorites(ctx context.Context, userID string, limit, offset int) ([]*domain.UserFavorite, error) {
	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"limit":   limit,
		"offset":  offset,
	}).Info("Getting user favorites")

	if userID == "" {
		return nil, domain.ErrInvalidUserID
	}

	favorites, err := s.repo.GetUserFavorites(userID, limit, offset)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get user favorites")
		return nil, err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id": userID,
		"count":   len(favorites),
	}).Info("Successfully retrieved user favorites")

	return favorites, nil
}

// AddFavorite adds an asset to user's favorites
func (s *FavoritesService) AddFavorite(ctx context.Context, userID string, asset domain.Asset) error {
	s.logger.WithFields(logrus.Fields{
		"user_id":    userID,
		"asset_id":   asset.GetID(),
		"asset_type": asset.GetType(),
	}).Info("Adding asset to favorites")

	if userID == "" {
		return domain.ErrInvalidUserID
	}

	if err := asset.Validate(); err != nil {
		s.logger.WithError(err).WithField("asset_id", asset.GetID()).Error("Asset validation failed")
		return err
	}

	// Check if asset exists, if not create it
	if _, err := s.repo.GetAsset(asset.GetID()); err == domain.ErrAssetNotFound {
		if err := s.repo.CreateAsset(asset); err != nil {
			s.logger.WithError(err).WithField("asset_id", asset.GetID()).Error("Failed to create asset")
			return err
		}
	}

	if err := s.repo.AddFavorite(userID, asset); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":  userID,
			"asset_id": asset.GetID(),
		}).Error("Failed to add favorite")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"asset_id": asset.GetID(),
	}).Info("Successfully added asset to favorites")

	return nil
}

// RemoveFavorite removes an asset from user's favorites
func (s *FavoritesService) RemoveFavorite(ctx context.Context, userID, assetID string) error {
	s.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"asset_id": assetID,
	}).Info("Removing asset from favorites")

	if userID == "" {
		return domain.ErrInvalidUserID
	}

	if assetID == "" {
		return domain.ErrInvalidInput
	}

	if err := s.repo.RemoveFavorite(userID, assetID); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":  userID,
			"asset_id": assetID,
		}).Error("Failed to remove favorite")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"asset_id": assetID,
	}).Info("Successfully removed asset from favorites")

	return nil
}

// UpdateFavoriteDescription updates the description of a favorite asset
func (s *FavoritesService) UpdateFavoriteDescription(ctx context.Context, userID, assetID, description string) error {
	s.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"asset_id": assetID,
	}).Info("Updating favorite asset description")

	if userID == "" {
		return domain.ErrInvalidUserID
	}

	if assetID == "" {
		return domain.ErrInvalidInput
	}

	// Check if it's a favorite
	isFavorite, err := s.repo.IsFavorite(userID, assetID)
	if err != nil {
		return err
	}
	if !isFavorite {
		return domain.ErrFavoriteNotFound
	}

	// Get the asset
	asset, err := s.repo.GetAsset(assetID)
	if err != nil {
		return err
	}

	// Update description
	asset.SetDescription(description)

	// Update in repository
	if err := s.repo.UpdateAsset(asset); err != nil {
		s.logger.WithError(err).WithFields(logrus.Fields{
			"user_id":  userID,
			"asset_id": assetID,
		}).Error("Failed to update asset description")
		return err
	}

	s.logger.WithFields(logrus.Fields{
		"user_id":  userID,
		"asset_id": assetID,
	}).Info("Successfully updated favorite asset description")

	return nil
}

// GetFavoriteCount returns the count of user's favorites
func (s *FavoritesService) GetFavoriteCount(ctx context.Context, userID string) (int, error) {
	if userID == "" {
		return 0, domain.ErrInvalidUserID
	}

	count, err := s.repo.GetFavoriteCount(userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to get favorite count")
		return 0, err
	}

	return count, nil
}

// IsFavorite checks if an asset is in user's favorites
func (s *FavoritesService) IsFavorite(ctx context.Context, userID, assetID string) (bool, error) {
	if userID == "" || assetID == "" {
		return false, domain.ErrInvalidInput
	}

	return s.repo.IsFavorite(userID, assetID)
}
