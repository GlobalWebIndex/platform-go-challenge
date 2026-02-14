package memory

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/gwi/platform-go-challenge/favourites/application"
	"github.com/gwi/platform-go-challenge/favourites/domain"
)

var _ application.FavouritesRepository = (*Repository)(nil)

type userAssetPair struct {
	userID  string
	assetID string
}

type Repository struct {
	mu               sync.RWMutex
	users            map[string]map[string]domain.Asset
	sourceAssetIndex map[string][]userAssetPair
	counter          int64
}

func NewRepository() *Repository {
	return &Repository{
		users:            make(map[string]map[string]domain.Asset),
		sourceAssetIndex: make(map[string][]userAssetPair),
	}
}

func (r *Repository) ListPaginated(userID string, limit, offset int) ([]domain.Asset, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	byID := r.users[userID]
	if len(byID) == 0 {
		return nil, 0
	}
	assets := make([]domain.Asset, 0, len(byID))
	for _, a := range byID {
		assets = append(assets, a)
	}
	sort.Slice(assets, func(i, j int) bool {
		return assets[i].CreatedAt.Before(assets[j].CreatedAt)
	})
	total := len(assets)
	if offset >= total {
		return []domain.Asset{}, total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	page := assets[offset:end]
	out := make([]domain.Asset, len(page))
	copy(out, page)
	return out, total
}

func (r *Repository) Get(userID string, assetID string) (*domain.Asset, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	byID := r.users[userID]
	if byID == nil {
		return nil, false
	}
	a, ok := byID[assetID]
	if !ok {
		return nil, false
	}
	return &a, true
}

// Add generates ID and timestamps if not set.
func (r *Repository) Add(userID string, asset domain.Asset) domain.Asset {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := time.Now()
	if asset.ID == "" {
		r.counter++
		asset.ID = fmt.Sprintf("asset-%d-%d", now.UnixNano(), r.counter)
	}
	asset.CreatedAt = now
	asset.UpdatedAt = now
	if r.users[userID] == nil {
		r.users[userID] = make(map[string]domain.Asset)
	}
	r.users[userID][asset.ID] = asset
	
	if asset.SourceAssetID != "" {
		r.sourceAssetIndex[asset.SourceAssetID] = append(
			r.sourceAssetIndex[asset.SourceAssetID],
			userAssetPair{userID: userID, assetID: asset.ID},
		)
	}
	
	return asset
}

func (r *Repository) Remove(userID string, assetID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	byID := r.users[userID]
	if byID == nil {
		return false
	}
	asset, ok := byID[assetID]
	if !ok {
		return false
	}
	
	if asset.SourceAssetID != "" {
		r.removeFromIndex(asset.SourceAssetID, userID, assetID)
	}
	
	delete(byID, assetID)
	return true
}

func (r *Repository) UpdateDescription(userID string, assetID string, description string) (*domain.Asset, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	byID := r.users[userID]
	if byID == nil {
		return nil, false
	}
	a, ok := byID[assetID]
	if !ok {
		return nil, false
	}
	a.Description = description
	a.UpdatedAt = time.Now()
	byID[assetID] = a
	return &a, true
}

// RemoveFavouritesBySourceAssetID is used by asset-deleted event to avoid orphans.
func (r *Repository) RemoveFavouritesBySourceAssetID(sourceAssetID string) int {
	if sourceAssetID == "" {
		return 0
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	
	pairs := r.sourceAssetIndex[sourceAssetID]
	if len(pairs) == 0 {
		return 0
	}
	
	for _, pair := range pairs {
		if byID := r.users[pair.userID]; byID != nil {
			delete(byID, pair.assetID)
		}
	}
	
	removed := len(pairs)
	delete(r.sourceAssetIndex, sourceAssetID)
	return removed
}

func (r *Repository) removeFromIndex(sourceAssetID, userID, assetID string) {
	pairs := r.sourceAssetIndex[sourceAssetID]
	for i, pair := range pairs {
		if pair.userID == userID && pair.assetID == assetID {
			r.sourceAssetIndex[sourceAssetID] = append(pairs[:i], pairs[i+1:]...)
			if len(r.sourceAssetIndex[sourceAssetID]) == 0 {
				delete(r.sourceAssetIndex, sourceAssetID)
			}
			return
		}
	}
}
