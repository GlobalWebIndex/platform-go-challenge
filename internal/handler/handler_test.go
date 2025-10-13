package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"platform-go-challenge/internal/middleware"
	"platform-go-challenge/internal/model"
	"platform-go-challenge/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockErrorStore simulates store errors for testing
type mockErrorStore struct {
	storage.Store
	getFavoritesError error
	addFavoriteError  error
	removeFavoriteError error
	updateFavoriteError error
}

func (m *mockErrorStore) GetUserFavorites(userID string) ([]*model.Favorite, error) {
	if m.getFavoritesError != nil {
		return nil, m.getFavoritesError
	}
	return []*model.Favorite{}, nil
}

func (m *mockErrorStore) AddFavorite(userID string, asset *model.Asset) error {
	return m.addFavoriteError
}

func (m *mockErrorStore) RemoveFavorite(userID string, assetID string) error {
	return m.removeFavoriteError
}

func (m *mockErrorStore) UpdateFavoriteDescription(userID string, assetID string, description string) error {
	return m.updateFavoriteError
}

func TestHealthCheck(t *testing.T) {
	t.Run("returns health status", func(t *testing.T) {
		metrics := middleware.NewMetrics()

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()

		HealthCheck(metrics)(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	})
}

func TestGetFavorites(t *testing.T) {
	t.Run("successfully retrieves favorites for existing user", func(t *testing.T) {
		store := storage.NewMemoryStore()
		asset := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart},
				Title:     "Test Chart",
			},
		}
		store.AddFavorite("user-123", asset)

		req := httptest.NewRequest(http.MethodGet, "/users/user-123/favorites", nil)
		req.SetPathValue("id", "user-123")
		w := httptest.NewRecorder()

		GetFavorites(store)(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "user-123", response["user_id"])
		assert.Equal(t, float64(1), response["count"])
	})

	t.Run("returns empty array for user with no favorites", func(t *testing.T) {
		store := storage.NewMemoryStore()

		req := httptest.NewRequest(http.MethodGet, "/users/user-123/favorites", nil)
		req.SetPathValue("id", "user-123")
		w := httptest.NewRecorder()

		GetFavorites(store)(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, float64(0), response["count"])
	})

	t.Run("handles store errors", func(t *testing.T) {
		mockStore := &mockErrorStore{
			getFavoritesError: errors.New("database error"),
		}

		req := httptest.NewRequest(http.MethodGet, "/users/user-123/favorites", nil)
		req.SetPathValue("id", "user-123")
		w := httptest.NewRecorder()

		GetFavorites(mockStore)(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "Failed to retrieve favorites")
	})
}

