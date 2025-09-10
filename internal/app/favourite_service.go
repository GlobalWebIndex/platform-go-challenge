package app

import (
	"context"

	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/ports"
	"github.com/google/uuid"
)

// FavouritesService coordinates multiple repositories and applies business rules.
type FavouritesService struct {
	userRepo  ports.UserRepository
	assetRepo ports.AssetRepository
	favRepo   ports.FavouriteRepository
}

func NewFavouritesService(
	userRepo ports.UserRepository,
	assetRepo ports.AssetRepository,
	favRepo ports.FavouriteRepository,
) *FavouritesService {
	return &FavouritesService{
		userRepo:  userRepo,
		assetRepo: assetRepo,
		favRepo:   favRepo,
	}
}

// Add validates user and asset existence, prevents duplicates, then creates a favourite.
func (favService *FavouritesService) Add(ctx context.Context, userID, assetID uuid.UUID) (domain.Favourite, error) {
	// Ensure user exists.
	ok, err := favService.userRepo.Exists(ctx, userID)
	if err != nil {
		return domain.Favourite{}, err
	}
	if !ok {
		return domain.Favourite{}, domain.ErrUserNotFound
	}

	// Ensure asset exists.
	if _, err := favService.assetRepo.Get(ctx, assetID); err != nil {
		return domain.Favourite{}, err // expected: domain.ErrAssetNotFound
	}

	// Prevent duplicates.
	exists, err := favService.favRepo.Exists(ctx, userID, assetID)
	if err != nil {
		return domain.Favourite{}, err
	}
	if exists {
		return domain.Favourite{}, domain.ErrFavouriteAlreadyExists
	}

	// Create favourite.
	favToReturn := domain.NewFavourite(userID, assetID)
	if err := favService.favRepo.Create(ctx, &favToReturn); err != nil {
		return domain.Favourite{}, err
	}
	return favToReturn, nil
}

// Remove deletes a favourite. Missing pair should return domain.ErrFavouriteNotFound.
func (favService *FavouritesService) Remove(ctx context.Context, userID, assetID uuid.UUID) error {
	return favService.favRepo.Delete(ctx, userID, assetID)
}

// ListByUserKeyset returns a user's favourited assets using a keyset cursor.
func (favService *FavouritesService) ListByUserKeyset(ctx context.Context, userID uuid.UUID, limit int, after string) ([]domain.Asset, *string, error) {
	// Validate the user exists (same as ListByUser does).
	ok, err := favService.userRepo.Exists(ctx, userID)
	if err != nil {
		return nil, nil, err
	}
	if !ok {
		return nil, nil, domain.ErrUserNotFound
	}
	// Delegate to repository (repo clamps limit to defaults/caps just like your offset path).
	return favService.favRepo.ListAssetsFavouritedByUserKeyset(ctx, userID, limit, after)
}
