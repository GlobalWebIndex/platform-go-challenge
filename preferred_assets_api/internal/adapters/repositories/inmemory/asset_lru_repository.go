package inmemory

import (
	"errors"
	"fmt"
	"sync"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/ports"
	lru "github.com/hashicorp/golang-lru/v2"
)

var (
	ErrAssetNotFound = errors.New("asset not found")
)

var _ ports.AssetRepository = (*LRUAssetRepositoryImpl)(nil)

type LRUAssetRepositoryImpl struct {
	cache *lru.Cache[string, entities.AssetEntity]
	mu    sync.RWMutex
}

func NewAssetRepository(cache *lru.Cache[string, entities.AssetEntity]) *LRUAssetRepositoryImpl {
	return &LRUAssetRepositoryImpl{cache: cache}
}

func (r *LRUAssetRepositoryImpl) Exists(id string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	exists := r.cache.Contains(id)
	return exists, nil
}

func (r *LRUAssetRepositoryImpl) Save(asset entities.AssetEntity) (entities.AssetEntity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := asset.Validate(); err != nil {
		return entities.AssetBaseEntity{}, err
	}

	val, ok := r.cache.Get(asset.GetID())
	if !ok {
		return entities.AssetBaseEntity{}, ErrAssetNotFound
	}

	return val, nil
}

func (r *LRUAssetRepositoryImpl) GetByID(id string) (entities.AssetEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, ok := r.cache.Get(id)
	if !ok {
		return nil, ErrAssetNotFound
	}
	return val, nil
}

func (r *LRUAssetRepositoryImpl) GetByIDs(ids []string) ([]entities.AssetEntity, error) {
	assets := make([]entities.AssetEntity, 0, len(ids))

	for _, id := range ids {
		asset, err := r.GetByID(id)
		if err != nil {
			continue
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (r *LRUAssetRepositoryImpl) GetAll() ([]entities.AssetEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	assets := make([]entities.AssetEntity, 0, r.cache.Len())
	for _, key := range r.cache.Keys() {
		if val, ok := r.cache.Peek(key); ok {
			assets = append(assets, val)
		}
	}
	return assets, nil
}

func (r *LRUAssetRepositoryImpl) GetByType(typeId entities.AssetType) ([]entities.AssetEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	assets := make([]entities.AssetEntity, 0)
	for _, key := range r.cache.Keys() {
		if val, ok := r.cache.Peek(key); ok && val.GetType() == typeId {
			assets = append(assets, val)
		}
	}
	return assets, nil
}

func (r *LRUAssetRepositoryImpl) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.cache.Contains(id) {
		return ErrAssetNotFound
	}

	r.cache.Remove(id)
	return nil
}

func (r *LRUAssetRepositoryImpl) Update(asset entities.AssetEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.cache.Contains(asset.GetID()) {
		return ErrAssetNotFound
	}

	if err := asset.Validate(); err != nil {
		return err
	}

	r.cache.Add(asset.GetID(), asset)
	return nil
}

func (r *LRUAssetRepositoryImpl) GetAudienceByID(id string) (*entities.AudienceEntity, error) {
	asset, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	audience, ok := asset.(*entities.AudienceEntity)
	if !ok {
		return nil, fmt.Errorf("asset with ID %s is not an Audience", id)
	}

	return audience, nil
}

func (r *LRUAssetRepositoryImpl) GetChartByID(id string) (*entities.ChartEntity, error) {
	asset, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	chart, ok := asset.(*entities.ChartEntity)
	if !ok {
		return nil, fmt.Errorf("asset with ID %s is not a Chart", id)
	}

	return chart, nil
}

func (r *LRUAssetRepositoryImpl) GetInsightByID(id string) (*entities.InsightEntity, error) {
	asset, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	insight, ok := asset.(*entities.InsightEntity)
	if !ok {
		return nil, fmt.Errorf("asset with ID %s is not an Insight", id)
	}

	return insight, nil
}
