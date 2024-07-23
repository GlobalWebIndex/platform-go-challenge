package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arhsxro/platform-go-challenge/models"
	"github.com/stretchr/testify/assert"
)

type MockStore struct {
	GetUserFavoritesFunc  func(ctx context.Context, userID, filterType string, page, pageSize int) ([]models.Asset, error)
	AddFavoriteFunc       func(ctx context.Context, userID string, asset models.Asset) error
	RemoveFavoriteFunc    func(ctx context.Context, userID, assetId string) error
	UpdateDescriptionFunc func(ctx context.Context, userID, assetID, newDescription string) error
}

func (m *MockStore) GetUserFavorites(ctx context.Context, userID, filterType string, page, pageSize int) ([]models.Asset, error) {
	if m.GetUserFavoritesFunc != nil {
		return m.GetUserFavoritesFunc(ctx, userID, filterType, page, pageSize)
	}
	assets := []models.Asset{
		{ID: "1", Type: "Chart", Description: "Test Asset 1", Data: []byte(`{"title": "Chart 1"}`)},
		{ID: "2", Type: "Chart", Description: "Test Asset 2", Data: []byte(`{"text": "Insight 2"}`)},
	}
	return assets, nil
}

func (m *MockStore) GetUserFavoritesInvalidType(ctx context.Context, userID, filterType string, page, pageSize int) ([]models.Asset, error) {
	if filterType != "" && !isValidAssetType(filterType) {
		return nil, errors.New("invalid asset type")
	}

	return []models.Asset{}, nil
}

func (m *MockStore) GetUserFavoritesTimeout(ctx context.Context, userID, filterType string, page, pageSize int) ([]models.Asset, error) {
	time.Sleep(5 * time.Second)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		assets := []models.Asset{
			{ID: "1", Type: "Chart", Description: "Test Asset 1", Data: []byte(`{"title": "Chart 1"}`)},
			{ID: "2", Type: "Chart", Description: "Test Asset 2", Data: []byte(`{"text": "Insight 2"}`)},
		}
		return assets, nil
	}
}

func (m *MockStore) GetUserFavoritesQueryFailed(ctx context.Context, userID, filterType string, page, pageSize int) ([]models.Asset, error) {
	return nil, errors.New("database error: maximum connections reached")
}

func (m *MockStore) AddFavorite(ctx context.Context, userID string, asset models.Asset) error {
	if m.AddFavoriteFunc != nil {
		return m.AddFavoriteFunc(ctx, userID, asset)
	}
	return nil
}

