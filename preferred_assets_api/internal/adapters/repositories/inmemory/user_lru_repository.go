package inmemory

import (
	"errors"
	"sync"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/ports"
	lru "github.com/hashicorp/golang-lru/v2"
)

var _ ports.UserRepository = (*LRUUserRepositoryImpl)(nil)

type LRUUserRepositoryImpl struct {
	cache         *lru.Cache[string, *entities.UserEntity]
	favouriteRepo ports.FavouriteRepository
	mu            sync.RWMutex
}

func NewUserRepository(cache *lru.Cache[string, *entities.UserEntity], favouriteRepo ports.FavouriteRepository) *LRUUserRepositoryImpl {
	return &LRUUserRepositoryImpl{cache: cache,
		favouriteRepo: favouriteRepo}
}

func (r *LRUUserRepositoryImpl) Save(u entities.UserEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cache.Add(u.Id, &u)

	return nil
}

func (r *LRUUserRepositoryImpl) GetByID(id string) (entities.UserEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	val, ok := r.cache.Get(id)
	if !ok {
		return entities.UserEntity{}, errors.New("user not found")
	}

	return *val, nil
}

func (r *LRUUserRepositoryImpl) GetAll() ([]entities.UserEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]entities.UserEntity, 0, r.cache.Len())

	for _, key := range r.cache.Keys() {
		if val, ok := r.cache.Peek(key); ok {
			users = append(users, *val)
		}
	}

	return users, nil
}

func (r *LRUUserRepositoryImpl) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.cache.Remove(id)

	return nil
}

func (r *LRUUserRepositoryImpl) Update(u entities.UserEntity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.cache.Get(u.Id); !ok {
		return errors.New("user not found")
	}

	r.cache.Add(u.Id, &u)
	return nil
}

func (r *LRUUserRepositoryImpl) GetFavouritesByID(id string) ([]entities.FavouriteEntity, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// First verify user exists
	if _, ok := r.cache.Get(id); !ok {
		return nil, errors.New("user not found")
	}

	// Get favourites directly from favourite repository
	return r.favouriteRepo.GetByUserID(id)
}
