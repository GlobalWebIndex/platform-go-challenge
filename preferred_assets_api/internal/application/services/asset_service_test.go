package services_test

import (
	"errors"
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/services"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

// Mocks
type mockAssetServiceRepo struct {
	saveCalled   bool
	saveErr      error
	deleteCalled bool
	deleteErr    error
}

func (m *mockAssetServiceRepo) Save(asset entities.AssetEntity) (entities.AssetEntity, error) {
	m.saveCalled = true
	return nil, m.saveErr
}

func (m *mockAssetServiceRepo) Delete(id string) error {
	m.deleteCalled = true
	return m.deleteErr
}

func (m *mockAssetServiceRepo) GetByID(id string) (entities.AssetEntity, error) {
	return entities.AssetBaseEntity{}, nil
}
func (m *mockAssetServiceRepo) GetByIDs(ids []string) ([]entities.AssetEntity, error) {
	return nil, nil
}
func (m *mockAssetServiceRepo) GetAll() ([]entities.AssetEntity, error) { return nil, nil }
func (m *mockAssetServiceRepo) GetByType(assetType entities.AssetType) ([]entities.AssetEntity, error) {
	return nil, nil
}
func (m *mockAssetServiceRepo) Update(asset entities.AssetEntity) error { return nil }
func (m *mockAssetServiceRepo) Exists(id string) (bool, error)          { return false, nil }

func newValidInsight() *domain.Insight {
	return &domain.Insight{
		AssetBase: domain.AssetBase{
			ID:    "1",
			Type:  domain.AssetType(domain.AssetTypeInsight),
			Title: "Example Insight",
		},
		Text: "Valid insight text",
	}
}

// Tests

func TestCreateAsset_Success(t *testing.T) {
	// Arrange
	mockRepo := &mockAssetServiceRepo{}
	service := services.NewAssetService(mockRepo)
	asset := newValidInsight()

	// Act
	start := time.Now().UTC()
	_, err := service.CreateAsset(asset)
	end := time.Now().UTC()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mockRepo.saveCalled {
		t.Error("expected Save to be called")
	}

	// Check that CreatedAt was set properly
	createdAt := asset.GetCreatedAt()
	if createdAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
	if createdAt.Before(start) || createdAt.After(end) {
		t.Errorf("CreatedAt not within expected range: got %v, expected between %v and %v", createdAt, start, end)
	}
}

func TestCreateAsset_SaveFails(t *testing.T) {
	// Arrange
	mockRepo := &mockAssetServiceRepo{saveErr: errors.New("save failed")}
	service := services.NewAssetService(mockRepo)
	asset := newValidInsight()

	// Act
	_, err := service.CreateAsset(asset)

	// Assert
	if err == nil {
		t.Error("expected error when Save fails")
	}
}

func TestDeleteAsset_Success(t *testing.T) {
	// Arrange
	mockRepo := &mockAssetServiceRepo{}
	service := services.NewAssetService(mockRepo)

	// Act
	err := service.DeleteAsset("asset1")

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !mockRepo.deleteCalled {
		t.Error("expected Delete to be called")
	}
}

func TestDeleteAsset_DeleteFails(t *testing.T) {
	// Arrange
	mockRepo := &mockAssetServiceRepo{deleteErr: errors.New("delete failed")}
	service := services.NewAssetService(mockRepo)

	// Act
	err := service.DeleteAsset("asset1")

	// Assert
	if err == nil {
		t.Error("expected error when Delete fails")
	}
}
