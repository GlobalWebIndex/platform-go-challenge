package handlers

import (
	"bytes"
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

func (m *MockFavouriteService) GetFavourites(userID string, page, size int) ([]entities.Asset, error) {
	args := m.Called(userID, page, size)
	return args.Get(0).([]entities.Asset), nil
}

func (m *MockFavouriteService) AddFavourite(userID string, assetID string) error {
	args := m.Called(userID, assetID)
	return args.Error(0)
}

func (m *MockFavouriteService) UpdateFavouriteDescription(userID string, assetID string, description string) error {
	args := m.Called(userID, assetID, description)
	return args.Error(0)
}

func (m *MockFavouriteService) RemoveFavourite(userID string, assetID string) error {
	args := m.Called(userID, assetID)
	return args.Error(0)
}

// getTestContext creates a new Gin context for testing with a given user ID and role.
func getTestContext(userID, role string) *gin.Context {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, userID)
	c.Set(middleware.JWT_ROLE_KEY, role)
	return c
}

// TestGetFavourites_Success tests the happy path of GetFavourites.
func TestGetFavourites_Success(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	expectedAssets := []entities.Asset{
		&entities.Audience{},
		&entities.Audience{},
	}
	mockService.On("GetFavourites", "user1", 1, 10).Return(expectedAssets)

	// Set up a Gin context with a user and make the request.
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("GET", "/users/user1/favourites?page=1&size=10", nil)
	c.Params = gin.Params{gin.Param{Key: "userID", Value: "user1"}}

	// Act
	handler.GetFavourites(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	// You should also assert the response body here
}

// TestGetFavourites_Unauthorized tests that a non-admin user cannot view another user's favourites.
func TestGetFavourites_Unauthorized(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	// Set up a Gin context for an unauthorized user.
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("GET", "/users/user2/favourites", nil)
	c.Params = gin.Params{gin.Param{Key: "userID", Value: "user2"}}

	// Act
	handler.GetFavourites(c)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertNotCalled(t, "GetFavourites")
}

// TestGetFavourites_AsAdmin tests that an admin can view any user's favourites.
func TestGetFavourites_AsAdmin(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	expectedAssets := []entities.Asset{
		&entities.Audience{},
	}
	mockService.On("GetFavourites", "user2", 1, 10).Return(expectedAssets)

	// Set up a Gin context for an admin user.
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "admin_user_id")
	c.Set(middleware.JWT_ROLE_KEY, "admin")
	c.Request, _ = http.NewRequest("GET", "/users/user2/favourites", nil)
	c.Params = gin.Params{gin.Param{Key: "userID", Value: "user2"}}

	// Act
	handler.GetFavourites(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertCalled(t, "GetFavourites", "user2", 1, 10)
}

// TestAddFavourite_Success tests the happy path of AddFavourite.
func TestAddFavourite_Success(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}
	mockService.On("AddFavourite", "user1", "asset1").Return(nil)

	w := httptest.NewRecorder()
	jsonBody := strings.NewReader(`{"assetID": "asset1"}`)
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("POST", "/users/user1/favourites", jsonBody)
	c.Params = gin.Params{gin.Param{Key: "userID", Value: "user1"}}
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.AddFavourite(c)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
}

// TestAddFavourite_Unauthorized tests that a user cannot add a favourite for another user.
func TestAddFavourite_Unauthorized(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	w := httptest.NewRecorder()
	jsonBody := strings.NewReader(`{"assetID": "asset1"}`)
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("POST", "/users/user2/favourites", jsonBody)
	c.Params = gin.Params{gin.Param{Key: "userID", Value: "user2"}}
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.AddFavourite(c)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertNotCalled(t, "AddFavourite")
}

