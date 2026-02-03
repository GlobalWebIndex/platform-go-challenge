package service

import (
	"context"
	"encoding/json"
	"errors"
	"iter"
	"log/slog"
	"testing"

	"github.com/google/uuid"

	"go-favorites-app/internal/core/domain/favorites"

	"github.com/stretchr/testify/mock"
)

// Mocks

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, asset favorites.Asset) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (favorites.Asset, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(favorites.Asset), args.Error(1)
}

func (m *MockRepository) FindAll(ctx context.Context, limit, offset int) (iter.Seq2[favorites.Asset, error], error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(iter.Seq2[favorites.Asset, error]), args.Error(1)
}

func (m *MockRepository) FindByUser(ctx context.Context, userID string, limit, offset int) (iter.Seq2[favorites.Asset, error], error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(iter.Seq2[favorites.Asset, error]), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) UpdateDescription(ctx context.Context, id, description string) (favorites.Asset, error) {
	args := m.Called(ctx, id, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(favorites.Asset), args.Error(1)
}

type MockCache struct {
	mock.Mock
}

func (m *MockCache) AddToSet(ctx context.Context, id string, score float64) error {
	args := m.Called(ctx, id, score)
	return args.Error(0)
}

func (m *MockCache) Set(ctx context.Context, id string, data []byte) error {
	args := m.Called(ctx, id, data)
	return args.Error(0)
}

func (m *MockCache) GetBatch(ctx context.Context, ids []string) (map[string][]byte, error) {
	args := m.Called(ctx, ids)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string][]byte), args.Error(1)
}

