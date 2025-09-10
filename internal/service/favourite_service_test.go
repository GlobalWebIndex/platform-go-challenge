package service

import (
	"context"
	"errors"
	"platform-go-challenge/internal/domain/entities"
	"platform-go-challenge/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// MockRepository is a mock implementation of the Repository interface for testing.
type MockRepository struct {
	repository.Repository
	mockAddFavourite               func(userID string, assetID string) error
	mockRemoveFavourite            func(userID string, assetID string) error
	mockUpdateFavouriteDescription func(userID string, assetID string, newDescription string) error
	mockGetFavourites              func(userID string, page, size int) ([]entities.Asset, error)
}

func (m *MockRepository) AddFavourite(userID string, assetID string) error {
	return m.mockAddFavourite(userID, assetID)
}

func (m *MockRepository) RemoveFavourite(userID string, assetID string) error {
	return m.mockRemoveFavourite(userID, assetID)
}

func (m *MockRepository) UpdateFavouriteDescription(userID string, assetID string, newDescription string) error {
	return m.mockUpdateFavouriteDescription(userID, assetID, newDescription)
}

func (m *MockRepository) GetFavourites(userID string, page, size int) ([]entities.Asset, error) {
	return m.mockGetFavourites(userID, page, size)
}

func TestGetFavourites(t *testing.T) {
	// Setup
	mockAssets := []entities.Asset{
		&entities.Audience{
			AssetBase: entities.AssetBase{
				ID:          "id1",
				Type:        "audience",
				Description: "a description",
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
			},
			Gender:         "male",
			BirthCountry:   "greece",
			AgeGroup:       "27",
			HoursOnSocial:  1,
			PurchasesMonth: 2,
		},
		&entities.Audience{
			AssetBase: entities.AssetBase{
				ID:          "id2",
				Type:        "audience",
				Description: "a description",
				CreatedAt:   time.Time{},
				UpdatedAt:   time.Time{},
			},
			Gender:         "female",
			BirthCountry:   "greece",
			AgeGroup:       "27",
			HoursOnSocial:  1,
			PurchasesMonth: 2,
		},
	}
	mockRepo := &MockRepository{
		mockGetFavourites: func(userID string, page, size int) ([]entities.Asset, error) {
			return mockAssets, nil
		},
	}
	logger := zap.NewNop()
	service := NewFavouriteService(mockRepo, logger)

	// Test case
	assets := service.GetFavourites(context.Background(), "user1", 1, 10)
	assert.Equal(t, mockAssets, assets)
}

// Separate test for a scenario where repository returns an error.
func TestGetFavouritesError(t *testing.T) {
	// Setup
	mockRepo := &MockRepository{
		mockGetFavourites: func(userID string, page, size int) ([]entities.Asset, error) {
			return nil, errors.New("repository error")
		},
	}
	logger := zap.NewNop()
	service := NewFavouriteService(mockRepo, logger)

	// Test case
	assets := service.GetFavourites(context.Background(), "user1", 0, 10)
	assert.Nil(t, assets)
}

// TestAddFavourite_Success tests the successful path of the AddFavourite method.
func TestAddFavourite_Success(t *testing.T) {
	// Mock Setup
	mockRepo := &MockRepository{
		mockAddFavourite: func(userID string, assetID string) error {
			return nil
		},
	}
	// Create service
	service := NewFavouriteService(mockRepo, zap.NewNop())

	// Act
	err := service.AddFavourite(context.Background(), "user1", "asset1")

	// Assert
	assert.NoError(t, err)
}

// TestAddFavourite_Error tests the error path of the AddFavourite method.
func TestAddFavourite_Error(t *testing.T) {
	// Mock Setup
	expectedErr := errors.New("database error")
	mockRepo := &MockRepository{
		mockAddFavourite: func(userID string, assetID string) error {
			return expectedErr
		},
	}

	//Create service
	service := NewFavouriteService(mockRepo, zap.NewNop())

	// Act
	err := service.AddFavourite(context.Background(), "user1", "asset1")

	// Assert
	assert.EqualError(t, err, expectedErr.Error())
}

// TestRemoveFavourite_Success tests the successful path of the RemoveFavourite method.
func TestRemoveFavourite_Success(t *testing.T) {
	// MockSetup
	mockRepo := &MockRepository{
		mockRemoveFavourite: func(userID string, assetID string) error {
			return nil
		},
	}
	// Create service
	service := NewFavouriteService(mockRepo, zap.NewNop())

	// Act
	err := service.RemoveFavourite(context.Background(), "user1", "asset1")

	// Assert
	assert.NoError(t, err)
}

// TestRemoveFavourite_Error tests the error path of the RemoveFavourite method.
func TestRemoveFavourite_Error(t *testing.T) {
	// Mock Setup
	expectedErr := errors.New("database error")
	mockRepo := &MockRepository{
		mockRemoveFavourite: func(userID string, assetID string) error {
			return expectedErr
		},
	}
	// Create service
	service := NewFavouriteService(mockRepo, zap.NewNop())

	// Act
	err := service.RemoveFavourite(context.Background(), "user1", "asset1")

	// Assert
	assert.EqualError(t, err, expectedErr.Error())
}

// TestUpdateFavouriteDescription_Success tests the successful path of the UpdateFavouriteDescription method.
func TestUpdateFavouriteDescription_Success(t *testing.T) {
	// Mock Setup
	mockRepo := &MockRepository{
		mockUpdateFavouriteDescription: func(userID string, assetID string, newDescription string) error {
			return nil
		},
	}
	// Create Service
	service := NewFavouriteService(mockRepo, zap.NewNop())

	// Act
	err := service.UpdateFavouriteDescription(context.Background(), "user1", "asset1", "new desc")

	// Assert.
	assert.NoError(t, err)
}

// TestUpdateFavouriteDescription_Error tests the error path of the UpdateFavouriteDescription method.
func TestUpdateFavouriteDescription_Error(t *testing.T) {
	// Mock Setup
	expectedErr := errors.New("database error")
	mockRepo := &MockRepository{
		mockUpdateFavouriteDescription: func(userID string, assetID string, newDescription string) error {
			return expectedErr
		},
	}
	// Create service.
	service := NewFavouriteService(mockRepo, zap.NewNop())

	// Act
	err := service.UpdateFavouriteDescription(context.Background(), "user1", "asset1", "new desc")

	// Assert
	assert.EqualError(t, err, expectedErr.Error())
}
