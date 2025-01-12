package service

import (
	"context"
	"errors"
	"platform-go-challenge/internal/domain"
	"platform-go-challenge/internal/service/types"
	"time"

	"github.com/google/uuid"
)

type assetService struct {
	repo domain.Repository
}

func AssetService(repo domain.Repository) domain.AssetService {
	return &assetService{repo: repo}
}

func (s *assetService) AddFavorite(ctx context.Context, userID, assetID string) error {
	err := validateAssetId(assetID)
	if err != nil {
		return err
	}
	// Check if asset exists
	if _, err := s.repo.GetAssets(ctx, []any{assetID}); err != nil {
		return err
	}

	favorite := types.Favorite{
		UserID:    userID,
		AssetID:   assetID,
		CreatedAt: time.Now(),
	}

	return s.repo.AddFavorite(ctx, favorite)
}

func (s *assetService) RemoveFavorite(ctx context.Context, userID, assetID string) error {
	err := validateAssetId(assetID)
	if err != nil {
		return err
	}
	return s.repo.RemoveFavorite(ctx, userID, assetID)
}

func (s *assetService) ListFavorites(ctx context.Context, userID string, page, limit int) ([]types.Asset, error) {
	if limit <= 0 {
		limit = 20 // Default limit
	}
	if page <= 0 {
		page = 1
	}

	favorites, err := s.repo.ListFavorites(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	ids := make([]any, 0, len(favorites))
	for _, fav := range favorites {
		ids = append(ids, fav.AssetID)
	}

	assets, err := s.repo.GetAssets(ctx, ids)
	if err != nil {
		return nil, err
	}

	return assets, nil
}

func (s *assetService) UpdateDescription(ctx context.Context, userID, assetID, description string) error {
	err := validateAssetId(assetID)
	if err != nil {
		return err
	}

	assets, err := s.repo.GetAssets(ctx, []any{assetID})
	if err != nil {
		return err
	}

	if len(assets) > 1 {
		// impossible case, assetIDs are unique
		return errors.New("non unique asset id")
	}

	asset := assets[0]

	if err := asset.SetDescription(description); err != nil {
		return err
	}

	return s.repo.UpdateAsset(ctx, asset)
}

func validateAssetId(assetId string) error {
	if _, err := uuid.Parse(assetId); err != nil {
		return domain.ErrInvalidInput
	}
	return nil
}