func (m *MockStore) AddFavoriteTimeout(ctx context.Context, userID string, asset models.Asset) error {
	time.Sleep(5 * time.Second)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (m *MockStore) AddFavoriteQueryFailed(ctx context.Context, userID string, asset models.Asset) error {
	return errors.New("database error: maximum connections reached")
}

func (m *MockStore) RemoveFavorite(ctx context.Context, userID, assetID string) error {
	if m.RemoveFavoriteFunc != nil {
		return m.RemoveFavoriteFunc(ctx, userID, assetID)
	}
	return nil
}

func (m *MockStore) RemoveFavoriteTimeout(ctx context.Context, userID, assetID string) error {
	time.Sleep(5 * time.Second)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (m *MockStore) RemoveFavoriteQueryFailed(ctx context.Context, userID, assetID string) error {
	return errors.New("database error: maximum connections reached")
}

// Default mock implementation for UpdateDescription
func (m *MockStore) UpdateDescription(ctx context.Context, userID, assetID, newDescription string) error {
	if m.UpdateDescriptionFunc != nil {
		return m.UpdateDescriptionFunc(ctx, userID, assetID, newDescription)
	}
	return nil
}

func (m *MockStore) UpdateDescriptionTimeout(ctx context.Context, userID, assetID, newDescription string) error {
	time.Sleep(5 * time.Second)
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (m *MockStore) UpdateDescriptionQueryFailed(ctx context.Context, userID, assetID, newDescription string) error {
	return errors.New("database error: maximum connections reached")
}

func (m *MockStore) Close() error {

	return nil
}

func isValidAssetType(paramType string) bool {
	for _, validType := range models.ValidAssetTypes {
		if models.AssetType(paramType) == validType {
			return true
		}
	}
	log.Println("false")
	return false
}

func TestHandleGetFavorites_ValidRequestWithFilter(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	req, err := http.NewRequest("GET", "/favorites/test_user?type=Chart", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.HandleFunc("/favorites/{user_id}", api.HandleGetFavorites)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "status codes do not match")

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "content type is not JSON")

	expectedBody := `[{"id":"1","type":"Chart","description":"Test Asset 1","data":{"title": "Chart 1"}},{"id":"2","type":"Chart","description":"Test Asset 2","data":{"text": "Insight 2"}}]`
	assert.JSONEq(t, expectedBody, rr.Body.String(), "response body does not match expected JSON")
}

func TestHandleGetFavorites_ValidRequestWithoutFilter(t *testing.T) {
	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	req, err := http.NewRequest("GET", "/favorites/test_user", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.HandleFunc("/favorites/{user_id}", api.HandleGetFavorites)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "status codes do not match")

	assert.Equal(t, "application/json", rr.Header().Get("Content-Type"), "content type is not JSON")

	expectedBody := `[{"id":"1","type":"Chart","description":"Test Asset 1","data":{"title": "Chart 1"}},{"id":"2","type":"Chart","description":"Test Asset 2","data":{"text": "Insight 2"}}]`
	assert.JSONEq(t, expectedBody, rr.Body.String(), "response body does not match expected JSON")
}

func TestHandleGetFavorites_InvalidAssetType(t *testing.T) {
	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	mockStore.GetUserFavoritesFunc = mockStore.GetUserFavoritesInvalidType

	req, err := http.NewRequest("GET", "/favorites/test_user?type=InvalidType&page=1&pageSize=2", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.HandleFunc("/favorites/{user_id}", api.HandleGetFavorites)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "status codes do not match")

	assert.Equal(t, "invalid asset type\n", rr.Body.String())
}

func TestHandleGetFavorites_ContextRequestTimeout(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	mockStore.GetUserFavoritesFunc = mockStore.GetUserFavoritesTimeout

	req, err := http.NewRequest("GET", "/favorites/test_user?type=Chart", nil)
	if err != nil {
		t.Fatal(err)
	}

	//context with a timeout shorter than the sleep in GetUserFavoritesTimeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	router.HandleFunc("/favorites/{user_id}", api.HandleGetFavorites)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusGatewayTimeout, rr.Code, "status code is not Gateway Timeout")
}

func TestHandleGetFavorites_QueryFailed(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	mockStore.GetUserFavoritesFunc = mockStore.GetUserFavoritesQueryFailed

	req, err := http.NewRequest("GET", "/favorites/test_user?type=Chart&page=1&pageSize=2", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.HandleFunc("/favorites/{user_id}", api.HandleGetFavorites)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status codes do not match")
}

//Tests for HandleAddFavorite

func TestHandleAddFavorites_NormalFlow(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	requestBody := models.Asset{
		ID:          "insight2",
		Type:        "Insight",
		Description: "An insightful text",
		Data:        json.RawMessage(`{"text": "only 15% of the people in Greece watch One Piece"}`),
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/favorites/test_user", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	router.HandleFunc("/favorites/{user_id}", api.HandleAddFavorite)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "status codes do not match")
}

func TestHandleAddFavorites_InvalidPayload(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	invalidRequestBody := `{"id":123, "type": 456, "description": 789, "data": "not a valid json"}`

	req, err := http.NewRequest("POST", "/favorites/test_user", bytes.NewBufferString(invalidRequestBody))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	router.HandleFunc("/favorites/{user_id}", api.HandleAddFavorite)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "status codes do not match")

	assert.Equal(t, "Invalid request payload\n", rr.Body.String())

}

func TestHandleAddFavorites_ContextTimeout(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	mockStore.AddFavoriteFunc = mockStore.AddFavoriteTimeout

	requestBody := models.Asset{
		ID:          "insight2",
		Type:        "Insight",
		Description: "An insightful text",
		Data:        json.RawMessage(`{"text": "only 15% of the people in Greece watch One Piece"}`),
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/favorites/test_user", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req.Header.Set("Content-Type", "application/json")
	router.HandleFunc("/favorites/{user_id}", api.HandleAddFavorite)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusGatewayTimeout, rr.Code, "status code is not Gateway Timeout")
}

func TestHandleAddFavorites_QueryFailed(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	mockStore.AddFavoriteFunc = mockStore.AddFavoriteQueryFailed

	requestBody := models.Asset{
		ID:          "insight2",
		Type:        "Insight",
		Description: "An insightful text",
		Data:        json.RawMessage(`{"text": "only 15% of the people in Greece watch One Piece"}`),
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/favorites/test_user", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	router.HandleFunc("/favorites/{user_id}", api.HandleAddFavorite)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status codes do not match")
}

//Tests for RemoveFavorite Handler

func TestHandleRemoveFavorites_NormalFlow(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	req, err := http.NewRequest("DELETE", "/favorites/test_user/test_asset", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.HandleFunc("/favorites/{user_id}/{asset_id}", api.HandleRemoveFavorite)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "status codes do not match")
}

func TestHandleRemoveFavorites_ContextTimeout(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	mockStore.RemoveFavoriteFunc = mockStore.RemoveFavoriteTimeout

	req, err := http.NewRequest("DELETE", "/favorites/test_user/test_asset", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	router.HandleFunc("/favorites/{user_id}/{asset_id}", api.HandleRemoveFavorite)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusGatewayTimeout, rr.Code, "status codes do not match")
}

func TestHandleRemoveFavorites_QueryFailed(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	mockStore.RemoveFavoriteFunc = mockStore.RemoveFavoriteQueryFailed

	req, err := http.NewRequest("DELETE", "/favorites/test_user/test_asset", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.HandleFunc("/favorites/{user_id}/{asset_id}", api.HandleRemoveFavorite)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status codes do not match")
}

//Tests for EditDesciption Handler

func TestHandleEditDescription_NormalFlow(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	requestBody := `{"description": "Updated description for the asset"}`
	req, err := http.NewRequest("PUT", "/favorites/test_user/test_asset", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	router.HandleFunc("/favorites/{user_id}/{asset_id}", api.HandleEditDescription)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "status codes do not match")
}

func TestHandleEditDescription_InvalidPayload(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	invalidRequestBody := `{"description": "Updated description for the asset"`
	req, err := http.NewRequest("PUT", "/favorites/test_user/asset1", bytes.NewBufferString(invalidRequestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	router.HandleFunc("/favorites/{user_id}/{asset_id}", api.HandleEditDescription)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "status codes do not match")

	assert.Equal(t, "Invalid request payload\n", rr.Body.String())
}

func TestHandleEditDescription_ContextTimeout(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	mockStore.UpdateDescriptionFunc = mockStore.UpdateDescriptionTimeout

	requestBody := `{"description": "Updated description for the asset"}`
	req, err := http.NewRequest("PUT", "/favorites/test_user/test_asset", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	router.HandleFunc("/favorites/{user_id}/{asset_id}", api.HandleEditDescription)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusGatewayTimeout, rr.Code, "status code is not Gateway Timeout")
}

func TestHandleEditDescription_QueryFailed(t *testing.T) {

	mockStore := &MockStore{}
	api := InitApi(mockStore)
	router := api.InitRoutes()

	mockStore.UpdateDescriptionFunc = mockStore.UpdateDescriptionQueryFailed

	requestBody := `{"description": "Updated description for the asset"}`
	req, err := http.NewRequest("PUT", "/favorites/test_user/test_asset", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	router.HandleFunc("/favorites/{user_id}/{asset_id}", api.HandleEditDescription)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "status codes do not match")
}
