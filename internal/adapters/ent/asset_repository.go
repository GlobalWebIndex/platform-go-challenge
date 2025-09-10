package entadapter

import (
	"context"

	"github.com/SokratisChaimanas/platform-go-challenge/ent"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/ports"
	"github.com/google/uuid"
)

// Ensure  ports.AssetRepository interface implementation.
var _ ports.AssetRepository = (*AssetRepo)(nil)

type AssetRepo struct {
	client *ent.Client
}

func NewAssetRepo(client *ent.Client) *AssetRepo {
	return &AssetRepo{client: client}
}

func (assetRepo *AssetRepo) Get(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	a, err := assetRepo.client.Asset.Get(ctx, id)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.ErrAssetNotFound
		}
		return nil, err
	}

	return &domain.Asset{
		ID:          a.ID,
		Type:        domain.AssetType(a.AssetType),
		Description: a.Description,
		Payload:     a.Payload,
	}, nil
}

// Update persists changes from the domain model.
// It does not update immutable fields.
func (assetRepo *AssetRepo) Update(ctx context.Context, updatedAsset *domain.Asset) error {
	_, err := assetRepo.client.Asset.
		UpdateOneID(updatedAsset.ID).
		SetDescription(updatedAsset.Description).
		SetPayload(updatedAsset.Payload).
		Save(ctx)

	if ent.IsNotFound(err) {
		return domain.ErrAssetNotFound
	}

	return err
}
