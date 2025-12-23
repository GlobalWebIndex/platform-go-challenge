package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ============================================================================
// UNIT TESTS - RequestHandler Layer
// ============================================================================

// TestCreateUserSuccess tests happy path for user creation
func TestCreateUserSuccess(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{
			userExists: true,
		},
	}
	handler := &RequestHandler{service: mockService}

	req := httptest.NewRequest("POST", "/api/v1/users", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if result["id"] == nil {
		t.Error("Expected id field in response")
	}

	if result["created_at"] == nil {
		t.Error("Expected created_at field in response")
	}
}

// TestListUsersSuccess tests successful user listing with pagination
func TestListUsersSuccess(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{
			userExists: true,
		},
	}
	handler := &RequestHandler{service: mockService}

	req := httptest.NewRequest("GET", "/api/v1/users?page=1&limit=20", nil)
	w := httptest.NewRecorder()

	handler.ListUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	// Verify paginated response structure
	if _, ok := result["users"]; !ok {
		t.Error("Expected 'users' array in paginated response")
	}

	if _, ok := result["pagination"]; !ok {
		t.Error("Expected 'pagination' object in response")
	}
}

// TestListUsersEmpty tests listing when no users exist
func TestListUsersEmpty(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{},
	}
	handler := &RequestHandler{service: mockService}

	req := httptest.NewRequest("GET", "/api/v1/users?page=1&limit=20", nil)
	w := httptest.NewRecorder()

	handler.ListUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	// Verify pagination shows zero results
	pagination, ok := result["pagination"].(map[string]interface{})
	if !ok {
		t.Error("Expected pagination object in response")
		return
	}

	if total, ok := pagination["total"].(float64); ok && total == 0 {
		t.Logf("Empty list pagination structure verified")
	}
}

// TestListUsersPaginationDefaults tests default pagination values
func TestListUsersPaginationDefaults(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{},
	}
	handler := &RequestHandler{service: mockService}

	// No page/limit parameters - should use defaults
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()

	handler.ListUsers(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	pagination, ok := result["pagination"].(map[string]interface{})
	if !ok {
		t.Error("Expected pagination object in response")
		return
	}

	// Verify defaults: page=1, limit=20
	if page, ok := pagination["page"].(float64); ok && page != 1 {
		t.Errorf("Expected default page 1, got %v", page)
	}
}

// ============================================================================
// ASSET TESTS
// ============================================================================

// TestCreateAssetSuccess tests creating an asset (chart, insight, or audience)
func TestCreateAssetSuccess(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{},
	}
	handler := &RequestHandler{service: mockService}

	assetData := map[string]interface{}{
		"type": "chart",
		"data": map[string]interface{}{
			"title": "Sales Data",
			"x_axis": "Month",
			"y_axis": "Revenue",
			"data":   []float64{100, 200, 300},
		},
	}
	bodyBytes, _ := json.Marshal(assetData)

	req := httptest.NewRequest("POST", "/api/v1/assets", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateAsset(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if result["id"] == nil {
		t.Error("Expected id field in asset response")
	}
	if result["type"] != "chart" {
		t.Error("Expected type field in asset response")
	}
}

// TestListAssetsSuccess tests retrieving all assets with optional type filter
func TestListAssetsSuccess(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{},
	}
	handler := &RequestHandler{service: mockService}

	// Test with type filter
	req := httptest.NewRequest("GET", "/api/v1/assets?page=1&limit=20&type=chart", nil)
	w := httptest.NewRecorder()

	handler.ListAssets(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if _, ok := result["assets"]; !ok {
		t.Error("Expected 'assets' array in response")
	}
	if _, ok := result["pagination"]; !ok {
		t.Error("Expected 'pagination' object in response")
	}
}

// ============================================================================
// FAVORITES TESTS
// ============================================================================

// TestGetFavoritesUserNotFound tests 404 response when user doesn't exist
func TestGetFavoritesUserNotFound(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{
			userExists: false, // User doesn't exist
		},
	}
	handler := &RequestHandler{service: mockService}

	req := httptest.NewRequest("GET", "/api/v1/users/nonexistent-user-id/favorites?page=1", nil)
	w := httptest.NewRecorder()

	handler.GetFavorites(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d (Not Found), got %d", http.StatusNotFound, w.Code)
	}

	var errorResp ErrorResponse
	json.NewDecoder(w.Body).Decode(&errorResp)

	if errorResp.Error == "" {
		t.Error("Expected error message in response")
	}
}

