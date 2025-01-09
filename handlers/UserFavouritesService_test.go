package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/models"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #region MockUserDataService

// MockUserDataService is a mock implementation of UserDataService for testing purposes.
type MockUserDataService struct {
	mock.Mock
}

func (m *MockUserDataService) GetUser(userID string) (*models.User, error) {
	args := m.Called(userID)
	if args.Get(0) != nil {
		return args.Get(0).(*models.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserDataService) AddUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserDataService) AddAssetToFavourites(userID string, asset models.Asset) error {
	args := m.Called(userID, asset)
	return args.Error(0)
}

func (m *MockUserDataService) RemoveAssetFromFavourites(userID string, assetID string) error {
	args := m.Called(userID, assetID)
	return args.Error(0)
}

func (m *MockUserDataService) EditAssetDescription(userID string, assetID string, description string) error {
	args := m.Called(userID, assetID, description)
	return args.Error(0)
}

// #endregion MockUserDataService

// #region TestGetUserFavourites

func TestGetUserFavourites(t *testing.T) {
	mockDataService := new(MockUserDataService)
	handler := NewUserFavouritesService(mockDataService)

	t.Run("successful response", func(t *testing.T) {
		userID := "1"
		user := &models.User{
			ID: userID,
			Favourites: []models.Asset{
				{ID: "1", Description: "Asset 1"},
			},
		}
		mockDataService.On("GetUser", userID).Return(user, nil)

		req, err := http.NewRequest("GET", "/users/1/favourites", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/users/{id}/favourites", handler.GetUserFavourites).Methods("GET")
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var favourites []models.Asset
		err = json.NewDecoder(rr.Body).Decode(&favourites)
		assert.NoError(t, err)
		assert.Equal(t, user.Favourites, favourites)
		mockDataService.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := "2"
		mockDataService.On("GetUser", userID).Return(nil, errors.New("user not found"))

		req, err := http.NewRequest("GET", "/users/2/favourites", nil)
		assert.NoError(t, err)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/users/{id}/favourites", handler.GetUserFavourites).Methods("GET")
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		mockDataService.AssertExpectations(t)
	})
}

// #endregion TestGetUserFavourites
