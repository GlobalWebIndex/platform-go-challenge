package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/domain/entities"
	"platform-go-challenge/internal/middleware"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFavouriteService is a mock implementation of the FavouriteService.
type MockFavouriteService struct {
	mock.Mock
}

func (m *MockFavouriteService) GetFavourites(ctx context.Context, userID string, page, size int) []entities.Asset {
	args := m.Called(ctx, userID, page, size)
	return args.Get(0).([]entities.Asset)
}

func (m *MockFavouriteService) AddFavourite(ctx context.Context, userID string, assetID string) error {
	args := m.Called(ctx, userID, assetID)
	return args.Error(0)
}

func (m *MockFavouriteService) UpdateFavouriteDescription(ctx context.Context, userID string, assetID string, description string) error {
	args := m.Called(ctx, userID, assetID, description)
	return args.Error(0)
}

func (m *MockFavouriteService) RemoveFavourite(ctx context.Context, userID string, assetID string) error {
	args := m.Called(ctx, userID, assetID)
	return args.Error(0)
}

func TestGetFavourites_Success(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	// Create the handler and INJECT the mock service.
	handler := &FavouriteHandler{FavouriteService: mockService}

	// Mock the service call
	expectedAssets := []entities.Asset{
		&entities.Audience{},
		&entities.Audience{},
	}
	mockService.On("GetFavourites", mock.Anything, "user1", 1, 10).Return(expectedAssets)

	// Set up Gin
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Use the pre-configured handler directly.
	router.GET("/users/:userID/favourites", handler.GetFavourites)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/user1/favourites?page=1&size=10", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestGetFavourites_InvalidUserID tests the handler's response to an empty userID.
func TestGetFavourites_InvalidUserID(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.GET("/users/:userID/favourites", handler.GetFavourites)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users//favourites?page=1&size=10", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestAddFavourite_Success tests the happy path of AddFavourite.
func TestAddFavourite_Success(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}
	mockService.On("AddFavourite", mock.Anything, "user1", "asset1").Return(nil)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/users/:userID/favourites", handler.AddFavourite)

	jsonBody := strings.NewReader(`{"assetID": "asset1"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/user1/favourites", jsonBody)
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, `{"message":"favourite added successfully"}`, strings.TrimSpace(w.Body.String()))
}

// TestAddFavourite_ServiceError tests the error path of AddFavourite.
func TestAddFavourite_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	mockService.On("AddFavourite", mock.Anything, "user1", "asset1").Return(errors.New("service error"))

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.POST("/users/:userID/favourites", handler.AddFavourite)

	jsonBody := strings.NewReader(`{"assetID": "asset1"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/user1/favourites", jsonBody)
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestAddFavourite_WhenUserIdIsEmpty_ThrowBadRequest tests that empty user id leads to 400 response code.
func TestAddFavourite_WhenUserIdIsEmpty_ThrowBadRequest(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.POST("/users/:userID/favourites", handler.AddFavourite)

	jsonBody := strings.NewReader(`{"assetID": "asset1"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users//favourites", jsonBody)
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response middleware.APIError
	json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Equal(t, "userID is required", response.GetMessage())
	assert.Equal(t, "userID", response.GetField())
}

// TestAddFavourite_WhenUserIdIsEmpty_ThrowBadRequest tests that empty user id leads to 400 response code.
func TestAddFavourite_WhenAssetIdIsNotProvided_ThrowBadRequest(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.POST("/users/:userID/favourites", handler.AddFavourite)

	jsonBody := strings.NewReader(`{"flow": "asset1"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/usr1/favourites", jsonBody)
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response middleware.APIError
	json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Equal(t, "invalid json body", response.GetMessage())
	assert.Equal(t, "", response.GetField())
}

// TestUpdateFavouriteDescription_Success tests the happy path for updating a favourite's description.
func TestUpdateFavouriteDescription_Success(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	// Mock the service call to return nil (success).
	mockService.On("UpdateFavouriteDescription", mock.Anything, "user1", "asset1", "new description").Return(nil)

	// Set up Gin for testing.
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.PATCH("/users/:userID/favourites/:assetID", handler.UpdateFavouriteDescription)

	// Prepare the request payload.
	requestBody := []byte(`{"description": "new description"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/user1/favourites/asset1", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	// Unmarshal the response body.
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert the message from the unmarshaled JSON.
	assert.Equal(t, "favourite description updated successfully", response["message"])
	mockService.AssertExpectations(t)
}

// TestUpdateFavouriteDescription_InvalidJSON tests the handler's response to an invalid JSON payload.
func TestUpdateFavouriteDescription_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	// Set up Gin with the error handler middleware.
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.PATCH("/users/:userID/favourites/:assetID", handler.UpdateFavouriteDescription)

	// Prepare the invalid request.
	requestBody := strings.NewReader(`{"description": 123}`) // Invalid type for description
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/user1/favourites/asset1", requestBody)
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Unmarshal the response body into the APIError struct.
	var response middleware.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert the fields of the unmarshaled struct.
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "invalid json body", response.Message)
	assert.Equal(t, "", response.Field)

	// Ensure the service method was not called.
	mockService.AssertNotCalled(t, "UpdateFavouriteDescription", mock.Anything, "user1", "asset1", mock.Anything)
}

// TestUpdateFavouriteDescription_ServiceError tests the handler's response when the service layer returns an error.
func TestUpdateFavouriteDescription_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	// Mock the service call to return a specific error.
	serviceError := errors.New("database connection failed")
	mockService.On("UpdateFavouriteDescription", mock.Anything, "user1", "asset1", "new description").Return(serviceError)

	// Set up Gin with the error handler middleware.
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.PATCH("/users/:userID/favourites/:assetID", handler.UpdateFavouriteDescription)

	// Prepare the request.
	requestBody := []byte(`{"description": "new description"}`)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PATCH", "/users/user1/favourites/asset1", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Unmarshal the response body into the APIError struct.
	var response middleware.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Assert the fields of the unmarshaled struct.
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "database connection failed", response.Message)

	mockService.AssertExpectations(t)
}

// TestDeleteFavourite_Success tests the happy path for deleting a favourite.
func TestDeleteFavourite_Success(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	// Mock the service call to return nil (success).
	mockService.On("RemoveFavourite", mock.Anything, "user1", "asset1").Return(nil)

	// Set up Gin for testing.
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.DELETE("/users/:userID/favourites/:assetID", handler.DeleteFavourite)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/user1/favourites/asset1", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "favourite removed successfully", response["message"])
	mockService.AssertExpectations(t)
}

// TestDeleteFavourite_ServiceError tests the handler's response when the service layer returns an error.
func TestDeleteFavourite_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	// Mock the service call to return a specific error.
	serviceError := errors.New("database error")
	mockService.On("RemoveFavourite", mock.Anything, "user1", "asset1").Return(serviceError)

	// Set up Gin with the error handler middleware.
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.DELETE("/users/:userID/favourites/:assetID", handler.DeleteFavourite)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/user1/favourites/asset1", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response middleware.APIError
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.Code)
	assert.Equal(t, "database error", response.Message)
	mockService.AssertExpectations(t)
}

// TestDeleteFavourite_InvalidUserID tests the handler's response to an empty userID.
func TestDeleteFavourite_InvalidUserID(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	// Set up Gin with the error handler middleware.
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler())
	router.DELETE("/users/:userID/favourites/:assetID", handler.DeleteFavourite)

	w := httptest.NewRecorder()
	// Test with an empty userID.
	req, _ := http.NewRequest("DELETE", "/users//favourites/asset1", nil)

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response middleware.APIError
	json.Unmarshal([]byte(w.Body.String()), &response)
	assert.Equal(t, "userID is required", response.GetMessage())
	assert.Equal(t, "userID", response.GetField())
}
