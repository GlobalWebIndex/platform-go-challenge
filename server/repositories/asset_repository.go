package repositories

import (
	"hash/fnv"
	"sync"

	"github.com/google/uuid"
	"favourite_assets/server/models"
	"favourite_assets/server/errors"
)

const assetShardCount = 16

type assetShard struct {
	mu     sync.RWMutex
	assets map[uuid.UUID]models.Asset
}

type AssetRepository struct {
	shards [assetShardCount]*assetShard
}

func NewAssetRepository() *AssetRepository {
	r := &AssetRepository{}
	for i := 0; i < assetShardCount; i++ {
		r.shards[i] = &assetShard{
			assets: make(map[uuid.UUID]models.Asset),
		}
	}
	return r
}

func (r *AssetRepository) pickShard(assetID uuid.UUID) *assetShard {
	h := fnv.New32a()
	h.Write(assetID[:])
	return r.shards[uint(h.Sum32())%assetShardCount]
}

func (r *AssetRepository) Create(asset models.Asset) error {
	shard := r.pickShard(asset.GetID())
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, exists := shard.assets[asset.GetID()]; exists {
		return errors.ErrAssetExists
	}

	shard.assets[asset.GetID()] = asset
	return nil
}

func (r *AssetRepository) GetByID(id uuid.UUID) (models.Asset, error) {
	shard := r.pickShard(id)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	asset, ok := shard.assets[id]
	if !ok {
		return nil, errors.ErrAssetNotFound
	}
	return asset, nil
}

func (r *AssetRepository) Update(asset models.Asset) error {
	shard := r.pickShard(asset.GetID())
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, ok := shard.assets[asset.GetID()]; !ok {
		return errors.ErrAssetNotFound
	}

	shard.assets[asset.GetID()] = asset
	return nil
}

func (r *AssetRepository) Delete(id uuid.UUID) error {
	shard := r.pickShard(id)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, ok := shard.assets[id]; !ok {
		return errors.ErrAssetNotFound
	}

	delete(shard.assets, id)
	return nil
}

func (r *AssetRepository) ListAll() []models.Asset {
	result := []models.Asset{}
	for i := 0; i < assetShardCount; i++ {
		shard := r.shards[i]
		shard.mu.RLock()
		for _, asset := range shard.assets {
			result = append(result, asset)
		}
		shard.mu.RUnlock()
	}
	return result
}
