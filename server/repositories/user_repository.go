package repositories

import (
	"hash/fnv"
	"sync"
	"time"

	"github.com/google/uuid"
	"favourite_assets/server/models"
	"favourite_assets/server/errors"
)

const userShardCount = 16

type userShard struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*models.User
}

type UserRepository struct {
	shards [userShardCount]*userShard
}

// NewUserRepository initializes the shards
func NewUserRepository() *UserRepository {
	r := &UserRepository{}
	for i := 0; i < userShardCount; i++ {
		r.shards[i] = &userShard{
			users: make(map[uuid.UUID]*models.User),
		}
	}
	return r
}

// pickShard selects a shard based on userID
func (r *UserRepository) pickShard(userID uuid.UUID) *userShard {
	h := fnv.New32a()
	h.Write(userID[:])
	return r.shards[uint(h.Sum32())%userShardCount]
}

func (r *UserRepository) Create(user *models.User) error {
	shard := r.pickShard(user.ID)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, exists := shard.users[user.ID]; exists {
		return errors.ErrUserExists
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	shard.users[user.ID] = user
	return nil
}

func (r *UserRepository) GetByID(userID uuid.UUID) (*models.User, error) {
	shard := r.pickShard(userID)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	user, ok := shard.users[userID]
	if !ok {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	shard := r.pickShard(user.ID)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	existing, ok := shard.users[user.ID]
	if !ok {
		return errors.ErrUserNotFound
	}

	existing.Name = user.Name
	existing.Email = user.Email
	existing.UpdatedAt = time.Now()
	return nil
}

func (r *UserRepository) Delete(userID uuid.UUID) error {
	shard := r.pickShard(userID)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	if _, ok := shard.users[userID]; !ok {
		return errors.ErrUserNotFound
	}

	delete(shard.users, userID)
	return nil
}

func (r *UserRepository) List() []*models.User {
	result := make([]*models.User, 0)
	for i := 0; i < userShardCount; i++ {
		shard := r.shards[i]
		shard.mu.RLock()
		for _, user := range shard.users {
			result = append(result, user)
		}
		shard.mu.RUnlock()
	}
	return result
}