package service

import (
	"context"
	"platform-go-challenge/internal/domain/entities"
	"platform-go-challenge/internal/repository"

	"go.uber.org/zap"
)

type FavouriteService struct {
	repo   repository.Repository
	logger *zap.Logger
}

func NewFavouriteService(r repository.Repository, l *zap.Logger) *FavouriteService {
	return &FavouriteService{
		repo:   r,
		logger: l,
	}
}

// GetFavourites now returns a slice of a different type and an error.
func (s *FavouriteService) GetFavourites(ctx context.Context, userID string, page, size int) []entities.Asset {
	// A real implementation would call the repository here.
	s.logger.Info("retrieving favourites", zap.String("userID", userID))

	// For now, it returns a dummy value.
	result, _ := s.repo.GetFavourites(userID)
	return result
}

func (s *FavouriteService) AddFavourite(ctx context.Context, userID string, assetID string) error {
	s.logger.Info("adding favourite", zap.String("userID", userID), zap.String("assetID", ""))
	// A real implementation would call the repository here.
	if err := s.repo.AddFavourite(userID, assetID); err != nil {
		return err
	}
	return nil
}

func (s *FavouriteService) UpdateFavouriteDescription(ctx context.Context, userID string, assetID string, description string) error {
	s.logger.Info("updating favourite description", zap.String("userID", userID), zap.String("assetID", assetID))

	if err := s.repo.UpdateFavouriteDescription(userID, assetID, description); err != nil {
		return err
	}
	return nil
}

func (s *FavouriteService) RemoveFavourite(ctx context.Context, userID string, assetID string) error {
	s.logger.Info("removing favourite", zap.String("userID", userID), zap.String("assetID", assetID))

	if err := s.repo.RemoveFavourite(userID, assetID); err != nil {
		return err
	}
	return nil
}
