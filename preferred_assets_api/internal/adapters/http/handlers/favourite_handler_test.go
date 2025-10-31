package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/http/middleware"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/ports"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFavouriteService is a mock implementation of ports.FavouriteService for testing
type MockFavouriteService struct {
	mock.Mock
}

func (m *MockFavouriteService) CreateFavourite(favourite domain.Favourite) error {
	args := m.Called(favourite)
	return args.Error(0)
}

func (m *MockFavouriteService) DeleteFavourite(userID string, assetID string) error {
	args := m.Called(userID, assetID)
	return args.Error(0)
}

func TestFavouriteHandler_Create(t *testing.T) {
	// Save the original Body getter and restore after test
	originalBodyGetter := middleware.Body
	defer func() {
		middleware.Body = originalBodyGetter
	}()

	tests := []struct {
		name                string
		method              string
		requestBody         interface{}
		setupMock           func(*MockFavouriteService)
		expectedStatus      int
		expectedBody        string
		validateBodySucceed bool
	}{
		{
			name:   "Happy Path - Successfully creates favourite",
			method: http.MethodPost,
			requestBody: dto.FavouriteRequest{
				UserId:  "user-123",
				AssetId: "asset-456",
			},
			setupMock: func(m *MockFavouriteService) {
				expectedFavourite := domain.Favourite{
					UserID:  "user-123",
					AssetID: "asset-456",
				}
				m.On("CreateFavourite", expectedFavourite).Return(nil)
			},
			expectedStatus:      http.StatusCreated,
			expectedBody:        "",
			validateBodySucceed: true,
		},
		{
			name:           "Unhappy Path - Wrong HTTP method",
			method:         http.MethodGet,
			requestBody:    nil,
			setupMock:      func(m *MockFavouriteService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
		{
			name:           "Unhappy Path - Missing validated body",
			method:         http.MethodPost,
			requestBody:    nil,
			setupMock:      func(m *MockFavouriteService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing validated body\n",
		},
		{
			name:   "Unhappy Path - Service returns error",
			method: http.MethodPost,
			requestBody: dto.FavouriteRequest{
				UserId:  "user-123",
				AssetId: "asset-456",
			},
			setupMock: func(m *MockFavouriteService) {
				expectedFavourite := domain.Favourite{
					UserID:  "user-123",
					AssetID: "asset-456",
				}
				m.On("CreateFavourite", expectedFavourite).Return(errors.New("database error"))
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedBody:        "database error\n",
			validateBodySucceed: true,
		},
		{
			name:   "Unhappy Path - Empty user ID in request",
			method: http.MethodPost,
			requestBody: dto.FavouriteRequest{
				UserId:  "",
				AssetId: "asset-456",
			},
			setupMock:           func(m *MockFavouriteService) {},
			expectedStatus:      http.StatusBadRequest,
			expectedBody:        "missing validated body\n", // This depends on your validation middleware
			validateBodySucceed: false,
		},
		{
			name:   "Unhappy Path - Empty asset ID in request",
			method: http.MethodPost,
			requestBody: dto.FavouriteRequest{
				UserId:  "user-123",
				AssetId: "",
			},
			setupMock:           func(m *MockFavouriteService) {},
			expectedStatus:      http.StatusBadRequest,
			expectedBody:        "missing validated body\n", // This depends on your validation middleware
			validateBodySucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockFavouriteService)
			tt.setupMock(mockService)
			handler := NewFavouriteHandler(mockService)

			var body []byte
			if tt.requestBody != nil {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(tt.method, "/favourites", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			// Set up the mock middleware based on whether we have a request body
			if tt.requestBody != nil {
				middleware.Body = MockBodyGetter{
					MockedBody:    tt.requestBody,
					ShouldSucceed: tt.validateBodySucceed,
				}
			} else {
				middleware.Body = MockBodyGetter{
					MockedBody:    nil,
					ShouldSucceed: false,
				}
			}

			// Act
			handler.Create(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestFavouriteHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		userID         string
		assetID        string
		setupMock      func(*MockFavouriteService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "Happy Path - Successfully deletes favourite",
			method:  http.MethodDelete,
			userID:  "user-123",
			assetID: "asset-456",
			setupMock: func(m *MockFavouriteService) {
				m.On("DeleteFavourite", "user-123", "asset-456").Return(nil)
			},
			expectedStatus: http.StatusOK, // Note: Your handler doesn't set status, so it defaults to 200
			expectedBody:   "",
		},
		{
			name:           "Unhappy Path - Wrong HTTP method",
			method:         http.MethodPost,
			userID:         "user-123",
			assetID:        "asset-456",
			setupMock:      func(m *MockFavouriteService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
		{
			name:           "Unhappy Path - Missing user ID",
			method:         http.MethodDelete,
			userID:         "",
			assetID:        "asset-456",
			setupMock:      func(m *MockFavouriteService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing user id\n",
		},
		{
			name:           "Unhappy Path - Missing asset ID",
			method:         http.MethodDelete,
			userID:         "user-123",
			assetID:        "",
			setupMock:      func(m *MockFavouriteService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing asset id\n",
		},
		{
			name:    "Unhappy Path - Service returns not found error",
			method:  http.MethodDelete,
			userID:  "user-999",
			assetID: "asset-999",
			setupMock: func(m *MockFavouriteService) {
				m.On("DeleteFavourite", "user-999", "asset-999").Return(errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "user not found\n",
		},
		{
			name:    "Unhappy Path - Service returns other error",
			method:  http.MethodDelete,
			userID:  "user-123",
			assetID: "asset-456",
			setupMock: func(m *MockFavouriteService) {
				m.On("DeleteFavourite", "user-123", "asset-456").Return(errors.New("database error"))
			},
			expectedStatus: http.StatusNotFound, // Your handler returns 404 for all errors
			expectedBody:   "user not found\n",
		},
		{
			name:           "Unhappy Path - Both user ID and asset ID missing",
			method:         http.MethodDelete,
			userID:         "",
			assetID:        "",
			setupMock:      func(m *MockFavouriteService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing user id\n", // First check fails on user ID
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockFavouriteService)
			tt.setupMock(mockService)
			handler := NewFavouriteHandler(mockService)

			req := httptest.NewRequest(tt.method, "/favourites/"+tt.userID+"/"+tt.assetID, nil)
			rr := httptest.NewRecorder()

			// Set up chi router context for URL parameters
			rctx := chi.NewRouteContext()
			if tt.userID != "" {
				rctx.URLParams.Add("userId", tt.userID)
			}
			if tt.assetID != "" {
				rctx.URLParams.Add("assetId", tt.assetID)
			}
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Act
			handler.Delete(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestNewFavouriteHandler(t *testing.T) {
	t.Run("should create new favourite handler", func(t *testing.T) {
		// Arrange
		mockService := new(MockFavouriteService)

		// Act
		handler := NewFavouriteHandler(mockService)

		// Assert
		assert.NotNil(t, handler)
		assert.Equal(t, mockService, handler.service)
	})
}

// Test interface implementation
func TestFavouriteHandler_InterfaceImplementation(t *testing.T) {
	// Arrange & Act
	var handler ports.FavouriteHandler = NewFavouriteHandler(nil)

	// Assert
	assert.NotNil(t, handler)
}

// Test specific scenarios with detailed assertions
func TestFavouriteHandler_Create_Detailed(t *testing.T) {
	// Save the original Body getter and restore after test
	originalBodyGetter := middleware.Body
	defer func() {
		middleware.Body = originalBodyGetter
	}()

	t.Run("should call service with correct domain favourite", func(t *testing.T) {
		// Arrange
		mockService := new(MockFavouriteService)
		handler := NewFavouriteHandler(mockService)

		requestBody := dto.FavouriteRequest{
			UserId:  "test-user",
			AssetId: "test-asset",
		}

		// Expect the exact domain favourite that mapping should create
		expectedFavourite := domain.Favourite{
			UserID:  "test-user",
			AssetID: "test-asset",
		}

		mockService.On("CreateFavourite", expectedFavourite).Return(nil)
		middleware.Body = MockBodyGetter{
			MockedBody:    requestBody,
			ShouldSucceed: true,
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/favourites", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		// Act
		handler.Create(rr, req)

		// Assert
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Empty(t, rr.Body.String())
		mockService.AssertCalled(t, "CreateFavourite", expectedFavourite)
	})

	t.Run("should handle service validation errors", func(t *testing.T) {
		// Arrange
		mockService := new(MockFavouriteService)
		handler := NewFavouriteHandler(mockService)

		requestBody := dto.FavouriteRequest{
			UserId:  "user-123",
			AssetId: "asset-456",
		}

		mockService.On("CreateFavourite", mock.AnythingOfType("domain.Favourite")).
			Return(errors.New("duplicate favourite"))

		middleware.Body = MockBodyGetter{
			MockedBody:    requestBody,
			ShouldSucceed: true,
		}

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest(http.MethodPost, "/favourites", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		// Act
		handler.Create(rr, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, "duplicate favourite\n", rr.Body.String())
		mockService.AssertExpectations(t)
	})
}

func TestFavouriteHandler_Delete_EdgeCases(t *testing.T) {
	t.Run("should handle special characters in IDs", func(t *testing.T) {
		// Arrange
		mockService := new(MockFavouriteService)
		handler := NewFavouriteHandler(mockService)

		userID := "user-123-with-special-chars"
		assetID := "asset-456/special/path"

		mockService.On("DeleteFavourite", userID, assetID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/favourites/"+userID+"/"+assetID, nil)
		rr := httptest.NewRecorder()

		// Set up chi router context for URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("userId", userID)
		rctx.URLParams.Add("assetId", assetID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		handler.Delete(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code)
		mockService.AssertCalled(t, "DeleteFavourite", userID, assetID)
	})

	t.Run("should handle very long IDs", func(t *testing.T) {
		// Arrange
		mockService := new(MockFavouriteService)
		handler := NewFavouriteHandler(mockService)

		userID := "user-" + strings.Repeat("a", 1000)   // Very long user ID
		assetID := "asset-" + strings.Repeat("b", 1000) // Very long asset ID

		mockService.On("DeleteFavourite", userID, assetID).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/favourites/"+userID+"/"+assetID, nil)
		rr := httptest.NewRecorder()

		// Set up chi router context for URL parameters
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("userId", userID)
		rctx.URLParams.Add("assetId", assetID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		// Act
		handler.Delete(rr, req)

		// Assert
		assert.Equal(t, http.StatusOK, rr.Code)
		mockService.AssertCalled(t, "DeleteFavourite", userID, assetID)
	})
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}
