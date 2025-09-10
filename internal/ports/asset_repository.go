package ports

import (
	"context"

	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/google/uuid"
)

// AssetRepository stores and retrieves assets.
type AssetRepository interface {
	// Get returns the asset or domain.ErrAssetNotFound.
	Get(ctx context.Context, id uuid.UUID) (*domain.Asset, error)

	// Update persists changes to an asset.
	Update(ctx context.Context, a *domain.Asset) error
}
