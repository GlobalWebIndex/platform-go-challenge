//https://dev.to/team3/repository-pattern-in-golang-a-practical-guide-1kla

package repository

import (
	"errors"
	"go-challenge/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

type FavouriteRepository interface {
	GetUserFavourites(userID string, limit, offset int) ([]models.Asset, int, error)
	AddFavourite(userID, assetID string) (*models.Favourite, error)
	RemoveFavourite(userID, assetID string) error
	UpdateAssetDescription(assetID, description string) (*models.Asset, error)
}

// stores favourites in memory (RAM).
type InMemoryFavouriteRepository struct {
	// favourites stores the relationship between users and their favourited assets.
	// Structure: map[userID][]Favourite
	// Example: {"user-001": [{ID: "fav1", AssetID: "asset-001"}, {ID: "fav2", AssetID: "asset-002"}]}
	favourites map[string][]models.Favourite

	// assets stores all available assets in the system.
	// Structure: map[assetID]Asset
	// Example: {"asset-001": {ID: "asset-001", Type: "insight", ...}}
	assets map[string]models.Asset

	// (Read-Write Mutex) that makes this repository thread-safe.
	mu sync.RWMutex
}

func NewInMemoryFavouriteRepository() *InMemoryFavouriteRepository {
	// Create the instance
	// Initializes the maps (or they'll be nil)
	favouritesRepo := &InMemoryFavouriteRepository{
		favourites: make(map[string][]models.Favourite),
		assets:     make(map[string]models.Asset),
	}

	// // Seed with mock assets (the shared catalog all users can access)
	// favouritesRepo.seedMockAssets()

	// Seed with mock user's favourites
	favouritesRepo.seedMockFavourites()

	return favouritesRepo
}

func (r *InMemoryFavouriteRepository) GetUserFavourites(userID string, limit, offset int) ([]models.Asset, int, error) {
	//allows multiple concurrent reads but blocks writes.
	r.mu.RLock()
	// run this when the function exit(even if there's an error)
	defer r.mu.RUnlock()

	// return the favourites list or an empty slice.
	userFavourites := r.favourites[userID]

	// calculate total count
	totalCount := len(userFavourites)

	// Check if offset is past the end of the list
	if offset >= totalCount {
		return []models.Asset{}, totalCount, nil
	}

	// Calculate the end index for slicing
	end := min(offset+limit, totalCount)

	//?????
	paginatedFavourites := userFavourites[offset:end]

	// Init slice with 0 items
	// slice capacity equals paginatedFavourites length.
	assets := make([]models.Asset, 0, len(paginatedFavourites))

	for _, fav := range paginatedFavourites {
		if asset, exists := r.assets[fav.AssetID]; exists {
			assets = append(assets, asset)
		}
		// If the asset doesn't exist, we skip it silently.
	}

	return assets, totalCount, nil
}

func (r *InMemoryFavouriteRepository) AddFavourite(userID, assetID string) (*models.Favourite, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if asset exists
	if _, exists := r.assets[assetID]; !exists {
		return nil, errors.New("asset not found")
	}

	// Check if already favorited
	for _, fav := range r.favourites[userID] {
		if fav.AssetID == assetID {
			return nil, errors.New("asset is already favourite")
		}
	}

	// create favourite asset
	favourite := models.Favourite{
		ID:        uuid.New().String(),
		UserID:    userID,
		AssetID:   assetID,
		CreatedAt: time.Now(),
	}

	r.favourites[userID] = append(r.favourites[userID], favourite)

	return &favourite, nil
}

func (r *InMemoryFavouriteRepository) RemoveFavourite(userID, assetID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	userFavorites := r.favourites[userID]

	for i, fav := range userFavorites {
		if fav.AssetID == assetID {
			// https://go.dev/tour/moretypes/15
			// https://medium.com/codezillas/building-a-restful-api-with-go-part-2-22d99285c130

			// Remove from slice
			r.favourites[userID] = append(userFavorites[:i], userFavorites[i+1:]...)
			return nil
		}
	}

	return errors.New("favourite could not be removed")
}

func (r *InMemoryFavouriteRepository) UpdateAssetDescription(assetID, description string) (*models.Asset, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Get the asset
	asset, exists := r.assets[assetID]
	if !exists {
		return nil, errors.New("asset not found")
	}

	// Update the description
	asset.Description = description
	asset.UpdatedAt = time.Now()

	// Save back to map
	r.assets[assetID] = asset

	return &asset, nil
}
