package storage

import (
	"challenge/models"
	"errors"
	"runtime"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// memoryHandler is the access handler for the in-memory storage option
type memoryHandler struct {
	Block    map[string]map[string]models.Asset // UserID: AssetID: Asset{}
	Users    map[string]models.User
	Released bool
}

// createMemoryHandler initializes the memory structure for the users/assets
func createMemoryHander() *memoryHandler {
	newMemory := memoryHandler{}
	newMemory.Block = make(map[string]map[string]models.Asset)
	newMemory.Users = make(map[string]models.User)
	return &newMemory
}

// initPrivilegedUser
func (m *memoryHandler) initPrivilegedUser(user string, pass string) error {
	if m.Released {
		err := errors.New("storage memory not allocated")
		return err
	}
	if user == "" || pass == "" {
		return nil
	}
	// Encrypt the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pass), 8)
	if err != nil {
		return err
	}
	m.Users[user] = models.User{
		ID:         user,
		Password:   string(hashedPassword),
		Cookie:     "",
		Privileged: true,
		Created:    time.Now(),
	}
	return nil
}

// RetrieveUserAssets returns the favourited assets for the selected user
func (m *memoryHandler) RetrieveUserAssets(userID string) (bool, []models.Asset, error) {
	if m.Released {
		err := errors.New("storage resource is not available")
		return false, nil, err
	}
	var result []models.Asset
	if user, userExists := m.Block[userID]; userExists {
		for _, value := range user {
			result = append(result, value)
		}
		return true, result, nil
	} else {
		return false, nil, nil
	}
}

// AddAssetToUser adds the provided asset to the selected user if it doesn't exist
func (m *memoryHandler) AddAssetToUser(userID string, asset models.Asset) (bool, error) {
	if m.Released {
		err := errors.New("storage resource is not available")
		return false, err
	}
	if user, userExists := m.Block[userID]; userExists {
		if _, assetExists := user[asset.ID]; assetExists {
			return false, nil
		} else {
			m.Block[userID][asset.ID] = asset
			return true, nil
		}
	} else {
		return false, nil
	}
}

// UpdateUserAssetDescription updates the description of a favourited asset for
// the selected user
func (m *memoryHandler) UpdateUserAssetDescription(userID string, assetID string, description string) (bool, error) {
	if m.Released {
		err := errors.New("storage resource is not available")
		return false, err
	}
	user, userExists := m.Block[userID]
	if userExists {
		_, assetExists := user[assetID]
		if assetExists {
			newVal := m.Block[userID][assetID]
			newVal.Description = description
			m.Block[userID][assetID] = newVal
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, nil
	}
}

// RemoveUser deletes a user from storage
func (m *memoryHandler) RemoveUser(userID string) (bool, error) {
	if m.Released {
		err := errors.New("storage resource is not available")
		return false, err
	}
	_, userExists := m.Users[userID]
	if userExists {
		delete(m.Block, userID) // If the user does not exist this is a no-op
		delete(m.Users, userID) // If the user does not exist this is a no-op
		return true, nil
	} else {
		return false, nil
	}
}

// RemoveUserAsset removes an asset from the selected user's favourites
func (m *memoryHandler) RemoveUserAsset(userID string, assetID string) (bool, error) {
	if m.Released {
		err := errors.New("storage resource is not available")
		return false, err
	}
	if _, userExists := m.Block[userID]; userExists {
		_, assetExists := m.Block[userID][assetID]
		if !assetExists {
			return false, nil
		}
		delete(m.Block[userID], assetID) // If the asset does not exist this is a no-op
		return true, nil
	} else {
		err := errors.New("user does not exist")
		return false, err
	}
}

// CheckUserExistence checks if a user exists
func (m *memoryHandler) CheckUserExistence(userID string) (bool, error) {
	if m.Released {
		err := errors.New("storage memory not allocated")
		return false, err
	}
	// _, userExists := m.Block[userID]
	_, userExists := m.Users[userID]
	return userExists, nil
}

// CheckUserAssetExistence checks if an asset is favourited by the selected user
func (m *memoryHandler) CheckUserAssetExistence(userID string, assetID string) (bool, error) {
	if m.Released {
		err := errors.New("user check error")
		return false, err
	}
	user, userExists := m.Block[userID]
	if userExists {
		_, assetExists := user[assetID]
		return assetExists, nil
	} else {
		return false, nil
	}
}

// CreateNewUser adds a new user to the list of users
func (m *memoryHandler) CreateNewUser(user models.User) (bool, error) {
	var t time.Time
	if user.Created == t {
		user.Created = time.Now()
	}
	if m.Released {
		err := errors.New("storage memory not allocated")
		return false, err
	}
	_, check := m.Block[user.ID]
	if check {
		return false, nil
	}
	m.Block[user.ID] = make(map[string]models.Asset)
	m.Users[user.ID] = models.User{
		ID:         user.ID,
		Password:   user.Password,
		Cookie:     user.Cookie,
		Privileged: user.Privileged,
		Created:    user.Created,
	}
	return true, nil
}

// CheckAdmin checks if a user is privileged (exists, admin, error)
func (m *memoryHandler) CheckAdmin(id string) (bool, bool, error) {
	if m.Released {
		err := errors.New("storage resource is not allocated")
		return false, false, err
	}
	usr, exists := m.Users[id]
	if !exists {
		return false, false, nil
	}
	if usr.Privileged {
		return true, true, nil
	} else {
		return true, false, nil
	}
}

// UpdateUserToken updates the user record with the latest cookie
func (m *memoryHandler) UpdateUserToken(userID string, cookie string) (bool, error) {
	if m.Released {
		err := errors.New("storage resource is not allocated")
		return false, err
	}
	prevUserInstance := m.Users[userID]
	newUserInstance := models.User{
		ID:         prevUserInstance.ID,
		Password:   prevUserInstance.Password,
		Cookie:     cookie,
		Privileged: prevUserInstance.Privileged,
		Created:    prevUserInstance.Created,
	}
	_, exists := m.Users[userID]
	if !exists {
		return false, nil
	}
	m.Users[userID] = newUserInstance
	// m.Users[userID].Cookie=cookie
	return true, nil
}

// FetchUserToken retrieves the current user token
func (m *memoryHandler) FetchUserToken(userID string) (bool, string, error) {
	if m.Released {
		err := errors.New("storage resource is not allocated")
		return false, "", err
	}
	var cookie string
	val, userExists := m.Users[userID]
	if userExists {
		cookie = val.Cookie
	}
	// cookie := m.Users[userID].Cookie
	return userExists, cookie, nil
}

// FetchUserPass retrieves the encoded password for the selected user (if it exists)
func (m *memoryHandler) FetchUserPass(userID string) (bool, string, error) {
	var pass string
	if m.Released {
		err := errors.New("storage memory not allocated")
		return false, pass, err
	}
	val, userExists := m.Users[userID]
	if userExists {
		pass = val.Password
	}
	return userExists, pass, nil
}

// FetchAllUsers retrieves all the currently registered users
func (m *memoryHandler) FetchAllUsers() ([]models.User, error) {
	if m.Released {
		err := errors.New("storage resource is not available")
		return nil, err
	}
	var results []models.User
	for _, val := range m.Users {
		results = append(results, val)
	}
	return results, nil
}

// Release closes the connection to the selected storage and manually calls GC
// to free allocated memory
func (m *memoryHandler) Release() error {
	if m.Released {
		err := errors.New("storage resource is already released")
		return err
	}
	for k := range m.Block {
		delete(m.Block, k)
	}
	for k := range m.Users {
		delete(m.Users, k)
	}
	m.Released = true
	runtime.GC()
	return nil
}

// MockStorageData generate mock storage data (user and favourites).
func MockStorageData() (models.User, models.Asset) {
	newUser := models.User{
		ID:         "testuser",
		Password:   "testpassword",
		Cookie:     "cookie",
		Privileged: false,
		Created:    time.Now(),
	}
	newAsset := models.Asset{
		ID:          "asset1",
		Type:        "chart",
		Description: "example description",
		AssetData:   []byte("asdasd"),
		Added:       time.Now(),
		Modified:    time.Now(),
	}

	return newUser, newAsset
}