func TestAddFavorite(t *testing.T) {
	tests := []struct {
		name           string
		payload        map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successfully adds Chart asset",
			payload: map[string]interface{}{
				"asset": map[string]interface{}{
					"chart": map[string]interface{}{
						"id":    "chart-1",
						"type":  "chart",
						"title": "Test Chart",
					},
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "successfully adds Insight asset",
			payload: map[string]interface{}{
				"asset": map[string]interface{}{
					"insight": map[string]interface{}{
						"id":      "insight-1",
						"type":    "insight",
						"title":   "Test Insight",
						"content": "Content here",
					},
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "successfully adds Audience asset",
			payload: map[string]interface{}{
				"asset": map[string]interface{}{
					"audience": map[string]interface{}{
						"id":        "audience-1",
						"type":      "audience",
						"gender":    "Male",
						"countries": []string{"US", "UK"},
						"age_min":   18,
						"age_max":   65,
					},
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "rejects invalid JSON",
			payload:        nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid JSON",
		},
		{
			name: "rejects empty asset",
			payload: map[string]interface{}{
				"asset": map[string]interface{}{},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Asset must have one of: chart, insight, or audience",
		},
		{
			name: "rejects missing asset ID",
			payload: map[string]interface{}{
				"asset": map[string]interface{}{
					"chart": map[string]interface{}{
						"type":  "chart",
						"title": "No ID",
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Asset ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := storage.NewMemoryStore()
			var body *bytes.Buffer
			
			if tt.payload != nil {
				data, _ := json.Marshal(tt.payload)
				body = bytes.NewBuffer(data)
			} else {
				body = bytes.NewBuffer([]byte("invalid json"))
			}

			req := httptest.NewRequest(http.MethodPost, "/users/user-123/favorites", body)
			req.SetPathValue("id", "user-123")
			w := httptest.NewRecorder()

			AddFavorite(store)(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			
			if tt.expectedError != "" {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response["error"], tt.expectedError)
			}
		})
	}

	t.Run("handles store error", func(t *testing.T) {
		mockStore := &mockErrorStore{
			addFavoriteError: errors.New("database error"),
		}

		payload := map[string]interface{}{
			"asset": map[string]interface{}{
				"chart": map[string]interface{}{
					"id":    "chart-1",
					"type":  "chart",
					"title": "Test",
				},
			},
		}
		data, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/users/user-123/favorites", bytes.NewBuffer(data))
		req.SetPathValue("id", "user-123")
		w := httptest.NewRecorder()

		AddFavorite(mockStore)(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "Failed to add favorite")
	})
}

func TestRemoveFavorite(t *testing.T) {
	t.Run("successfully removes existing favorite", func(t *testing.T) {
		store := storage.NewMemoryStore()
		asset := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart},
			},
		}
		store.AddFavorite("user-123", asset)

		req := httptest.NewRequest(http.MethodDelete, "/users/user-123/favorites/chart-1", nil)
		req.SetPathValue("id", "user-123")
		req.SetPathValue("assetId", "chart-1")
		w := httptest.NewRecorder()

		RemoveFavorite(store)(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "chart-1", response["asset_id"])
	})

	t.Run("returns 404 for non-existent user", func(t *testing.T) {
		store := storage.NewMemoryStore()

		req := httptest.NewRequest(http.MethodDelete, "/users/user-123/favorites/chart-1", nil)
		req.SetPathValue("id", "user-123")
		req.SetPathValue("assetId", "chart-1")
		w := httptest.NewRecorder()

		RemoveFavorite(store)(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 404 for non-existent favorite", func(t *testing.T) {
		store := storage.NewMemoryStore()
		asset := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart},
			},
		}
		store.AddFavorite("user-123", asset)

		req := httptest.NewRequest(http.MethodDelete, "/users/user-123/favorites/chart-999", nil)
		req.SetPathValue("id", "user-123")
		req.SetPathValue("assetId", "chart-999")
		w := httptest.NewRecorder()

		RemoveFavorite(store)(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("handles store error", func(t *testing.T) {
		mockStore := &mockErrorStore{
			removeFavoriteError: errors.New("database error"),
		}

		req := httptest.NewRequest(http.MethodDelete, "/users/user-123/favorites/chart-1", nil)
		req.SetPathValue("id", "user-123")
		req.SetPathValue("assetId", "chart-1")
		w := httptest.NewRecorder()

		RemoveFavorite(mockStore)(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "Failed to remove favorite")
	})
}

func TestEditFavorite(t *testing.T) {
	t.Run("successfully updates description", func(t *testing.T) {
		store := storage.NewMemoryStore()
		asset := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart, Description: "Old"},
			},
		}
		store.AddFavorite("user-123", asset)

		payload := map[string]string{"description": "New description"}
		data, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPatch, "/users/user-123/favorites/chart-1", bytes.NewBuffer(data))
		req.SetPathValue("id", "user-123")
		req.SetPathValue("assetId", "chart-1")
		w := httptest.NewRecorder()

		EditFavorite(store)(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		
		// Verify description was updated
		favorites, _ := store.GetUserFavorites("user-123")
		assert.Equal(t, "New description", favorites[0].Asset.GetBaseAsset().Description)
	})

	t.Run("rejects invalid JSON", func(t *testing.T) {
		store := storage.NewMemoryStore()

		req := httptest.NewRequest(http.MethodPatch, "/users/user-123/favorites/chart-1", bytes.NewBuffer([]byte("invalid")))
		req.SetPathValue("id", "user-123")
		req.SetPathValue("assetId", "chart-1")
		w := httptest.NewRecorder()

		EditFavorite(store)(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("returns 404 for non-existent user", func(t *testing.T) {
		store := storage.NewMemoryStore()

		payload := map[string]string{"description": "New"}
		data, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPatch, "/users/user-123/favorites/chart-1", bytes.NewBuffer(data))
		req.SetPathValue("id", "user-123")
		req.SetPathValue("assetId", "chart-1")
		w := httptest.NewRecorder()

		EditFavorite(store)(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("returns 404 for non-existent favorite", func(t *testing.T) {
		store := storage.NewMemoryStore()
		asset := &model.Asset{
			Chart: &model.Chart{
				BaseAsset: model.BaseAsset{ID: "chart-1", Type: model.AssetTypeChart},
			},
		}
		store.AddFavorite("user-123", asset)

		payload := map[string]string{"description": "New"}
		data, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPatch, "/users/user-123/favorites/chart-999", bytes.NewBuffer(data))
		req.SetPathValue("id", "user-123")
		req.SetPathValue("assetId", "chart-999")
		w := httptest.NewRecorder()

		EditFavorite(store)(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("handles store error", func(t *testing.T) {
		mockStore := &mockErrorStore{
			updateFavoriteError: errors.New("database error"),
		}

		payload := map[string]string{"description": "New"}
		data, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPatch, "/users/user-123/favorites/chart-1", bytes.NewBuffer(data))
		req.SetPathValue("id", "user-123")
		req.SetPathValue("assetId", "chart-1")
		w := httptest.NewRecorder()

		EditFavorite(mockStore)(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Contains(t, response["error"], "Failed to update favorite")
	})
}

// Integration test
func TestFavoriteWorkflow(t *testing.T) {
	t.Run("full workflow: add, get, edit, remove", func(t *testing.T) {
		store := storage.NewMemoryStore()

		// Add
		payload := map[string]interface{}{
			"asset": map[string]interface{}{
				"chart": map[string]interface{}{
					"id":    "chart-1",
					"type":  "chart",
					"title": "Test Chart",
				},
			},
		}
		data, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/users/user-123/favorites", bytes.NewBuffer(data))
		req.SetPathValue("id", "user-123")
		w := httptest.NewRecorder()
		AddFavorite(store)(w, req)
		require.Equal(t, http.StatusCreated, w.Code)

		// Get
		req = httptest.NewRequest(http.MethodGet, "/users/user-123/favorites", nil)
		req.SetPathValue("id", "user-123")
		w = httptest.NewRecorder()
		GetFavorites(store)(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// Edit
		editPayload := map[string]string{"description": "Updated"}
		data, _ = json.Marshal(editPayload)
		req = httptest.NewRequest(http.MethodPatch, "/users/user-123/favorites/chart-1", bytes.NewBuffer(data))
		req.SetPathValue("id", "user-123")
		req.SetPathValue("assetId", "chart-1")
		w = httptest.NewRecorder()
		EditFavorite(store)(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// Remove
		req = httptest.NewRequest(http.MethodDelete, "/users/user-123/favorites/chart-1", nil)
		req.SetPathValue("id", "user-123")
		req.SetPathValue("assetId", "chart-1")
		w = httptest.NewRecorder()
		RemoveFavorite(store)(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// Verify empty
		favorites, _ := store.GetUserFavorites("user-123")
		assert.Empty(t, favorites)
	})
}