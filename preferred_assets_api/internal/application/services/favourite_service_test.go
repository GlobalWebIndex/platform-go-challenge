package services_test

import (
	"errors"
	"testing"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/services"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

// Mocks
type mockFavouriteRepo struct {
	existsResult bool
	addErr       error
	deleteErr    error
}

func (m *mockFavouriteRepo) Exists(userID, assetID string) (bool, error) {
	return m.existsResult, nil
}

func (m *mockFavouriteRepo) Add(f entities.FavouriteEntity) error {
	return m.addErr
}

func (m *mockFavouriteRepo) Delete(userID, assetID string) error {
	return m.deleteErr
}

func (m *mockFavouriteRepo) GetByUserID(userID string) ([]entities.FavouriteEntity, error) {
	return nil, nil
}

// --- Tests ---

func TestCreateFavourite_Success(t *testing.T) {
	// Arrange
	mockRepo := &mockFavouriteRepo{}
	service := services.NewFavouriteService(mockRepo)
	fav := domain.Favourite{UserID: "u1", AssetID: "a1"}

	// Act
	err := service.CreateFavourite(fav)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateFavourite_AlreadyExists(t *testing.T) {
	// Arrange
	mockRepo := &mockFavouriteRepo{existsResult: true}
	service := services.NewFavouriteService(mockRepo)
	fav := domain.Favourite{UserID: "u1", AssetID: "a1"}

	// Act
	err := service.CreateFavourite(fav)

	// Assert
	if err == nil {
		t.Error("expected error when favourite already exists")
	}
}

func TestCreateFavourite_AddFails(t *testing.T) {
	// Arrange
	mockRepo := &mockFavouriteRepo{addErr: errors.New("db failed")}
	service := services.NewFavouriteService(mockRepo)
	fav := domain.Favourite{UserID: "u1", AssetID: "a1"}

	// Act
	err := service.CreateFavourite(fav)

	// Assert
	if err == nil {
		t.Error("expected error when Add fails")
	}
}

func TestDeleteFavourite_Success(t *testing.T) {
	// Arrange
	mockRepo := &mockFavouriteRepo{}
	service := services.NewFavouriteService(mockRepo)

	// Act
	err := service.DeleteFavourite("u1", "a1")

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteFavourite_Fails(t *testing.T) {
	// Arrange
	mockRepo := &mockFavouriteRepo{deleteErr: errors.New("delete failed")}
	service := services.NewFavouriteService(mockRepo)

	// Act
	err := service.DeleteFavourite("u1", "a1")

	// Assert
	if err == nil {
		t.Error("expected error when Delete fails")
	}
}
