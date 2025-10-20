// Package store: in-memory implementation.
package db

import (
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/mr-polymath/platform-go-challenge/internal/model"
)

// ErrNotFound is returned when an operation targets a missing asset.
var ErrNotFound = errors.New("not found")

// MemoryStore is a concurrency-safe in-memory Store implementation. It is
// intentionally simple for the challenge and stores assets in a nested map:
// map[userID]map[assetID]Asset.
type MemoryStore struct {
	mu sync.RWMutex
	// map[userID]map[assetID]Asset
	favourites map[uuid.UUID]map[uuid.UUID]model.Asset
}

// NewMemoryStore constructs a ready-to-use MemoryStore.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{favourites: make(map[uuid.UUID]map[uuid.UUID]model.Asset)}
}

// GetAll returns all assets favourited by the given user.
func (m *MemoryStore) GetAllFavourites(userID uuid.UUID, q PagedQuery) ([]model.Asset, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	assets := []model.Asset{}
	if userFavourites, ok := m.favourites[userID]; ok {
		for _, asset := range userFavourites {
			assets = append(assets, asset)
		}
	}

	// Sorting by q.Order or ID by default, because map does not
	// guarantee the order  of iteration
	sort.Slice(assets, func(i, j int) bool {
		switch strings.ToLower(q.SortBy) {
		case "description":
			if strings.ToLower(q.Order) == "desc" {
				return assets[i].Description > assets[j].Description
			}
			return assets[i].Description < assets[j].Description
		case "type":
			if strings.ToLower(q.Order) == "desc" {
				return assets[i].Type > assets[j].Type
			}
			return assets[i].Type < assets[j].Type
		default:
			if strings.ToLower(q.Order) == "desc" {
				return assets[i].ID.String() > assets[j].ID.String()
			}
			return assets[i].ID.String() < assets[j].ID.String()
		}
	})

	// Apply pagination
	start := q.Offset
	if start > len(assets) {
		return []model.Asset{}, nil
	}

	end := start + q.Limit
	if q.Limit <= 0 || end > len(assets) {
		end = len(assets)
	}

	return assets[start:end], nil
}

// Add inserts an asset into a user's favourites. If the asset has no ID it
// will be assigned a timestamp-based ID and CreatedAt will be set.
func (m *MemoryStore) AddFavourite(userID uuid.UUID, a *model.Asset) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.favourites[userID]; !ok {
		m.favourites[userID] = map[uuid.UUID]model.Asset{}
	}
	a.ID = uuid.New()
	a.CreatedAt = time.Now().UTC()
	a.UpdatedAt = a.CreatedAt
	m.favourites[userID][a.ID] = *a
	return nil
}

// Remove deletes an asset from a user's favourites.
func (m *MemoryStore) RemoveFavourite(userID uuid.UUID, assetID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if userMap, ok := m.favourites[userID]; ok {
		if _, ok := userMap[assetID]; ok {
			delete(userMap, assetID)
			return nil
		}
	}

	return ErrNotFound
}

// UpdateDescription updates only the Description field of the asset.
func (m *MemoryStore) PatchFavourite(userID uuid.UUID, assetID uuid.UUID, description string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if userMap, ok := m.favourites[userID]; ok {
		if a, ok := userMap[assetID]; ok {
			a.Description = description
			a.UpdatedAt = time.Now().UTC()
			userMap[assetID] = a
			return nil
		}
	}
	return ErrNotFound
}