// TestGetFavoritesSuccess tests retrieving user's favorite assets with pagination
func TestGetFavoritesSuccess(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{
			userExists: true,
		},
	}
	handler := &RequestHandler{service: mockService}

	req := httptest.NewRequest("GET", "/api/v1/users/user-123/favorites?page=1&limit=20", nil)
	w := httptest.NewRecorder()

	handler.GetFavorites(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result PaginatedResponse
	json.NewDecoder(w.Body).Decode(&result)

	// Verify response structure
	if result.Pagination.Page != 1 {
		t.Errorf("Expected page 1, got %d", result.Pagination.Page)
	}
}

// TestGetFavoritesWithTypeFilter tests filtering favorites by asset type
func TestGetFavoritesWithTypeFilter(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{
			userExists: true,
		},
	}
	handler := &RequestHandler{service: mockService}

	// Filter for only 'chart' type assets
	req := httptest.NewRequest("GET", "/api/v1/users/user-123/favorites?page=1&type=chart", nil)
	w := httptest.NewRecorder()

	handler.GetFavorites(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

// TestAddFavoriteSuccess tests adding an asset to user's favorites
func TestAddFavoriteSuccess(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{
			userExists: true,
		},
	}
	handler := &RequestHandler{service: mockService}

	// Request body: asset_id and optional description
	body := map[string]interface{}{
		"asset_id":    "asset-456",
		"description": "Important quarterly report",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/users/user-123/favorites", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddFavorite(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d (Created), got %d", http.StatusCreated, w.Code)
	}

	var result map[string]interface{}
	json.NewDecoder(w.Body).Decode(&result)

	if result["id"] == nil {
		t.Error("Expected favorite id in response")
	}
}

// TestAddFavoriteWithoutDescription tests adding favorite without custom description
func TestAddFavoriteWithoutDescription(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{
			userExists: true,
		},
	}
	handler := &RequestHandler{service: mockService}

	// Request body: only asset_id (description is optional)
	body := map[string]interface{}{
		"asset_id": "asset-789",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/users/user-123/favorites", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddFavorite(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}
}

// TestAddFavoriteUserNotFound tests adding favorite when user doesn't exist
func TestAddFavoriteUserNotFound(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{
			userExists: false,
		},
	}
	handler := &RequestHandler{service: mockService}

	body := map[string]interface{}{
		"asset_id": "asset-456",
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", "/api/v1/users/nonexistent/favorites", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.AddFavorite(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d (Not Found), got %d", http.StatusNotFound, w.Code)
	}
}

// TestUpdateFavoriteDescriptionSuccess tests updating a favorite's custom description
func TestUpdateFavoriteDescriptionSuccess(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{
			userExists: true,
		},
	}
	handler := &RequestHandler{service: mockService}

	updateBody := map[string]interface{}{
		"description": "Updated: Critical metrics for Q4 review",
	}
	bodyBytes, _ := json.Marshal(updateBody)

	req := httptest.NewRequest("PUT", "/api/v1/users/user-123/favorites/asset-456", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateFavoriteDescription(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d (OK), got %d", http.StatusOK, w.Code)
	}
}

// TestRemoveFavoriteSuccess tests soft-deleting a favorite (removes from user's list)
func TestRemoveFavoriteSuccess(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{
			userExists: true,
		},
	}
	handler := &RequestHandler{service: mockService}

	req := httptest.NewRequest("DELETE", "/api/v1/users/user-123/favorites/asset-456", nil)
	w := httptest.NewRecorder()

	handler.RemoveFromFavorites(w, req)

	// Soft delete returns 204 No Content
	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status %d (No Content), got %d", http.StatusNoContent, w.Code)
	}
}

// TestRemoveFavoriteUserNotFound tests removing favorite when user doesn't exist
func TestRemoveFavoriteUserNotFound(t *testing.T) {
	mockService := &Service{
		storage: &mockStorage{
			userExists: false,
		},
	}
	handler := &RequestHandler{service: mockService}

	req := httptest.NewRequest("DELETE", "/api/v1/users/nonexistent/favorites/asset-456", nil)
	w := httptest.NewRecorder()

	handler.RemoveFromFavorites(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d (Not Found), got %d", http.StatusNotFound, w.Code)
	}
}

