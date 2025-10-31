package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/http/middleware"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of ports.UserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetUserByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) GetAllUsers() ([]domain.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserService) UpdateUser(user domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) GetFavouritesByUser(id string) ([]domain.Favourite, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Favourite), args.Error(1)
}

type MockBodyGetter struct {
	MockedBody    any
	ShouldSucceed bool
}

func NewMockBodyGetter(mockedBody any) MockBodyGetter {
	return MockBodyGetter{
		MockedBody:    mockedBody,
		ShouldSucceed: true, // Default to true
	}
}

func (m MockBodyGetter) GetValidatedBody(r *http.Request) (any, bool) {
	if m.MockedBody == nil {
		return nil, false
	}
	return m.MockedBody, m.ShouldSucceed
}
func TestUserHandler_Create(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		requestBody    interface{}
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Happy Path - Successfully creates user",
			method: http.MethodPost,
			requestBody: dto.CreateUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "Password123",
			},
			setupMock: func(m *MockUserService) {
				m.On("CreateUser", mock.AnythingOfType("domain.User")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
		},
		{
			name:           "Unhappy Path - Wrong HTTP method",
			method:         http.MethodGet,
			requestBody:    nil,
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
		{
			name:           "Unhappy Path - Missing validated body",
			method:         http.MethodPost,
			requestBody:    nil,
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing validated body\n",
		},
		{
			name:   "Unhappy Path - Service returns error",
			method: http.MethodPost,
			requestBody: dto.CreateUserRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "Password123",
			},
			setupMock: func(m *MockUserService) {
				m.On("CreateUser", mock.AnythingOfType("domain.User")).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockUserService)
			tt.setupMock(mockService)
			handler := NewUserHandler(mockService)

			var body []byte
			if tt.requestBody != nil {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(tt.method, "/users", bytes.NewReader(body))
			rr := httptest.NewRecorder()

			// Set up the mock middleware based on whether we have a request body
			if tt.requestBody != nil {
				// Use your mock middleware to provide validated body
				middleware.Body = NewMockBodyGetter(tt.requestBody)
			} else {
				// Set to return false when no body is provided
				middleware.Body = NewMockBodyGetter(nil)
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

func TestUserHandler_Get(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		userID         string
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "Happy Path - Successfully gets user",
			method: http.MethodGet,
			userID: "user-123",
			setupMock: func(m *MockUserService) {
				user := &domain.User{
					Id:    "user-123",
					Name:  "John Doe",
					Email: "john@example.com",
				}
				m.On("GetUserByID", "user-123").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.UserResponse{
				ID:    "user-123",
				Name:  "John Doe",
				Email: "john@example.com",
			},
		},
		{
			name:           "Unhappy Path - Wrong HTTP method",
			method:         http.MethodPost,
			userID:         "user-123",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
		{
			name:           "Unhappy Path - Missing user ID",
			method:         http.MethodGet,
			userID:         "",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing user id\n",
		},
		{
			name:   "Unhappy Path - User not found",
			method: http.MethodGet,
			userID: "user-999",
			setupMock: func(m *MockUserService) {
				m.On("GetUserByID", "user-999").Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "user not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockUserService)
			tt.setupMock(mockService)
			handler := NewUserHandler(mockService)

			req := httptest.NewRequest(tt.method, "/users/"+tt.userID, nil)
			rr := httptest.NewRecorder()

			// Set up chi router context for URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Act
			handler.Get(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response dto.UserResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			} else {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_List(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:   "Happy Path - Successfully lists users",
			method: http.MethodGet,
			setupMock: func(m *MockUserService) {
				users := []domain.User{
					{Id: "user-1", Name: "John Doe", Email: "john@example.com"},
					{Id: "user-2", Name: "Jane Smith", Email: "jane@example.com"},
				}
				m.On("GetAllUsers").Return(users, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:   "Unhappy Path - Wrong HTTP method",
			method: http.MethodPost,
			setupMock: func(m *MockUserService) {
			},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCount:  0,
		},
		{
			name:   "Unhappy Path - Service returns error",
			method: http.MethodGet,
			setupMock: func(m *MockUserService) {
				m.On("GetAllUsers").Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
		{
			name:   "Happy Path - Empty user list",
			method: http.MethodGet,
			setupMock: func(m *MockUserService) {
				m.On("GetAllUsers").Return([]domain.User{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockUserService)
			tt.setupMock(mockService)
			handler := NewUserHandler(mockService)

			req := httptest.NewRequest(tt.method, "/users", nil)
			rr := httptest.NewRecorder()

			// Act
			handler.List(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response []dto.UserResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response, tt.expectedCount)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		userID         string
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Happy Path - Successfully deletes user",
			method: http.MethodDelete,
			userID: "user-123",
			setupMock: func(m *MockUserService) {
				m.On("DeleteUser", "user-123").Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "",
		},
		{
			name:           "Unhappy Path - Wrong HTTP method",
			method:         http.MethodPost,
			userID:         "user-123",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
		{
			name:           "Unhappy Path - Missing user ID",
			method:         http.MethodDelete,
			userID:         "",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing user id\n",
		},
		{
			name:   "Unhappy Path - User not found",
			method: http.MethodDelete,
			userID: "user-999",
			setupMock: func(m *MockUserService) {
				m.On("DeleteUser", "user-999").Return(errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "user not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockUserService)
			tt.setupMock(mockService)
			handler := NewUserHandler(mockService)

			req := httptest.NewRequest(tt.method, "/users/"+tt.userID, nil)
			rr := httptest.NewRecorder()

			// Set up chi router context for URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
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
func TestUserHandler_Update(t *testing.T) {
	// Save the original Body getter and restore after test
	originalBodyGetter := middleware.Body
	defer func() {
		middleware.Body = originalBodyGetter
	}()

	tests := []struct {
		name                string
		method              string
		userID              string
		requestBody         interface{}
		setupMock           func(*MockUserService)
		expectedStatus      int
		expectedBody        string
		validateBodySucceed bool
	}{
		{
			name:   "Happy Path - Successfully updates user with PUT",
			method: http.MethodPut,
			userID: "user-123",
			requestBody: dto.UpdateUserRequest{
				Name:  "John Updated",
				Email: "john.updated@example.com",
			},
			setupMock: func(m *MockUserService) {
				existingUser := &domain.User{
					Id:    "user-123",
					Name:  "John Doe",
					Email: "john@example.com",
				}
				m.On("GetUserByID", "user-123").Return(existingUser, nil)
				m.On("UpdateUser", mock.AnythingOfType("domain.User")).Return(nil)
			},
			expectedStatus:      http.StatusOK,
			expectedBody:        `{"status":"updated"}` + "\n",
			validateBodySucceed: true,
		},
		{
			name:   "Happy Path - Successfully updates user with PATCH",
			method: http.MethodPatch,
			userID: "user-123",
			requestBody: dto.UpdateUserRequest{
				Email: "john.patched@example.com",
			},
			setupMock: func(m *MockUserService) {
				existingUser := &domain.User{
					Id:    "user-123",
					Name:  "John Doe",
					Email: "john@example.com",
				}
				m.On("GetUserByID", "user-123").Return(existingUser, nil)
				m.On("UpdateUser", mock.AnythingOfType("domain.User")).Return(nil)
			},
			expectedStatus:      http.StatusOK,
			expectedBody:        `{"status":"updated"}` + "\n",
			validateBodySucceed: true,
		},
		{
			name:           "Unhappy Path - Wrong HTTP method",
			method:         http.MethodPost,
			userID:         "user-123",
			requestBody:    nil,
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   "method not allowed\n",
		},
		{
			name:           "Unhappy Path - Missing user ID",
			method:         http.MethodPut,
			userID:         "",
			requestBody:    nil,
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing user id\n",
		},
		{
			name:           "Unhappy Path - Missing validated body",
			method:         http.MethodPut,
			userID:         "user-123",
			requestBody:    nil,
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "missing validated body\n",
		},
		{
			name:   "Unhappy Path - User not found",
			method: http.MethodPut,
			userID: "user-999",
			requestBody: dto.UpdateUserRequest{
				Name: "Updated Name",
			},
			setupMock: func(m *MockUserService) {
				m.On("GetUserByID", "user-999").Return(nil, errors.New("not found"))
			},
			expectedStatus:      http.StatusNotFound,
			expectedBody:        "user not found\n",
			validateBodySucceed: true,
		},
		{
			name:   "Unhappy Path - Update service error",
			method: http.MethodPut,
			userID: "user-123",
			requestBody: dto.UpdateUserRequest{
				Name: "Updated Name",
			},
			setupMock: func(m *MockUserService) {
				existingUser := &domain.User{
					Id:    "user-123",
					Name:  "John Doe",
					Email: "john@example.com",
				}
				m.On("GetUserByID", "user-123").Return(existingUser, nil)
				m.On("UpdateUser", mock.AnythingOfType("domain.User")).Return(errors.New("update failed"))
			},
			expectedStatus:      http.StatusInternalServerError,
			expectedBody:        "update failed\n",
			validateBodySucceed: true,
		},
		{
			name:   "Unhappy Path - Invalid update request (empty fields)",
			method: http.MethodPut,
			userID: "user-123",
			requestBody: dto.UpdateUserRequest{
				Email: "invalid-email", // Invalid email format
			},
			setupMock: func(m *MockUserService) {
				// No service calls expected since validation should fail
			},
			expectedStatus:      http.StatusBadRequest,
			expectedBody:        "missing validated body\n", // This will depend on your validation middleware
			validateBodySucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockUserService)
			tt.setupMock(mockService)
			handler := NewUserHandler(mockService)

			var body []byte
			if tt.requestBody != nil {
				body, _ = json.Marshal(tt.requestBody)
			}

			req := httptest.NewRequest(tt.method, "/users/"+tt.userID, bytes.NewReader(body))
			rr := httptest.NewRecorder()

			// Set up chi router context for URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Set up the mock middleware based on whether we have a request body
			if tt.requestBody != nil {
				// Use your mock middleware to provide validated body
				middleware.Body = MockBodyGetter{
					MockedBody:    tt.requestBody,
					ShouldSucceed: tt.validateBodySucceed,
				}
			} else {
				// Set to return false when no body is provided
				middleware.Body = MockBodyGetter{
					MockedBody:    nil,
					ShouldSucceed: false, // This will make GetValidatedBody return false
				}
			}

			// Act
			handler.Update(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)
			if tt.expectedBody != "" {
				if strings.HasPrefix(tt.expectedBody, "{") {
					assert.JSONEq(t, tt.expectedBody, rr.Body.String())
				} else {
					assert.Equal(t, tt.expectedBody, rr.Body.String())
				}
			}
			mockService.AssertExpectations(t)
		})
	}
}
func TestUserHandler_GetFavourites(t *testing.T) {
	// Create sample time for consistent testing
	sampleTime := time.Date(2023, 10, 15, 14, 30, 0, 0, time.UTC)

	tests := []struct {
		name           string
		method         string
		userID         string
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:   "Happy Path - Successfully gets user favourites with mixed asset types",
			method: http.MethodGet,
			userID: "user-123",
			setupMock: func(m *MockUserService) {
				favourites := []domain.Favourite{
					{
						UserID:    "user-123",
						AssetID:   "aud-001",
						CreatedAt: sampleTime,
						AssetType: domain.AssetTypeAudience,
						Audience: &domain.Audience{
							AssetBase: domain.AssetBase{
								ID:          "aud-001",
								Type:        domain.AssetTypeAudience,
								Title:       "Premium Users Audience",
								Description: "High-value customer segment",
								CreatedAt:   sampleTime,
								UpdatedAt:   sampleTime,
							},
							Gender:          "All",
							BirthCountry:    "Various",
							AgeGroup:        "25-40",
							HoursSocial:     15,
							PurchasesLastMo: 5,
						},
					},
					{
						UserID:    "user-123",
						AssetID:   "chart-001",
						CreatedAt: sampleTime,
						AssetType: domain.AssetTypeChart,
						Chart: &domain.Chart{
							AssetBase: domain.AssetBase{
								ID:          "chart-001",
								Type:        domain.AssetTypeChart,
								Title:       "Monthly Sales Chart",
								Description: "Shows monthly sales performance across regions",
								CreatedAt:   sampleTime,
								UpdatedAt:   sampleTime,
							},
							AxesTitles: []string{"Months", "Revenue"},
							Data:       [][]float64{{1, 2, 3}, {4, 5, 6}},
						},
					},
					{
						UserID:    "user-123",
						AssetID:   "insight-001",
						CreatedAt: sampleTime,
						AssetType: domain.AssetTypeInsight,
						Insight: &domain.Insight{
							AssetBase: domain.AssetBase{
								ID:          "insight-001",
								Type:        domain.AssetTypeInsight,
								Title:       "Revenue Growth Insight",
								Description: "Analysis of Q1 revenue growth patterns",
								CreatedAt:   sampleTime,
								UpdatedAt:   sampleTime,
							},
							Text: "Revenue has grown by 15% in Q1 compared to previous quarter.",
						},
					},
				}
				m.On("GetFavouritesByUser", "user-123").Return(favourites, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  3,
		},
		{
			name:   "Happy Path - Successfully gets user favourites with only charts",
			method: http.MethodGet,
			userID: "user-456",
			setupMock: func(m *MockUserService) {
				favourites := []domain.Favourite{
					{
						UserID:    "user-456",
						AssetID:   "chart-002",
						CreatedAt: sampleTime,
						AssetType: domain.AssetTypeChart,
						Chart: &domain.Chart{
							AssetBase: domain.AssetBase{
								ID:          "chart-002",
								Type:        domain.AssetTypeChart,
								Title:       "User Engagement Chart",
								Description: "Daily active users over time",
								CreatedAt:   sampleTime,
								UpdatedAt:   sampleTime,
							},
							AxesTitles: []string{"Days", "Active Users"},
							Data:       [][]float64{{1, 2, 3}, {1000, 1500, 1200}},
						},
					},
				}
				m.On("GetFavouritesByUser", "user-456").Return(favourites, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:   "Happy Path - Successfully gets user favourites with only audiences",
			method: http.MethodGet,
			userID: "user-789",
			setupMock: func(m *MockUserService) {
				favourites := []domain.Favourite{
					{
						UserID:    "user-789",
						AssetID:   "aud-002",
						CreatedAt: sampleTime,
						AssetType: domain.AssetTypeAudience,
						Audience: &domain.Audience{
							AssetBase: domain.AssetBase{
								ID:          "aud-002",
								Type:        domain.AssetTypeAudience,
								Title:       "New Customers Audience",
								Description: "Customers acquired in the last 30 days",
								CreatedAt:   sampleTime,
								UpdatedAt:   sampleTime,
							},
							Gender:          "All",
							BirthCountry:    "US",
							AgeGroup:        "18-25",
							HoursSocial:     20,
							PurchasesLastMo: 2,
						},
					},
				}
				m.On("GetFavouritesByUser", "user-789").Return(favourites, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:   "Happy Path - Successfully gets user favourites with only insights",
			method: http.MethodGet,
			userID: "user-999",
			setupMock: func(m *MockUserService) {
				favourites := []domain.Favourite{
					{
						UserID:    "user-999",
						AssetID:   "insight-002",
						CreatedAt: sampleTime,
						AssetType: domain.AssetTypeInsight,
						Insight: &domain.Insight{
							AssetBase: domain.AssetBase{
								ID:          "insight-002",
								Type:        domain.AssetTypeInsight,
								Title:       "Market Trends Insight",
								Description: "Emerging market trends analysis",
								CreatedAt:   sampleTime,
								UpdatedAt:   sampleTime,
							},
							Text: "Consumer preferences are shifting towards sustainable products.",
						},
					},
				}
				m.On("GetFavouritesByUser", "user-999").Return(favourites, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:   "Happy Path - Favourite with nil asset pointers (edge case)",
			method: http.MethodGet,
			userID: "user-555",
			setupMock: func(m *MockUserService) {
				favourites := []domain.Favourite{
					{
						UserID:    "user-555",
						AssetID:   "chart-003",
						CreatedAt: sampleTime,
						AssetType: domain.AssetTypeChart,
						Chart:     nil, // nil pointer to test resilience
					},
				}
				m.On("GetFavouritesByUser", "user-555").Return(favourites, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "Unhappy Path - Wrong HTTP method",
			method:         http.MethodPost,
			userID:         "user-123",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCount:  0,
		},
		{
			name:           "Unhappy Path - Missing user ID",
			method:         http.MethodGet,
			userID:         "",
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name:   "Unhappy Path - Favourites not found",
			method: http.MethodGet,
			userID: "user-999",
			setupMock: func(m *MockUserService) {
				m.On("GetFavouritesByUser", "user-999").Return(nil, errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedCount:  0,
		},
		{
			name:   "Happy Path - Empty favourites list",
			method: http.MethodGet,
			userID: "user-123",
			setupMock: func(m *MockUserService) {
				m.On("GetFavouritesByUser", "user-123").Return([]domain.Favourite{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockUserService)
			tt.setupMock(mockService)
			handler := NewUserHandler(mockService)

			req := httptest.NewRequest(tt.method, "/users/"+tt.userID+"/favourites", nil)
			rr := httptest.NewRecorder()

			// Set up chi router context for URL parameters
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.userID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Act
			handler.GetFavourites(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var response []dto.FavouriteResponse
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response, tt.expectedCount)

				// Additional assertions for specific test cases
				if tt.name == "Happy Path - Successfully gets user favourites with mixed asset types" && tt.expectedCount > 0 {
					// Verify that the response contains the expected data
					// Note: The mapping function should handle converting domain assets to dto responses
					for _, fav := range response {
						assert.NotEmpty(t, fav.AssetID)
						assert.NotEmpty(t, fav.UserID)
					}
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetFavourites_HappyPath_SuccessfullyGetsUserFavouritesWithMixedAssetTypes(t *testing.T) {
	// Arrange
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	// Create test favourites with proper assets
	favourites := []domain.Favourite{
		{
			UserID:    "user-123",
			AssetID:   "aud-1",
			AssetType: domain.AssetTypeAudience,
			Audience: &domain.Audience{
				AssetBase: domain.AssetBase{
					ID:          "aud-1",
					Type:        domain.AssetTypeAudience,
					Title:       "Test Audience",
					Description: "Audience Description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				Gender:       "F",
				BirthCountry: "US",
				AgeGroup:     "25-34",
			},
		},
		{
			UserID:    "user-123",
			AssetID:   "chart-1",
			AssetType: domain.AssetTypeChart,
			Chart: &domain.Chart{
				AssetBase: domain.AssetBase{
					ID:          "chart-1",
					Type:        domain.AssetTypeChart,
					Title:       "Test Chart",
					Description: "Chart Description",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				},
				AxesTitles: []string{"X", "Y"},
				Data:       [][]float64{{1, 2}, {3, 4}},
			},
		},
	}

	mockService.On("GetFavouritesByUser", "user-123").Return(favourites, nil)

	req := httptest.NewRequest("GET", "/users/user-123/favourites", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "user-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
	rr := httptest.NewRecorder()

	// Act
	handler.GetFavourites(rr, req)

	// Debug: Print the actual response
	fmt.Printf("Response Status: %d\n", rr.Code)
	fmt.Printf("Response Body: %s\n", rr.Body.String())

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotEmpty(t, rr.Body.String(), "Response body should not be empty")

	// Parse and verify the JSON response
	var response []dto.FavouriteResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2, "Should return 2 favourites")

	// Verify first favourite (Audience)
	assert.Equal(t, "user-123", response[0].UserID)
	assert.Equal(t, "aud-1", response[0].AssetID)
	assert.NotNil(t, response[0].Asset)

	// Verify second favourite (Chart)
	assert.Equal(t, "user-123", response[1].UserID)
	assert.Equal(t, "chart-1", response[1].AssetID)
	assert.NotNil(t, response[1].Asset)

	mockService.AssertExpectations(t)
}
