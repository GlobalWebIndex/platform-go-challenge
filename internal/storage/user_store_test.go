package storage

import (
	"strconv"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"platform-go-challenge/internal/common"
	"platform-go-challenge/internal/models"
)

func newTestUser(username, password string, role models.Role) models.User {
	user, _ := models.CreateUser(username, password, role)
	return user
}

func TestUserStore(t *testing.T) {
	store := NewUserStore()

	t.Run("AddUser and GetUserByUserID", func(t *testing.T) {
		user := newTestUser("alice", "pass123", models.RoleUser)
		err := store.AddUser(&user)
		assert.NoError(t, err)

		// Retrieve by ID
		gotUser, exists := store.GetUserByUserID(user.ID)
		assert.True(t, exists)
		assert.Equal(t, user.Username, gotUser.Username)
	})

	t.Run("AddUser duplicate username", func(t *testing.T) {
		user := newTestUser("alice", "pass456", models.RoleUser)
		err := store.AddUser(&user)
		assert.Error(t, err)
		assert.EqualError(t, err, common.ErrUsernameExists)
	})

	t.Run("GetUserByUsername", func(t *testing.T) {
		user, exists := store.GetUserByUsername("alice")
		assert.True(t, exists)
		assert.Equal(t, "alice", user.Username)

		_, exists = store.GetUserByUsername("nonexistent")
		assert.False(t, exists)
	})

	t.Run("GetUserID", func(t *testing.T) {
		userID, exists := store.GetUserID("alice")
		assert.True(t, exists)
		assert.NotEqual(t, uuid.Nil, userID)

		userID, exists = store.GetUserID("nonexistent")
		assert.False(t, exists)
		assert.Equal(t, uuid.Nil, userID)
	})

	t.Run("GetAllUsernames", func(t *testing.T) {
		usernames := store.GetAllUsernames()
		assert.Contains(t, usernames, "alice")
		assert.Len(t, usernames, 1)
	})

	t.Run("DeleteUser", func(t *testing.T) {
		user, _ := store.GetUserByUsername("alice")
		err := store.DeleteUser(user.ID)
		assert.NoError(t, err)

		// Ensure user is gone
		_, exists := store.GetUserByUserID(user.ID)
		assert.False(t, exists)

		_, exists = store.GetUserByUsername("alice")
		assert.False(t, exists)
	})

	t.Run("DeleteUser not found", func(t *testing.T) {
		err := store.DeleteUser(uuid.New())
		assert.Error(t, err)
		assert.EqualError(t, err, common.ErrUserNotFound)
	})
}

func TestUserStoreParallelWithWaitGroup(t *testing.T) {
	store := NewUserStore()
	numUsers := 50

	var wg sync.WaitGroup

	// Add users concurrently
	for i := 0; i < numUsers; i++ {
		i := i // capture loop variable
		wg.Add(1)
		go func() {
			defer wg.Done()
			username := "user" + strconv.Itoa(i)
			user := newTestUser(username, "password", models.RoleUser)
			err := store.AddUser(&user)
			assert.NoError(t, err)
		}()
	}

	// Wait for all adds to complete
	wg.Wait()

	// Verify all users were added
	allUsernames := store.GetAllUsernames()
	assert.Len(t, allUsernames, numUsers)

	// Delete users concurrently
	wg = sync.WaitGroup{}
	for _, id := range allUsernames {
		id := id // capture loop variable
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := store.DeleteUser(id)
			assert.NoError(t, err)
		}()
	}

	// Wait for all deletes to complete
	wg.Wait()

	// Verify store is empty
	finalUsers := store.GetAllUsernames()
	assert.Len(t, finalUsers, 0)
}
