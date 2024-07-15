package service_test

import (
	"testing"

	"github.com/ceciivanov/platform-go-challenge/internal/models"
	"github.com/ceciivanov/platform-go-challenge/internal/repository"
	"github.com/ceciivanov/platform-go-challenge/internal/service"
	"github.com/stretchr/testify/assert"
)

func setup() *service.UserService {
	repo := repository.NewInMemoryUserRepository()
	repo.GenerateSampleUsers(1, 1) // Create 1 user with 1 asset
	return service.NewUserService(repo)
}

func TestGetUserFavorites(t *testing.T) {
	s := setup()

	// Test existing user
	favorites, err := s.GetUserFavorites(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, favorites)

	// Test non-existing user
	_, err = s.GetUserFavorites(999)
	assert.Error(t, err)
}

func TestAddUserFavorite(t *testing.T) {
	s := setup()
	newAsset := models.Insight{
		ID:          2,
		Type:        models.InsightType,
		Description: "New Insight",
		Text:        "Some text",
	}

	// Test adding asset to existing user
	err := s.AddUserFavorite(1, newAsset)
	assert.NoError(t, err)

	// Test adding existing asset to user
	err = s.AddUserFavorite(1, newAsset)
	assert.Error(t, err)

	// Test adding asset to non-existing user
	err = s.AddUserFavorite(999, newAsset)
	assert.Error(t, err)
}

func TestDeleteUserFavorite(t *testing.T) {
	s := setup()

	// Test deleting existing asset
	err := s.DeleteUserFavorite(1, 1)
	assert.NoError(t, err)

	// Test deleting non-existing asset
	err = s.DeleteUserFavorite(1, 999)
	assert.Error(t, err)

	// Test deleting asset from non-existing user
	err = s.DeleteUserFavorite(999, 1)
	assert.Error(t, err)
}

func TestEditUserFavorite(t *testing.T) {
	s := setup()
	editedAsset := models.Insight{
		ID:          1,
		Type:        models.InsightType,
		Description: "Edited Insight",
		Text:        "Edited text",
	}

	// Test editing existing asset
	err := s.EditUserFavorite(1, 1, editedAsset)
	assert.NoError(t, err)

	// Test editing non-existing asset
	err = s.EditUserFavorite(1, 999, editedAsset)
	assert.Error(t, err)

	// Test editing asset for non-existing user
	err = s.EditUserFavorite(999, 1, editedAsset)
	assert.Error(t, err)

	// Test editing asset with mismatched ID
	editedAsset.ID = 999
	err = s.EditUserFavorite(1, 1, editedAsset)
	assert.Error(t, err)

	// Test editing asset with mismatched type
	editedAsset.ID = 1
	editedAsset.Type = models.AudienceType
	err = s.EditUserFavorite(1, 1, editedAsset)
	assert.Error(t, err)
}
