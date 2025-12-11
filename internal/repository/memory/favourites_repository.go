package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/domain/favourites"
)

// favouritesRepository is an in-memory implementation of favourites.Repository.
// Safe for concurrent use.
type favouritesRepository struct {
	mu sync.RWMutex

	// assets stores assets by ID.
	assets map[string]favourites.Asset

	// userFavourites stores favourites per user:
	//   userID -> assetID -> Favourite
	userFavourites map[string]map[string]favourites.Favourite
}

// NewFavouritesRepository creates a new in-memory favourites repository.
func NewFavouritesRepository() favourites.Repository {
	return &favouritesRepository{
		assets:         make(map[string]favourites.Asset),
		userFavourites: make(map[string]map[string]favourites.Favourite),
	}
}

// CreateFavourite creates an asset and marks it as a favourite.
// asset.ID must be empty; a new UUID v4 is generated.
func (r *favouritesRepository) CreateFavourite(
	ctx context.Context,
	userID string,
	asset favourites.Asset,
) (favourites.FavouriteEntry, error) {
	if err := ctx.Err(); err != nil {
		return favourites.FavouriteEntry{}, err
	}

	if asset.ID != "" {
		return favourites.FavouriteEntry{}, favourites.ErrInvalidAssetPayload
	}

	if err := asset.Validate(); err != nil {
		return favourites.FavouriteEntry{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate a new unique asset ID (UUID v4).
	assetID := uuid.NewString()
	asset.ID = assetID

	// Ensure per-user favourites map exists.
	uf, ok := r.userFavourites[userID]
	if !ok {
		uf = make(map[string]favourites.Favourite)
		r.userFavourites[userID] = uf
	}

	now := time.Now().UTC()
	fav := favourites.Favourite{
		UserID:    userID,
		AssetID:   assetID,
		CreatedAt: now,
	}

	r.assets[assetID] = asset
	uf[assetID] = fav

	return favourites.FavouriteEntry{
		Favourite: fav,
		Asset:     asset,
	}, nil
}

// ListFavourites returns favourites for the given user, optionally filtered by type,
// and applies pagination based on Limit/Offset in the filter. It returns the page
// of results plus the total number of matching entries (ignoring pagination).
func (r *favouritesRepository) ListFavourites(
	ctx context.Context,
	userID string,
	filter favourites.ListFilter,
) ([]favourites.FavouriteEntry, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	uf, ok := r.userFavourites[userID]
	if !ok {
		return []favourites.FavouriteEntry{}, 0, nil
	}

	entries := make([]favourites.FavouriteEntry, 0, len(uf))

	for _, fav := range uf {
		asset, ok := r.assets[fav.AssetID]
		if !ok {
			return nil, 0, fmt.Errorf("invariant violation: asset %q not found for favourite", fav.AssetID)
		}

		if filter.Type != nil && asset.Type != *filter.Type {
			continue
		}

		entries = append(entries, favourites.FavouriteEntry{
			Favourite: fav,
			Asset:     asset,
		})
	}

	total := len(entries)
	if total == 0 {
		return []favourites.FavouriteEntry{}, 0, nil
	}

	// Ensure deterministic ordering for pagination: sort by CreatedAt ascending.
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Favourite.CreatedAt.Before(entries[j].Favourite.CreatedAt)
	})

	// Apply pagination.
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}
	if offset >= total {
		return []favourites.FavouriteEntry{}, total, nil
	}

	limit := filter.Limit
	if limit <= 0 {
		// Treat non-positive limit as "no explicit limit" from this layer.
		limit = total - offset
	}

	end := offset + limit
	if end > total {
		end = total
	}

	page := entries[offset:end]
	return page, total, nil
}

// RemoveFavourite removes the favourite link. Idempotent: removing a favourite
// that does not exist is a no-op and does not return ErrFavouriteNotFound.
func (r *favouritesRepository) RemoveFavourite(
	ctx context.Context,
	userID, assetID string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	uf, ok := r.userFavourites[userID]
	if !ok {
		return nil
	}

	delete(uf, assetID)
	if len(uf) == 0 {
		delete(r.userFavourites, userID)
	}

	// Assets are kept even if no user references them (no cleanup needed for this challenge).
	return nil
}

// UpdateDescriptionForUser updates an asset's description if the user has it as a favourite.
// Returns ErrFavouriteNotFound if the favourite doesn't exist.
func (r *favouritesRepository) UpdateDescriptionForUser(
	ctx context.Context,
	userID, assetID, description string,
) (favourites.FavouriteEntry, error) {
	if err := ctx.Err(); err != nil {
		return favourites.FavouriteEntry{}, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	uf, ok := r.userFavourites[userID]
	if !ok {
		return favourites.FavouriteEntry{}, favourites.ErrFavouriteNotFound
	}

	fav, ok := uf[assetID]
	if !ok {
		return favourites.FavouriteEntry{}, favourites.ErrFavouriteNotFound
	}

	asset, ok := r.assets[assetID]
	if !ok {
		return favourites.FavouriteEntry{}, fmt.Errorf("invariant violation: asset %q not found for favourite", assetID)
	}

	asset.Description = description
	r.assets[assetID] = asset

	return favourites.FavouriteEntry{
		Favourite: fav,
		Asset:     asset,
	}, nil
}
