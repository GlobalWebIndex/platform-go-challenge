// internal/adapters/http/handlers/asset_handler_test.go
package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/http/middleware"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/mapping"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/ports"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAssetService is a mock implementation of ports.AssetService for testing
type MockAssetService struct {
	mock.Mock
}

func (m *MockAssetService) CreateAsset(asset domain.Asset) (domain.Asset, error) {
	args := m.Called(asset)
	return args.Get(0).(domain.Asset), args.Error(1)
}

func (m *MockAssetService) DeleteAsset(assetID string) error {
	args := m.Called(assetID)
	return args.Error(0)
}

func TestAssetHandler_Create(t *testing.T) {
	// Save the original Body getter and restore after test
	originalBodyGetter := middleware.Body
	defer func() {
		middleware.Body = originalBodyGetter
	}()

	tests := []struct {
		name                string
		method              string
		requestBody         interface{}
		setupMock           func(*MockAssetService)
		expectedStatus      int
		expectedBody        string
		validateBodySucceed bool
	}{
		{
			name:   "Happy Path - Successfully creates audience asset",
			method: http.MethodPost,
			requestBody: dto.AssetRequest{
				ID:           "asset-123",
				Type:         "audience",
				Title:        "Test Audience",
				Description:  "Test Audience Description",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Gender:       stringPtr("female"),
				AgeGroup:     stringPtr("18-24"),
				BirthCountry: stringPtr("US"),
			},
			setupMock: func(m *MockAssetService) {
				asset, _ := mapping.AssetReqToDomain(dto.AssetRequest{
					ID:           "asset-123",
					Type:         "audience",
					Title:        "Test Audience",
					Description:  "Test Audience Description",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					Gender:       stringPtr("female"),
					AgeGroup:     stringPtr("18-24"),
					BirthCountry: stringPtr("US"),
				})
				m.On("CreateAsset", mock.AnythingOfType("*domain.Audience")).Return(asset, nil)
			},
			expectedStatus:      http.StatusCreated,
			expectedBody:        "", // We'll validate JSON structure separately
			validateBodySucceed: true,
		},
		{
			name:   "Happy Path - Successfully creates chart asset",
			method: http.MethodPost,
			requestBody: dto.AssetRequest{
				ID:          "chart-123",
				Type:        "chart",
				Title:       "Test Chart",
				Description: "Test Chart Description",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				AxesTitles:  []string{"X Axis", "Y Axis"},
				Data:        [][]float64{{1}, {2}, {3}},
			},
			setupMock: func(m *MockAssetService) {
				asset, _ := mapping.AssetReqToDomain(dto.AssetRequest{
					ID:          "chart-123",
					Type:        "chart",
					Title:       "Test Chart",
					Description: "Test Chart Description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					AxesTitles:  []string{"X Axis", "Y Axis"},
					Data:        [][]float64{{1}, {2}, {3}},
				})
				m.On("CreateAsset", mock.AnythingOfType("*domain.Chart")).Return(asset, nil)
			},
			expectedStatus:      http.StatusCreated,
			expectedBody:        "",
			validateBodySucceed: true,
		},
		{
			name:           "Unhappy Path - Wrong HTTP method",
			method:         http.MethodGet,
			requestBody:    nil,
			setupMock:      func(m *MockAssetService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
		{
			name:           "Unhappy Path - Missing validated body",
			method:         http.MethodPost,
			requestBody:    nil,
			setupMock:      func(m *MockAssetService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing validated body\n",
		},
		{
			name:   "Unhappy Path - Invalid asset type",
			method: http.MethodPost,
			requestBody: dto.AssetRequest{
				ID:    "asset-123",
				Type:  "invalid-type",
				Title: "Test Asset",
			},
			setupMock:           func(m *MockAssetService) {},
			expectedStatus:      http.StatusBadRequest,
			expectedBody:        "invalid asset type: invalid-type\n",
			validateBodySucceed: true,
		},
		{
			name:   "Unhappy Path - Service returns error",
			method: http.MethodPost,
			requestBody: dto.AssetRequest{
				ID:    "asset-123",
				Type:  "audience",
				Title: "Test Asset",
			},
			setupMock: func(m *MockAssetService) {
				asset, _ := mapping.AssetReqToDomain(dto.AssetRequest{
					ID:    "asset-123",
					Type:  "audience",
					Title: "Test Asset",
				})
				m.On("CreateAsset", mock.AnythingOfType("*domain.Audience")).Return(asset, errors.New("service error"))
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedBody:        "service error\n",
			validateBodySucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockAssetService)
			tt.setupMock(mockService)
			handler := NewAssetHandler(mockService)

			var body []byte
			if tt.requestBody != nil {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(tt.method, "/assets", bytes.NewReader(body))
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
			} else if tt.expectedStatus == http.StatusCreated {
				// Validate JSON response structure for success cases
				var response dto.AssetCreationResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotEmpty(t, response.ID)
				assert.NotEmpty(t, response.Type)
				assert.NotEmpty(t, response.Title)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestAssetHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		assetID        string
		setupMock      func(*MockAssetService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:    "Happy Path - Successfully deletes asset",
			method:  http.MethodDelete,
			assetID: "asset-123",
			setupMock: func(m *MockAssetService) {
				m.On("DeleteAsset", "asset-123").Return(nil)
			},
			expectedStatus: http.StatusNoContent,
			expectedBody:   "",
		},
		{
			name:           "Unhappy Path - Wrong HTTP method",
			method:         http.MethodPost,
			assetID:        "asset-123",
			setupMock:      func(m *MockAssetService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
		{
			name:           "Unhappy Path - Missing asset ID",
			method:         http.MethodDelete,
			assetID:        "",
			setupMock:      func(m *MockAssetService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing asset id\n",
		},
		{
			name:    "Unhappy Path - Service returns error",
			method:  http.MethodDelete,
			assetID: "asset-999",
			setupMock: func(m *MockAssetService) {
				m.On("DeleteAsset", "asset-999").Return(errors.New("delete failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "delete failed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockAssetService)
			tt.setupMock(mockService)
			handler := NewAssetHandler(mockService)

			req := httptest.NewRequest(tt.method, "/assets/"+tt.assetID, nil)
			rr := httptest.NewRecorder()

			// Set up chi router context for URL parameters
			if tt.assetID != "" {
				rctx := chi.NewRouteContext()
				rctx.URLParams.Add("assetId", tt.assetID)
				req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
			}

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

func TestNewAssetHandler(t *testing.T) {
	t.Run("should create new asset handler", func(t *testing.T) {
		// Arrange
		mockService := new(MockAssetService)

		// Act
		handler := NewAssetHandler(mockService)

		// Assert
		assert.NotNil(t, handler)
		assert.Equal(t, mockService, handler.service)
	})
}

func TestAssetHandler_Create_AssetTypes(t *testing.T) {
	// Save the original Body getter and restore after test
	originalBodyGetter := middleware.Body
	defer func() {
		middleware.Body = originalBodyGetter
	}()

	assetTypes := []struct {
		name        string
		assetType   string
		requestBody dto.AssetRequest
	}{
		{
			name:      "Audience Asset",
			assetType: "audience",
			requestBody: dto.AssetRequest{
				ID:          "aud-123",
				Type:        "audience",
				Title:       "Test Audience",
				Description: "Audience Description",
				Gender:      stringPointer("female"),
				AgeGroup:    stringPointer("18-24"),
			},
		},
		{
			name:      "Chart Asset",
			assetType: "chart",
			requestBody: dto.AssetRequest{
				ID:          "chart-123",
				Type:        "chart",
				Title:       "Test Chart",
				Description: "Chart Description",
				AxesTitles:  []string{"X", "Y"},
			},
		},
		{
			name:      "Insight Asset",
			assetType: "insight",
			requestBody: dto.AssetRequest{
				ID:          "insight-123",
				Type:        "insight",
				Title:       "Test Insight",
				Description: "Insight Description",
			},
		},
	}

	for _, tt := range assetTypes {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockAssetService)
			handler := NewAssetHandler(mockService)

			// Create expected domain asset
			expectedAsset, err := mapping.AssetReqToDomain(tt.requestBody)
			assert.NoError(t, err)

			mockService.On("CreateAsset", mock.Anything).Return(expectedAsset, nil)
			middleware.Body = MockBodyGetter{
				MockedBody:    tt.requestBody,
				ShouldSucceed: true,
			}

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/assets", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			// Act
			handler.Create(rr, req)

			// Assert
			assert.Equal(t, http.StatusCreated, rr.Code)

			var response dto.AssetCreationResponse
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.requestBody.ID, response.ID)
			assert.Equal(t, tt.requestBody.Type, response.Type)
			assert.Equal(t, tt.requestBody.Title, response.Title)

			mockService.AssertExpectations(t)
		})
	}
}

// Test interface implementation
func TestAssetHandler_InterfaceImplementation(t *testing.T) {
	// Arrange & Act
	var handler ports.AssetHandler = NewAssetHandler(nil)

	// Assert
	assert.NotNil(t, handler)
}

// Helper functions
func stringPointer(s string) *string {
	return &s
}