func (m *MockCache) GetIdsFromSet(ctx context.Context, start, stop int64) ([]string, error) {
	args := m.Called(ctx, start, stop)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockCache) Remove(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCache) Invalidate(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockEnricher struct {
	mock.Mock
}

func (m *MockEnricher) Enrich(ctx context.Context, asset favorites.Asset) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

// Helper to silence logs
type testWriter struct{}

func (tw *testWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func mustMarshal(v any) []byte {
	data, _ := json.Marshal(v)
	return data
}

func TestService_Save(t *testing.T) {
	repo := new(MockRepository)
	cache := new(MockCache)
	enricher := new(MockEnricher)
	logger := slog.New(slog.NewTextHandler(&testWriter{}, nil))

	t.Run("successful save", func(t *testing.T) {
		svc := NewService(repo, cache, enricher, logger)

		asset := favorites.Insight{
			BaseAsset: favorites.BaseAsset{ID: "1", Name: "Test", Type: favorites.AssetTypeInsight},
			Content:   "Knowledge",
		}

		repo.On("Save", mock.Anything, asset).Return(nil).Once()

		// Synchronous Enriched + Cache
		enricher.On("Enrich", mock.Anything, asset).Return(nil).Once()
		cache.On("AddToSet", mock.Anything, "1", mock.Anything).Return(nil).Once()
		cache.On("Set", mock.Anything, "1", mock.Anything).Return(nil).Once()

		err := svc.Save(context.Background(), asset)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		repo.AssertExpectations(t)
		enricher.AssertExpectations(t)
		cache.AssertExpectations(t)
	})

	t.Run("validation failure", func(t *testing.T) {
		svc := NewService(repo, cache, enricher, logger)

		asset := favorites.Insight{
			BaseAsset: favorites.BaseAsset{ID: "1", Name: "Test", Type: favorites.AssetTypeInsight},
		}

		err := svc.Save(context.Background(), asset)
		if err == nil {
			t.Error("expected validation error, got nil")
		}
	})

	t.Run("repo failure", func(t *testing.T) {
		svc := NewService(repo, cache, enricher, logger)

		asset := favorites.Insight{
			BaseAsset: favorites.BaseAsset{ID: "1", Name: "Test", Type: favorites.AssetTypeInsight},
			Content:   "Knowledge",
		}

		repo.On("Save", mock.Anything, asset).Return(errors.New("db error")).Once()

		err := svc.Save(context.Background(), asset)
		if err == nil {
			t.Error("expected repo error, got nil")
		}
	})
}

func TestService_FindByID(t *testing.T) {
	repo := new(MockRepository)
	cache := new(MockCache)
	enricher := new(MockEnricher)
	logger := slog.New(slog.NewTextHandler(&testWriter{}, nil))

	t.Run("cache hit", func(t *testing.T) {
		svc := NewService(repo, cache, enricher, logger)

		asset := favorites.Insight{
			BaseAsset: favorites.BaseAsset{ID: "1", Name: "Test", Type: favorites.AssetTypeInsight},
			Content:   "Knowledge",
		}

		cache.On("GetBatch", mock.Anything, []string{"1"}).Return(map[string][]byte{
			"1": mustMarshal(asset),
		}, nil).Once()

		found, err := svc.FindByID(context.Background(), "1")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if found.GetID() != "1" {
			t.Errorf("expected ID 1, got %s", found.GetID())
		}
	})

	t.Run("cache miss - read repair with enrichment", func(t *testing.T) {
		svc := NewService(repo, cache, enricher, logger)

		asset := favorites.Insight{
			BaseAsset: favorites.BaseAsset{ID: "2", Name: "Test Miss", Type: favorites.AssetTypeInsight},
			Content:   "Knowledge",
		}

		// Cache Miss
		cache.On("GetBatch", mock.Anything, []string{"2"}).Return(map[string][]byte{}, nil).Once()
		// Repo Hit
		repo.On("FindByID", mock.Anything, "2").Return(asset, nil).Once()

		// Read Repair Expectations (Synchronous)
		// We expect Enrich to be called
		enricher.On("Enrich", mock.Anything, asset).Return(nil).Once()
		// We expect Cache update
		cache.On("AddToSet", mock.Anything, "2", mock.Anything).Return(nil).Once()
		cache.On("Set", mock.Anything, "2", mock.Anything).Return(nil).Once()

		res, err := svc.FindByID(context.Background(), "2")
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if res.GetID() != "2" {
			t.Errorf("got id %s, want 2", res.GetID())
		}

		repo.AssertExpectations(t)
		enricher.AssertExpectations(t)
		cache.AssertExpectations(t)
	})
}

func TestService_FindAll(t *testing.T) {
	repo := new(MockRepository)
	cache := new(MockCache)
	enricher := new(MockEnricher)
	logger := slog.New(slog.NewTextHandler(&testWriter{}, nil))

	t.Run("cache hit FindAll", func(t *testing.T) {
		svc := NewService(repo, cache, enricher, logger)

		cache.On("GetIdsFromSet", mock.Anything, int64(0), int64(9)).Return([]string{"1"}, nil).Once()
		cache.On("GetBatch", mock.Anything, []string{"1"}).Return(map[string][]byte{
			"1": mustMarshal(favorites.Insight{
				BaseAsset: favorites.BaseAsset{ID: "1", Name: "Test", Type: favorites.AssetTypeInsight},
				Content:   "Knowledge",
			}),
		}, nil).Once()

		results, err := svc.FindAll(context.Background(), 10, 0)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		count := 0
		for range results {
			count++
		}
		if count != 1 {
			t.Errorf("expected 1 item, got %d", count)
		}
	})
}

func TestService_Delete(t *testing.T) {
	repo := new(MockRepository)
	cache := new(MockCache)
	enricher := new(MockEnricher)
	logger := slog.New(slog.NewTextHandler(&testWriter{}, nil))

	id := "1"
	userID := uuid.NewString()

	t.Run("successful delete", func(t *testing.T) {
		svc := NewService(repo, cache, enricher, logger)

		// Setup FindByID callback check
		repo.On("FindByID", mock.Anything, id).Return(favorites.Insight{
			BaseAsset: favorites.BaseAsset{ID: id, UserID: userID, Type: favorites.AssetTypeInsight},
		}, nil).Once()

		repo.On("Delete", mock.Anything, id).Return(nil).Once()
		cache.On("Remove", mock.Anything, id).Return(nil).Once()

		err := svc.Delete(context.Background(), id, userID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("delete forbidden", func(t *testing.T) {
		svc := NewService(repo, cache, enricher, logger)

		otherUser := uuid.NewString()

		// Return asset owned by different user
		repo.On("FindByID", mock.Anything, id).Return(favorites.Insight{
			BaseAsset: favorites.BaseAsset{ID: id, UserID: otherUser, Type: favorites.AssetTypeInsight},
		}, nil).Once()

		err := svc.Delete(context.Background(), id, userID)
		if err == nil {
			t.Error("expected forbidden error, got nil")
		}
	})
}

func TestService_UpdateDescription(t *testing.T) {
	repo := new(MockRepository)
	cache := new(MockCache)
	enricher := new(MockEnricher)
	logger := slog.New(slog.NewTextHandler(&testWriter{}, nil))

	id := "1"
	userID := uuid.NewString()
	newDesc := "new desc"

	t.Run("successful update", func(t *testing.T) {
		svc := NewService(repo, cache, enricher, logger)

		asset := favorites.Insight{
			BaseAsset: favorites.BaseAsset{ID: id, UserID: userID, Name: "Test", Type: favorites.AssetTypeInsight, Description: newDesc},
			Content:   "Knowledge",
		}

		// Ownership check
		repo.On("FindByID", mock.Anything, id).Return(favorites.Insight{
			BaseAsset: favorites.BaseAsset{ID: id, UserID: userID, Type: favorites.AssetTypeInsight},
		}, nil).Once()

		repo.On("UpdateDescription", mock.Anything, id, newDesc).Return(asset, nil).Once()
		// Expect cache invalidation
		cache.On("Remove", mock.Anything, id).Return(nil).Once()

		updated, err := svc.UpdateDescription(context.Background(), id, newDesc, userID)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if updated.GetID() != id {
			t.Errorf("expected ID %s, got %s", id, updated.GetID())
		}
	})
}