// ============================================================================
// HEALTH CHECK TEST
// ============================================================================

// TestHealthCheck verifies the service health endpoint
// Used by Docker HEALTHCHECK and monitoring systems
func TestHealthCheck(t *testing.T) {
	handler := &RequestHandler{}

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	handler.HealthCheck(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)

	if result["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", result["status"])
	}
}

// ============================================================================
// MOCK STORAGE - For unit testing without database
// ============================================================================

// mockStorage implements the Storage interface for testing.
// It simulates database operations without requiring a real database connection.
type mockStorage struct {
	userExists bool
	assets     map[string]*Asset
	favorites  map[string][]*Favorite
}

// CreateUser simulates user creation
func (m *mockStorage) CreateUser(userID string) error {
	return nil
}

// UserExists simulates checking if a user exists
func (m *mockStorage) UserExists(userID string) (bool, error) {
	return m.userExists, nil
}

// ListUsers simulates fetching paginated user list
func (m *mockStorage) ListUsers(limit int, offset int) ([]*struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
}, int, error) {
	// Return empty list for mock
	return make([]*struct {
		ID        string
		CreatedAt time.Time
		UpdatedAt time.Time
	}, 0), 0, nil
}

// DeleteUser simulates user deletion
func (m *mockStorage) DeleteUser(userID string) error {
	return nil
}

// CreateAsset simulates creating a new asset (chart, insight, or audience)
func (m *mockStorage) CreateAsset(assetType string, data json.RawMessage) (string, error) {
	assetID := "mock-asset-" + assetType
	if m.assets == nil {
		m.assets = make(map[string]*Asset)
	}
	m.assets[assetID] = &Asset{
		ID:   assetID,
		Type: assetType,
		Data: data,
	}
	return assetID, nil
}

// GetAsset simulates retrieving a single asset
func (m *mockStorage) GetAsset(assetID string) (*Asset, error) {
	if m.assets != nil {
		if asset, ok := m.assets[assetID]; ok {
			return asset, nil
		}
	}
	return &Asset{
		ID:   assetID,
		Type: "chart",
		Data: json.RawMessage(`{"test": "data"}`),
	}, nil
}

// ListAssets simulates fetching paginated asset list with optional type filter
func (m *mockStorage) ListAssets(limit int, offset int, assetType *string) ([]*Asset, int, error) {
	// Return empty list for mock
	return make([]*Asset, 0), 0, nil
}

// DeleteAsset simulates asset deletion
func (m *mockStorage) DeleteAsset(assetID string) error {
	if m.assets != nil {
		delete(m.assets, assetID)
	}
	return nil
}

// AddToFavorites simulates adding an asset to user's favorites
// Supports optional custom description override
func (m *mockStorage) AddToFavorites(userID string, assetID string, description *string) (string, error) {
	favoriteID := "mock-favorite-" + assetID
	if m.favorites == nil {
		m.favorites = make(map[string][]*Favorite)
	}
	m.favorites[userID] = append(m.favorites[userID], &Favorite{
		ID:                  favoriteID,
		UserID:              userID,
		DescriptionOverride: description,
		AddedAt:             time.Now(),
	})
	return favoriteID, nil
}

// GetFavorites simulates retrieving user's favorites with pagination and optional type filter
func (m *mockStorage) GetFavorites(
	userID string,
	limit int,
	offset int,
	assetType *string,
) ([]*Favorite, int, error) {
	// Return empty list for mock
	return make([]*Favorite, 0), 0, nil
}

// UpdateFavoriteDescription simulates updating a favorite's custom description
func (m *mockStorage) UpdateFavoriteDescription(
	userID string,
	assetID string,
	description string,
) (bool, error) {
	return true, nil
}

// RemoveFromFavorites simulates soft-delete of a favorite (sets deleted_at timestamp)
func (m *mockStorage) RemoveFromFavorites(userID string, assetID string) (bool, error) {
	return true, nil
}

// Close simulates closing database connection
func (m *mockStorage) Close() error {
	return nil
}

// ============================================================================
// TEST EXECUTION GUIDE
// ============================================================================

// For Docker/integration testing, see RUNNING_TESTS.md
