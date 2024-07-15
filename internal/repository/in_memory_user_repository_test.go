package repository_test

import (
	"testing"

	"sync"

	"github.com/ceciivanov/platform-go-challenge/internal/models"
	"github.com/ceciivanov/platform-go-challenge/internal/repository"
	"github.com/stretchr/testify/assert"
)

func setup() *repository.InMemoryUserRepository {
	repo := repository.NewInMemoryUserRepository()
	repo.GenerateSampleUsers(1, 1) // Create 1 user with 1 asset
	return repo
}

func TestGenerateSampleUsers(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	repo.GenerateSampleUsers(2, 2)

	assert.Equal(t, 2, len(repo.Users))
	for _, user := range repo.Users {
		assert.Equal(t, 2, len(user.Favourites))
	}
}

func TestGetUserFavorites(t *testing.T) {
	repo := setup()

	// Test existing user
	favorites, err := repo.GetUserFavorites(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, favorites)

	// Test non-existing user
	_, err = repo.GetUserFavorites(999)
	assert.Error(t, err)
}

func TestAddUserFavorite(t *testing.T) {
	repo := setup()
	newAsset := models.Insight{
		ID:          2,
		Type:        models.InsightType,
		Description: "New Insight",
		Text:        "Some text",
	}

	// Test adding asset to existing user
	err := repo.AddUserFavorite(1, newAsset)
	assert.NoError(t, err)

	// Verify asset was added
	favorites, _ := repo.GetUserFavorites(1)
	assert.Equal(t, newAsset, favorites[2])

	// Test adding existing asset to user
	err = repo.AddUserFavorite(1, newAsset)
	assert.Error(t, err)

	// Test adding asset to non-existing user
	err = repo.AddUserFavorite(999, newAsset)
	assert.Error(t, err)
}

func TestDeleteUserFavorite(t *testing.T) {
	repo := setup()

	// Test deleting existing asset
	err := repo.DeleteUserFavorite(1, 1)
	assert.NoError(t, err)

	// Verify asset was deleted
	favorites, _ := repo.GetUserFavorites(1)
	assert.Empty(t, favorites)

	// Test deleting non-existing asset
	err = repo.DeleteUserFavorite(1, 999)
	assert.Error(t, err)

	// Test deleting asset from non-existing user
	err = repo.DeleteUserFavorite(999, 1)
	assert.Error(t, err)
}

func TestEditUserFavorite(t *testing.T) {
	repo := setup()
	editedAsset := models.Insight{
		ID:          1,
		Type:        models.InsightType,
		Description: "Edited Insight",
		Text:        "Edited text",
	}

	// Test editing existing asset
	err := repo.EditUserFavorite(1, 1, editedAsset)
	assert.NoError(t, err)

	// Verify asset was edited
	favorites, _ := repo.GetUserFavorites(1)
	assert.Equal(t, editedAsset, favorites[1])

	// Test editing non-existing asset
	err = repo.EditUserFavorite(1, 999, editedAsset)
	assert.Error(t, err)

	// Test editing asset for non-existing user
	err = repo.EditUserFavorite(999, 1, editedAsset)
	assert.Error(t, err)

	// Test editing asset with mismatched ID
	editedAsset.ID = 999
	err = repo.EditUserFavorite(1, 1, editedAsset)
	assert.Error(t, err)

	// Test editing asset with mismatched type
	editedAsset.ID = 1
	editedAsset.Type = models.AudienceType
	err = repo.EditUserFavorite(1, 1, editedAsset)
	assert.Error(t, err)
}

// TESTS FOR CONCURRENT OPERATIONS

// TestConcurrentGetUserFavorites tests getting a user's favorites concurrently
func TestConcurrentGetUserFavorites(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	repo.GenerateSampleUsers(10, 10)

	// 10 concurrent operations to get the same user's favorites
	var wg sync.WaitGroup
	numOperations := 10

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := repo.GetUserFavorites(1)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
}

// TestConcurrentAddUserFavorite tests adding a favorite asset to a user concurrently
func TestConcurrentAddUserFavorite(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	repo.GenerateSampleUsers(10, 10)

	userID := 1
	asset := models.Insight{
		ID:          20,
		Type:        models.InsightType,
		Description: "Test Insight",
		Text:        "Test text",
	}

	// 10 concurrent operations to add the same asset to the user
	var wg sync.WaitGroup
	numOperations := 10

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.AddUserFavorite(userID, asset)
			if err != nil && err.Error() != "asset already exists" {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	favorites, err := repo.GetUserFavorites(userID)
	assert.NoError(t, err)
	assert.Equal(t, 11, len(favorites)) // expect 10 existing assets + 1 new asset
}

// TestConcurrentDeleteUserFavorite tests deleting a favorite asset from a user concurrently
func TestConcurrentDeleteUserFavorite(t *testing.T) {
	repo := repository.NewInMemoryUserRepository()
	repo.GenerateSampleUsers(10, 10)

	userID := 1

	// 10 concurrent operations to delete the same asset from the user
	var wg sync.WaitGroup
	numOperations := 10

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.DeleteUserFavorite(userID, 1)
			if err != nil && err.Error() != "asset not found" {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	favorites, err := repo.GetUserFavorites(userID)
	assert.NoError(t, err)
	assert.Equal(t, 9, len(favorites)) // expect 10 existing assets - 1 deleted asset
}

// TestConcurrentEditUserFavorite tests editing a favorite asset from a user concurrently
func TestConcurrentEditUserFavorite(t *testing.T) {
	// create 1 user with 1 asset
	repo := repository.NewInMemoryUserRepository()
	repo.Users = map[int]models.User{
		1: {
			ID: 1,
			Favourites: map[int]models.Asset{
				1: &models.Insight{
					ID:          1,
					Type:        models.InsightType,
					Description: "Sample Insight",
					Text:        "Sample Insight Text",
				},
			},
		},
	}

	userID := 1
	assetID := 1
	editedAsset := models.Insight{
		ID:          assetID,
		Type:        models.InsightType,
		Description: "Edited Insight",
		Text:        "Edited text",
	}

	// 10 concurrent operations to edit the same asset from the user
	var wg sync.WaitGroup
	numOperations := 10

	for i := 0; i < numOperations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.EditUserFavorite(userID, assetID, editedAsset)
			if err != nil && err.Error() != "asset not found" {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()

	// get user's favorites and verify the asset was edited
	favorites, err := repo.GetUserFavorites(userID)
	assert.NoError(t, err)
	// expect 1 asset with the folowing edited description and text
	editedFavorite := favorites[assetID]
	assert.Equal(t, "Edited Insight", editedFavorite.GetDescription())
	assert.Equal(t, "Edited text", editedFavorite.(models.Insight).Text)

	// expect 1 asset in the user's favorites
	assert.Equal(t, 1, len(favorites))
}
