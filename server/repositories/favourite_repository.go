package repositories

import (
	"hash/fnv"
	"sync"

	"github.com/google/uuid"
	"favourite_assets/server/models"
	"favourite_assets/server/errors"
)

const favouriteShardCount = 16

type favouriteShard struct {
	mu        sync.RWMutex
	favourites map[uuid.UUID]*models.Favourite
}

type FavouriteRepository struct {
	shards [favouriteShardCount]*favouriteShard
}

// NewFavouriteRepository initializes shards
func NewFavoriteRepository() *FavouriteRepository {
	r := &FavouriteRepository{}
	for i := 0; i < favouriteShardCount; i++ {
		r.shards[i] = &favouriteShard{
			favourites: make(map[uuid.UUID]*models.Favourite),
		}
	}
	return r
}

// pickShard selects a shard based on favourite ID
func (r *FavouriteRepository) pickShard(favID uuid.UUID) *favouriteShard {
	h := fnv.New32a()
	h.Write(favID[:])
	return r.shards[uint(h.Sum32())%favouriteShardCount]
}

func (r *FavouriteRepository) Create(fav *models.Favourite) error {
	shard := r.pickShard(fav.ID)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, exists := shard.favourites[fav.ID]; exists {
		return errors.ErrFavouriteExists
	}

	shard.favourites[fav.ID] = fav
	return nil
}

func (r *FavouriteRepository) GetByID(favID uuid.UUID) (*models.Favourite, error) {
	shard := r.pickShard(favID)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	fav, ok := shard.favourites[favID]
	if !ok {
		return nil, errors.ErrFavouriteNotFound
	}
	return fav, nil
}

func (r *FavouriteRepository) Delete(favID uuid.UUID) error {
	shard := r.pickShard(favID)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, ok := shard.favourites[favID]; !ok {
		return errors.ErrFavouriteNotFound
	}

	delete(shard.favourites, favID)
	return nil
}

func (r *FavouriteRepository) ListByUser(userID uuid.UUID) []*models.Favourite {
	var result []*models.Favourite
	for i := 0; i < favouriteShardCount; i++ {
		shard := r.shards[i]
		shard.mu.RLock()
		for _, fav := range shard.favourites {
			if fav.UserID == userID {
				result = append(result, fav)
			}
		}
		shard.mu.RUnlock()
	}
	return result
}

func (r *FavouriteRepository) Get(favID uuid.UUID) (models.Favourite, error) {
	shard := r.pickShard(favID)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	fav, ok := shard.favourites[favID]
	if !ok {
		return models.Favourite{}, errors.ErrFavouriteNotFound
	}

	return *fav, nil
}