package store

import (
	"context"
	"sync"

	"github.com/gwi/favourites-api/internal/models"
)

// FavouriteStore defines the interface for favourite asset storage operations.
// This interface allows for multiple implementations (in-memory, PostgreSQL, MongoDB, etc.)
// while keeping the business logic decoupled from storage details.
type FavouriteStore interface {
	// GetUserFavourites retrieves all favourite assets for a given user
	GetUserFavourites(ctx context.Context, userID string) ([]models.Asset, error)

	// AddFavourite adds a new asset to a user's favourites
	AddFavourite(ctx context.Context, userID string, asset models.Asset) error

	// RemoveFavourite removes an asset from a user's favourites
	RemoveFavourite(ctx context.Context, userID string, assetID string) error

	// UpdateDescription updates the description of a favourited asset
	UpdateDescription(ctx context.Context, userID string, assetID string, description string) error

	// GetFavourite retrieves a specific favourite asset
	GetFavourite(ctx context.Context, userID string, assetID string) (models.Asset, error)
}

// InMemoryFavouriteStore provides a thread-safe in-memory implementation
// of the FavouriteStore interface. Uses RWMutex for concurrent read optimization.
type InMemoryFavouriteStore struct {
	mu   sync.RWMutex
	data map[string]map[string]models.Asset // userID -> assetID -> asset
}

func NewInMemoryFavouriteStore() *InMemoryFavouriteStore {
	return &InMemoryFavouriteStore{
		data: make(map[string]map[string]models.Asset),
	}
}

// GetUserFavourites retrieves all favourite assets for a user
func (s *InMemoryFavouriteStore) GetUserFavourites(
	ctx context.Context,
	userID string,
) ([]models.Asset, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	userFavs, ok := s.data[userID]
	if !ok {
		return []models.Asset{}, nil // Return empty slice for users with no favourites
	}

	result := make([]models.Asset, 0, len(userFavs))
	for _, asset := range userFavs {
		result = append(result, asset)
	}

	return result, nil
}

// AddFavourite adds a new asset to a user's favourites
func (s *InMemoryFavouriteStore) AddFavourite(
	ctx context.Context,
	userID string,
	asset models.Asset,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	// Initialize user's favourites map if it doesn't exist
	if _, ok := s.data[userID]; !ok {
		s.data[userID] = make(map[string]models.Asset)
	}

	// Check if asset already exists
	if _, exists := s.data[userID][asset.GetAssetID()]; exists {
		return ErrConflict
	}

	s.data[userID][asset.GetAssetID()] = asset

	return nil
}

// RemoveFavourite removes an asset from a user's favourites
func (s *InMemoryFavouriteStore) RemoveFavourite(
	ctx context.Context,
	userID string,
	assetID string,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	userFavs, ok := s.data[userID]
	if !ok {
		return ErrNotFound
	}

	if _, exists := userFavs[assetID]; !exists {
		return ErrNotFound
	}

	delete(userFavs, assetID)

	// Clean up empty user map
	if len(userFavs) == 0 {
		delete(s.data, userID)
	}

	return nil
}

// UpdateDescription updates the description of a favourited asset
func (s *InMemoryFavouriteStore) UpdateDescription(
	ctx context.Context,
	userID string,
	assetID string,
	description string,
) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	userFavs, ok := s.data[userID]
	if !ok {
		return ErrNotFound
	}

	asset, exists := userFavs[assetID]
	if !exists {
		return ErrNotFound
	}

	// Update the description
	asset.SetDescription(description)
	s.data[userID][assetID] = asset

	return nil
}

// GetFavourite retrieves a specific favourite asset
func (s *InMemoryFavouriteStore) GetFavourite(
	ctx context.Context,
	userID string,
	assetID string,
) (models.Asset, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	userFavs, ok := s.data[userID]
	if !ok {
		return nil, ErrNotFound
	}

	asset, exists := userFavs[assetID]
	if !exists {
		return nil, ErrNotFound
	}

	return asset, nil
}
