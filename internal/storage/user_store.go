package storage

import (
	"errors"
	"maps"
	"sync"

	"github.com/google/uuid"

	"platform-go-challenge/internal/common"
	"platform-go-challenge/internal/models"
)

// User storage with concurrency safety
type UserStore struct {
	users         map[uuid.UUID]models.User // Map of User IDs to Users
	usernameIndex map[string]uuid.UUID      // Map of Usernames to User IDs for quick lookups
	mutex         sync.RWMutex
}

func NewUserStore() *UserStore {
	return &UserStore{
		users:         make(map[uuid.UUID]models.User),
		usernameIndex: make(map[string]uuid.UUID),
	}
}

// ---------- User management methods ----------

func (us *UserStore) AddUser(user *models.User) error {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	// check if the username already exists
	if _, exists := us.usernameIndex[user.Username]; exists {
		return errors.New(common.ErrUsernameExists)
	}

	for {
		user.ID = uuid.New()
		if _, exists := us.users[user.ID]; !exists {
			break
		}
	}

	us.users[user.ID] = *user
	us.usernameIndex[user.Username] = user.ID

	return nil
}

func (us *UserStore) GetUserByUserID(userID uuid.UUID) (models.User, bool) {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	user, exists := us.users[userID]
	return user, exists
}

func (us *UserStore) GetUserByUsername(username string) (models.User, bool) {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	userID, exists := us.usernameIndex[username]
	if !exists {
		return models.User{}, false
	}

	user, exists := us.users[userID]
	return user, exists
}

func (us *UserStore) GetUserID(username string) (uuid.UUID, bool) {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	userID, exists := us.usernameIndex[username]
	return userID, exists
}

func (us *UserStore) GetAllUsernames() map[string]uuid.UUID {
	us.mutex.RLock()
	defer us.mutex.RUnlock()

	// Make a copy so the caller can't mutate the internal map
	result := make(map[string]uuid.UUID, len(us.usernameIndex))
	maps.Copy(result, us.usernameIndex)
	return result
}

func (us *UserStore) DeleteUser(userID uuid.UUID) error {
	us.mutex.Lock()
	defer us.mutex.Unlock()

	user, exists := us.users[userID]
	if !exists {
		return errors.New(common.ErrUserNotFound)
	}

	delete(us.users, userID)
	delete(us.usernameIndex, user.Username)

	return nil
}
