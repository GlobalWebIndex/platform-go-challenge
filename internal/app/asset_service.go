package app

import (
	"context"

	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/ports"
	"github.com/google/uuid"
)

// AssetService coordinates asset operations (load → mutate → save).
// Domain enforces rules; ent handles persistence.
type AssetService struct {
	assetRepo ports.AssetRepository
}

func NewAssetService(assets ports.AssetRepository) *AssetService {
	return &AssetService{assetRepo: assets}
}

// EditDescription loads the asset, edits the description (domain rule),
// then persists changes. Returns the updated asset.
func (assetService *AssetService) EditDescription(ctx context.Context, assetID uuid.UUID, newDesc string) (*domain.Asset, error) {
	a, err := assetService.assetRepo.Get(ctx, assetID)
	if err != nil {
		return nil, err // expected: domain.ErrAssetNotFound
	}
	if err := a.EditDescription(newDesc); err != nil {
		return nil, err // expected: domain.ErrEmptyDescription
	}
	if err := assetService.assetRepo.Update(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}
