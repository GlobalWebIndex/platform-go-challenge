package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"iter"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-favorites-app/internal/core/domain/favorites"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Save(ctx context.Context, asset favorites.Asset) error {
	args := m.Called(ctx, asset)
	return args.Error(0)
}

func (m *MockService) FindByID(ctx context.Context, id string) (favorites.Asset, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(favorites.Asset), args.Error(1)
}

func (m *MockService) FindAll(ctx context.Context, limit, offset int) (iter.Seq2[favorites.Asset, error], error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(iter.Seq2[favorites.Asset, error]), args.Error(1)
}

func (m *MockService) FindAllByUser(ctx context.Context, userID string, limit, offset int) (iter.Seq2[favorites.Asset, error], error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(iter.Seq2[favorites.Asset, error]), args.Error(1)
}

func (m *MockService) Delete(ctx context.Context, id, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func (m *MockService) UpdateDescription(ctx context.Context, id, description, userID string) (favorites.Asset, error) {
	args := m.Called(ctx, id, description, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(favorites.Asset), args.Error(1)
}

func (m *MockService) Shutdown() {
	m.Called()
}

func TestHandler_Create(t *testing.T) {
	mockSvc := new(MockService)
	logger := slog.Default()
	h := NewHandler(mockSvc, logger)

	t.Run("success", func(t *testing.T) {
		id := uuid.NewString()
		userID := uuid.NewString()
		reqBody := map[string]interface{}{
			"type": "audience",
			"id":   id,
			"name": "Test Audience",
			"rules": map[string]interface{}{
				"gender": "male",
			},
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/favorites", bytes.NewBuffer(body))

		// Inject UserID into context
		ctx := context.WithValue(req.Context(), userIDKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSvc.On("Save", mock.Anything, mock.MatchedBy(func(a favorites.Asset) bool {
			return a.GetID() == id &&
				a.GetType() == favorites.AssetTypeAudience &&
				a.GetUserID() == userID
		})).Return(nil)

		h.Create(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/favorites", nil)
		w := httptest.NewRecorder()
		h.Create(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestHandler_Get(t *testing.T) {
	mockSvc := new(MockService)
	logger := slog.Default()
	h := NewHandler(mockSvc, logger)

	t.Run("success", func(t *testing.T) {
		id := uuid.NewString()
		asset := &favorites.Audience{
			BaseAsset: favorites.BaseAsset{ID: id, Name: "Found", Type: favorites.AssetTypeAudience},
			Rules:     favorites.AudienceRules{Gender: "female"},
		}

		req := httptest.NewRequest(http.MethodGet, "/favorites/"+id, nil)
		req.SetPathValue("id", id)
		w := httptest.NewRecorder()

		mockSvc.On("FindByID", mock.Anything, id).Return(asset, nil)

		h.Get(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestHandler_Delete(t *testing.T) {
	mockSvc := new(MockService)
	logger := slog.Default()
	h := NewHandler(mockSvc, logger)
	id := uuid.NewString()
	userID := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/favorites/"+id, nil)
		req.SetPathValue("id", id)

		ctx := context.WithValue(req.Context(), userIDKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		mockSvc.On("Delete", mock.Anything, id, userID).Return(nil)

		h.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestHandler_UpdateDescription(t *testing.T) {
	mockSvc := new(MockService)
	logger := slog.Default()
	handler := NewHandler(mockSvc, logger)
	// We don't need NewRouter for unit testing handlers generally, but if we used it we need to bypass auth
	// Direct call is easier

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()
		userID := uuid.NewString()
		desc := "Updated Description"
		reqBody := map[string]string{"description": desc}
		body, _ := json.Marshal(reqBody)

		expectedAsset := favorites.Chart{
			BaseAsset: favorites.BaseAsset{
				ID:          id,
				UserID:      userID,
				Name:        "My Chart",
				Type:        favorites.AssetTypeChart,
				Description: desc,
			},
			XAxis: "time",
			YAxis: "value",
		}

		mockSvc.On("UpdateDescription", mock.Anything, id, desc, userID).Return(expectedAsset, nil).Once()

		req, _ := http.NewRequest("PATCH", "/favorites/"+id, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", id)

		ctx := context.WithValue(req.Context(), userIDKey, userID)
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		handler.UpdateDescription(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var respAsset favorites.Chart
		err := json.Unmarshal(w.Body.Bytes(), &respAsset)
		assert.NoError(t, err)
		assert.Equal(t, id, respAsset.ID)
		assert.Equal(t, desc, respAsset.Description)
	})
}