// TestAddFavourite_AsAdmin tests that an admin can add a favourite for any user.
func TestAddFavourite_AsAdmin(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}
	mockService.On("AddFavourite", "user2", "asset1").Return(nil)

	w := httptest.NewRecorder()
	jsonBody := strings.NewReader(`{"assetID": "asset1"}`)
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "admin_user_id")
	c.Set(middleware.JWT_ROLE_KEY, "admin")
	c.Request, _ = http.NewRequest("POST", "/users/user2/favourites", jsonBody)
	c.Params = gin.Params{gin.Param{Key: "userID", Value: "user2"}}
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.AddFavourite(c)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertCalled(t, "AddFavourite", "user2", "asset1")
}

// TestAddFavourite_ServiceError tests the error path of AddFavourite.
func TestAddFavourite_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}
	mockService.On("AddFavourite", "user1", "asset1").Return(errors.New("service error"))

	w := httptest.NewRecorder()
	jsonBody := strings.NewReader(`{"assetID": "asset1"}`)
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("POST", "/users/user1/favourites", jsonBody)
	c.Params = gin.Params{gin.Param{Key: "userID", Value: "user1"}}
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.AddFavourite(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestAddFavourite_WhenUserIdIsEmpty_ThrowBadRequest tests that empty user id leads to 400 response code.
func TestAddFavourite_WhenUserIdIsEmpty_ThrowBadRequest(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	w := httptest.NewRecorder()
	jsonBody := strings.NewReader(`{"assetID": "asset1"}`)
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("POST", "/users//favourites", jsonBody)
	c.Params = gin.Params{gin.Param{Key: "userID", Value: ""}}
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.AddFavourite(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestAddFavourite_WhenAssetIdIsNotProvided_ThrowBadRequest tests that empty user id leads to 400 response code.
func TestAddFavourite_WhenAssetIdIsNotProvided_ThrowBadRequest(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	w := httptest.NewRecorder()
	jsonBody := strings.NewReader(`{"flow": "asset1"}`)
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "usr1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("POST", "/users/usr1/favourites", jsonBody)
	c.Params = gin.Params{gin.Param{Key: "userID", Value: "usr1"}}
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.AddFavourite(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdateFavouriteDescription_Success tests the happy path for updating a favourite's description.
func TestUpdateFavouriteDescription_Success(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}
	mockService.On("UpdateFavouriteDescription", "user1", "asset1", "new description").Return(nil)

	w := httptest.NewRecorder()
	requestBody := []byte(`{"description": "new description"}`)
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("PATCH", "/users/user1/favourites/asset1", bytes.NewBuffer(requestBody))
	c.Params = gin.Params{
		gin.Param{Key: "userID", Value: "user1"},
		gin.Param{Key: "assetID", Value: "asset1"},
	}
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.UpdateFavouriteDescription(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestUpdateFavouriteDescription_Unauthorized tests that a non-admin user cannot update another user's favourite.
func TestUpdateFavouriteDescription_Unauthorized(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	w := httptest.NewRecorder()
	requestBody := []byte(`{"description": "new description"}`)
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("PATCH", "/users/user2/favourites/asset1", bytes.NewBuffer(requestBody))
	c.Params = gin.Params{
		gin.Param{Key: "userID", Value: "user2"},
		gin.Param{Key: "assetID", Value: "asset1"},
	}
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.UpdateFavouriteDescription(c)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertNotCalled(t, "UpdateFavouriteDescription")
}

// TestUpdateFavouriteDescription_AsAdmin tests that an admin can update any user's favourite.
func TestUpdateFavouriteDescription_AsAdmin(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}
	mockService.On("UpdateFavouriteDescription", "user2", "asset1", "new description").Return(nil)

	w := httptest.NewRecorder()
	requestBody := []byte(`{"description": "new description"}`)
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "admin_user_id")
	c.Set(middleware.JWT_ROLE_KEY, "admin")
	c.Request, _ = http.NewRequest("PATCH", "/users/user2/favourites/asset1", bytes.NewBuffer(requestBody))
	c.Params = gin.Params{
		gin.Param{Key: "userID", Value: "user2"},
		gin.Param{Key: "assetID", Value: "asset1"},
	}
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.UpdateFavouriteDescription(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertCalled(t, "UpdateFavouriteDescription", "user2", "asset1", "new description")
}

// TestUpdateFavouriteDescription_InvalidJSON tests the handler's response to an invalid JSON payload.
func TestUpdateFavouriteDescription_InvalidJSON(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	w := httptest.NewRecorder()
	requestBody := strings.NewReader(`{"description": 123}`)
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("PATCH", "/users/user1/favourites/asset1", requestBody)
	c.Params = gin.Params{
		gin.Param{Key: "userID", Value: "user1"},
		gin.Param{Key: "assetID", Value: "asset1"},
	}
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.UpdateFavouriteDescription(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdateFavouriteDescription_ServiceError tests the handler's response when the service layer returns an error.
func TestUpdateFavouriteDescription_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}
	serviceError := errors.New("database connection failed")
	mockService.On("UpdateFavouriteDescription", "user1", "asset1", "new description").Return(serviceError)

	w := httptest.NewRecorder()
	requestBody := []byte(`{"description": "new description"}`)
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("PATCH", "/users/user1/favourites/asset1", bytes.NewBuffer(requestBody))
	c.Params = gin.Params{
		gin.Param{Key: "userID", Value: "user1"},
		gin.Param{Key: "assetID", Value: "asset1"},
	}
	c.Request.Header.Set("Content-Type", "application/json")

	// Act
	handler.UpdateFavouriteDescription(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestDeleteFavourite_Success tests the happy path for deleting a favourite.
func TestDeleteFavourite_Success(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}
	mockService.On("RemoveFavourite", "user1", "asset1").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("DELETE", "/users/user1/favourites/asset1", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userID", Value: "user1"},
		gin.Param{Key: "assetID", Value: "asset1"},
	}

	// Act
	handler.DeleteFavourite(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

// TestDeleteFavourite_Unauthorized tests that a non-admin user cannot delete another user's favourite.
func TestDeleteFavourite_Unauthorized(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("DELETE", "/users/user2/favourites/asset1", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userID", Value: "user2"},
		gin.Param{Key: "assetID", Value: "asset1"},
	}

	// Act
	handler.DeleteFavourite(c)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertNotCalled(t, "RemoveFavourite")
}

// TestDeleteFavourite_AsAdmin tests that an admin can delete any user's favourite.
func TestDeleteFavourite_AsAdmin(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}
	mockService.On("RemoveFavourite", "user2", "asset1").Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "admin_user_id")
	c.Set(middleware.JWT_ROLE_KEY, "admin")
	c.Request, _ = http.NewRequest("DELETE", "/users/user2/favourites/asset1", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userID", Value: "user2"},
		gin.Param{Key: "assetID", Value: "asset1"},
	}

	// Act
	handler.DeleteFavourite(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertCalled(t, "RemoveFavourite", "user2", "asset1")
}

// TestDeleteFavourite_ServiceError tests the handler's response when the service layer returns an error.
func TestDeleteFavourite_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}
	serviceError := errors.New("database error")
	mockService.On("RemoveFavourite", "user1", "asset1").Return(serviceError)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("DELETE", "/users/user1/favourites/asset1", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userID", Value: "user1"},
		gin.Param{Key: "assetID", Value: "asset1"},
	}

	// Act
	handler.DeleteFavourite(c)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// TestDeleteFavourite_InvalidUserID tests the handler's response to an empty userID.
func TestDeleteFavourite_InvalidUserID(t *testing.T) {
	// Arrange
	mockService := new(MockFavouriteService)
	handler := &FavouriteHandler{FavouriteService: mockService}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set(middleware.JWT_USER_ID_KEY, "user1")
	c.Set(middleware.JWT_ROLE_KEY, "user")
	c.Request, _ = http.NewRequest("DELETE", "/users//favourites/asset1", nil)
	c.Params = gin.Params{
		gin.Param{Key: "userID", Value: ""},
		gin.Param{Key: "assetID", Value: "asset1"},
	}

	// Act
	handler.DeleteFavourite(c)

	// Assert
	assert.Equal(t, http.StatusForbidden, w.Code)
}
