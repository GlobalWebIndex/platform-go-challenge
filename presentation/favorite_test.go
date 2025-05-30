package presentation

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"pgc/app/mocks"
	"pgc/domain"
)

func TestFavoriteHandler_GetUserFavorites(t *testing.T) {
	tests := []struct {
		name               string
		userIDHeader       string
		favoriteService    func(mock *gomock.Controller) domain.FavoriteService
		expectedStatusCode int
		expectedResponse   any
	}{
		{
			name:         "Valid Request - With Favorites",
			userIDHeader: "1",
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					GetUserFavorites(gomock.Any(), 1).
					Return(
						[]domain.Favorite{
							{ID: 1, UserID: 1, AssetID: 1, Description: "Test favorite"},
							{ID: 2, UserID: 1, AssetID: 2, Description: "Another favorite"},
						},
						nil,
					)
				return service
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse: []domain.Favorite{
				{ID: 1, UserID: 1, AssetID: 1, Description: "Test favorite"},
				{ID: 2, UserID: 1, AssetID: 2, Description: "Another favorite"},
			},
		},
		{
			name:         "Valid Request - No Favorites Found",
			userIDHeader: "1",
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					GetUserFavorites(gomock.Any(), 1).
					Return([]domain.Favorite{}, nil)
				return service
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   nil,
		},
		{
			name:               "Missing User ID Header",
			userIDHeader:       "",
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "X-User-ID header is required"},
		},
		{
			name:               "Invalid User ID Format",
			userIDHeader:       "invalid",
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "invalid user ID format"},
		},
		{
			name:         "Invalid User ID Error from Service",
			userIDHeader: "0",
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					GetUserFavorites(gomock.Any(), 0).
					Return(nil, domain.ErrInvalidUserID)
				return service
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: domain.ErrInvalidUserID.Error()},
		},
		{
			name:         "Internal Server Error",
			userIDHeader: "1",
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					GetUserFavorites(gomock.Any(), 1).
					Return(nil, domain.ErrInternalServerError)
				return service
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   ErrorResponse{Error: domain.ErrInternalServerError.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := NewFavoriteHandler(tt.favoriteService(ctrl))

			req := httptest.NewRequest(http.MethodGet, "/favorite/", nil)
			if tt.userIDHeader != "" {
				req.Header.Set("X-User-ID", tt.userIDHeader)
			}

			w := httptest.NewRecorder()
			handler.GetUserFavorites(w, req)

			require.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedResponse != nil {
				var actualResponse interface{}
				if tt.expectedStatusCode == http.StatusOK {
					var favorites []domain.Favorite
					err := json.Unmarshal(w.Body.Bytes(), &favorites)
					require.NoError(t, err)
					actualResponse = favorites
				} else {
					var errorResp ErrorResponse
					err := json.Unmarshal(w.Body.Bytes(), &errorResp)
					require.NoError(t, err)
					actualResponse = errorResp
				}
				require.Equal(t, tt.expectedResponse, actualResponse)
			}
		})
	}
}

func TestFavoriteHandler_AddAssetToUserFavorites(t *testing.T) {
	tests := []struct {
		name               string
		userIDHeader       string
		assetIDParam       string
		requestBody        string
		favoriteService    func(mock *gomock.Controller) domain.FavoriteService
		expectedStatusCode int
		expectedResponse   any
	}{
		{
			name:         "Valid Request",
			userIDHeader: "1",
			assetIDParam: "1",
			requestBody:  `{"description":"Test description"}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Add(gomock.Any(), domain.Favorite{UserID: 1, AssetID: 1, Description: "Test description"}).
					Return(
						domain.Favorite{ID: 1, UserID: 1, AssetID: 1, Description: "Test description"},
						nil,
					)
				return service
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   domain.Favorite{ID: 1, UserID: 1, AssetID: 1, Description: "Test description"},
		},
		{
			name:         "Valid Request - Empty Description",
			userIDHeader: "1",
			assetIDParam: "1",
			requestBody:  `{"description":""}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Add(gomock.Any(), domain.Favorite{UserID: 1, AssetID: 1, Description: ""}).
					Return(
						domain.Favorite{ID: 1, UserID: 1, AssetID: 1, Description: ""},
						nil,
					)
				return service
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   domain.Favorite{ID: 1, UserID: 1, AssetID: 1, Description: ""},
		},
		{
			name:               "Missing User ID Header",
			userIDHeader:       "",
			assetIDParam:       "1",
			requestBody:        `{"description":"Test description"}`,
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "X-User-ID header is required"},
		},
		{
			name:               "Invalid User ID Format",
			userIDHeader:       "invalid",
			assetIDParam:       "1",
			requestBody:        `{"description":"Test description"}`,
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "invalid user ID format"},
		},
		{
			name:               "Invalid Asset ID Format",
			userIDHeader:       "1",
			assetIDParam:       "invalid",
			requestBody:        `{"description":"Test description"}`,
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "invalid asset ID format"},
		},
		{
			name:               "Invalid JSON Body",
			userIDHeader:       "1",
			assetIDParam:       "1",
			requestBody:        `{"description":}`,
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "invalid json body"},
		},
		{
			name:         "Duplicate Favorite Error",
			userIDHeader: "1",
			assetIDParam: "1",
			requestBody:  `{"description":"Test description"}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Add(gomock.Any(), domain.Favorite{UserID: 1, AssetID: 1, Description: "Test description"}).
					Return(domain.Favorite{}, domain.ErrDuplicateFavorite)
				return service
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   ErrorResponse{Error: domain.ErrDuplicateFavorite.Error()},
		},
		{
			name:         "Invalid User ID Error from Service",
			userIDHeader: "1",
			assetIDParam: "1",
			requestBody:  `{"description":"Test description"}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Add(gomock.Any(), domain.Favorite{UserID: 1, AssetID: 1, Description: "Test description"}).
					Return(domain.Favorite{}, domain.ErrInvalidUserID)
				return service
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: domain.ErrInvalidUserID.Error()},
		},
		{
			name:         "Invalid Asset ID Error from Service",
			userIDHeader: "1",
			assetIDParam: "1",
			requestBody:  `{"description":"Test description"}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Add(gomock.Any(), domain.Favorite{UserID: 1, AssetID: 1, Description: "Test description"}).
					Return(domain.Favorite{}, domain.ErrInvalidAssetID)
				return service
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: domain.ErrInvalidAssetID.Error()},
		},
		{
			name:         "Internal Server Error",
			userIDHeader: "1",
			assetIDParam: "1",
			requestBody:  `{"description":"Test description"}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Add(gomock.Any(), domain.Favorite{UserID: 1, AssetID: 1, Description: "Test description"}).
					Return(domain.Favorite{}, domain.ErrInternalServerError)
				return service
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   ErrorResponse{Error: domain.ErrInternalServerError.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := NewFavoriteHandler(tt.favoriteService(ctrl))

			var body *strings.Reader
			if tt.requestBody != "" {
				body = strings.NewReader(tt.requestBody)
			} else {
				body = strings.NewReader("")
			}

			req := httptest.NewRequest(http.MethodPost, "/favorite/"+tt.assetIDParam, body)
			req.SetPathValue(pathParamAssetID, tt.assetIDParam)
			if tt.userIDHeader != "" {
				req.Header.Set("X-User-ID", tt.userIDHeader)
			}

			w := httptest.NewRecorder()
			handler.AddAssetToUserFavorites(w, req)

			require.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedResponse != nil {
				var actualResponse interface{}
				if tt.expectedStatusCode == http.StatusCreated {
					var favorite domain.Favorite
					err := json.Unmarshal(w.Body.Bytes(), &favorite)
					require.NoError(t, err)
					actualResponse = favorite
				} else {
					var errorResp ErrorResponse
					err := json.Unmarshal(w.Body.Bytes(), &errorResp)
					require.NoError(t, err)
					actualResponse = errorResp
				}
				require.Equal(t, tt.expectedResponse, actualResponse)
			}
		})
	}
}

func TestFavoriteHandler_RemoveAssetFromUserFavorites(t *testing.T) {
	tests := []struct {
		name               string
		userIDHeader       string
		favoriteIDParam    string
		favoriteService    func(mock *gomock.Controller) domain.FavoriteService
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:            "Valid Request",
			userIDHeader:    "1",
			favoriteIDParam: "1",
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Delete(gomock.Any(), 1, 1).
					Return(nil)
				return service
			},
			expectedStatusCode: http.StatusNoContent,
			expectedResponse:   nil,
		},
		{
			name:               "Missing User ID Header",
			userIDHeader:       "",
			favoriteIDParam:    "1",
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "X-User-ID header is required"},
		},
		{
			name:               "Invalid User ID Format",
			userIDHeader:       "invalid",
			favoriteIDParam:    "1",
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "invalid user ID format"},
		},
		{
			name:               "Missing Favorite ID Parameter",
			userIDHeader:       "1",
			favoriteIDParam:    "",
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "asset ID is required in URL path"},
		},
		{
			name:               "Invalid Favorite ID Format",
			userIDHeader:       "1",
			favoriteIDParam:    "invalid",
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "invalid asset ID format"},
		},
		{
			name:            "Favorite Not Found Error",
			userIDHeader:    "1",
			favoriteIDParam: "1",
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Delete(gomock.Any(), 1, 1).
					Return(domain.ErrFavoriteNotFound)
				return service
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   nil,
		},
		{
			name:            "Invalid User ID Error from Service",
			userIDHeader:    "1",
			favoriteIDParam: "1",
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Delete(gomock.Any(), 1, 1).
					Return(domain.ErrInvalidUserID)
				return service
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: domain.ErrInvalidUserID.Error()},
		},
		{
			name:            "Internal Server Error",
			userIDHeader:    "1",
			favoriteIDParam: "1",
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Delete(gomock.Any(), 1, 1).
					Return(domain.ErrInternalServerError)
				return service
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   ErrorResponse{Error: domain.ErrInternalServerError.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := NewFavoriteHandler(tt.favoriteService(ctrl))

			req := httptest.NewRequest(http.MethodDelete, "/favorite/"+tt.favoriteIDParam, nil)
			req.SetPathValue(pathParamID, tt.favoriteIDParam)
			if tt.userIDHeader != "" {
				req.Header.Set("X-User-ID", tt.userIDHeader)
			}

			w := httptest.NewRecorder()
			handler.RemoveAssetFromUserFavorites(w, req)

			require.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedResponse != nil {
				var errorResp ErrorResponse
				err := json.Unmarshal(w.Body.Bytes(), &errorResp)
				require.NoError(t, err)
				require.Equal(t, tt.expectedResponse, errorResp)
			}
		})
	}
}

func TestFavoriteHandler_UpdateUserFavoriteAsset(t *testing.T) {
	tests := []struct {
		name               string
		userIDHeader       string
		favoriteIDParam    string
		requestBody        string
		favoriteService    func(mock *gomock.Controller) domain.FavoriteService
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:            "Valid Request",
			userIDHeader:    "1",
			favoriteIDParam: "1",
			requestBody:     `{"description":"Updated description"}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Edit(gomock.Any(), 1, 1, "Updated description").
					Return(
						domain.Favorite{ID: 1, UserID: 1, AssetID: 1, Description: "Updated description"},
						nil,
					)
				return service
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   domain.Favorite{ID: 1, UserID: 1, AssetID: 1, Description: "Updated description"},
		},
		{
			name:            "Valid Request - Empty Description",
			userIDHeader:    "1",
			favoriteIDParam: "1",
			requestBody:     `{"description":""}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Edit(gomock.Any(), 1, 1, "").
					Return(
						domain.Favorite{ID: 1, UserID: 1, AssetID: 1, Description: ""},
						nil,
					)
				return service
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   domain.Favorite{ID: 1, UserID: 1, AssetID: 1, Description: ""},
		},
		{
			name:               "Missing User ID Header",
			userIDHeader:       "",
			favoriteIDParam:    "1",
			requestBody:        `{"description":"Updated description"}`,
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "X-User-ID header is required"},
		},
		{
			name:               "Invalid User ID Format",
			userIDHeader:       "invalid",
			favoriteIDParam:    "1",
			requestBody:        `{"description":"Updated description"}`,
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "invalid user ID format"},
		},
		{
			name:               "Missing Favorite ID Parameter",
			userIDHeader:       "1",
			favoriteIDParam:    "",
			requestBody:        `{"description":"Updated description"}`,
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "asset ID is required in URL path"},
		},
		{
			name:               "Invalid Favorite ID Format",
			userIDHeader:       "1",
			favoriteIDParam:    "invalid",
			requestBody:        `{"description":"Updated description"}`,
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "invalid asset ID format"},
		},
		{
			name:               "Invalid JSON Body",
			userIDHeader:       "1",
			favoriteIDParam:    "1",
			requestBody:        `{"description":}`,
			favoriteService:    func(mock *gomock.Controller) domain.FavoriteService { return mocks.NewMockFavoriteService(mock) },
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: "invalid json body"},
		},
		{
			name:            "Favorite Not Found Error",
			userIDHeader:    "1",
			favoriteIDParam: "1",
			requestBody:     `{"description":"Updated description"}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Edit(gomock.Any(), 1, 1, "Updated description").
					Return(domain.Favorite{}, domain.ErrFavoriteNotFound)
				return service
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   nil,
		},
		{
			name:            "Invalid User ID Error from Service",
			userIDHeader:    "1",
			favoriteIDParam: "1",
			requestBody:     `{"description":"Updated description"}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Edit(gomock.Any(), 1, 1, "Updated description").
					Return(domain.Favorite{}, domain.ErrInvalidUserID)
				return service
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: domain.ErrInvalidUserID.Error()},
		},
		{
			name:            "Invalid Favorite ID Error from Service",
			userIDHeader:    "1",
			favoriteIDParam: "1",
			requestBody:     `{"description":"Updated description"}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Edit(gomock.Any(), 1, 1, "Updated description").
					Return(domain.Favorite{}, domain.ErrInvalidFavoriteID)
				return service
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   ErrorResponse{Error: domain.ErrInvalidFavoriteID.Error()},
		},
		{
			name:            "Internal Server Error",
			userIDHeader:    "1",
			favoriteIDParam: "1",
			requestBody:     `{"description":"Updated description"}`,
			favoriteService: func(mock *gomock.Controller) domain.FavoriteService {
				service := mocks.NewMockFavoriteService(mock)
				service.EXPECT().
					Edit(gomock.Any(), 1, 1, "Updated description").
					Return(domain.Favorite{}, domain.ErrInternalServerError)
				return service
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   ErrorResponse{Error: domain.ErrInternalServerError.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			handler := NewFavoriteHandler(tt.favoriteService(ctrl))

			var body *strings.Reader
			if tt.requestBody != "" {
				body = strings.NewReader(tt.requestBody)
			} else {
				body = strings.NewReader("")
			}

			req := httptest.NewRequest(http.MethodPatch, "/favorite/"+tt.favoriteIDParam, body)
			req.SetPathValue(pathParamID, tt.favoriteIDParam)
			if tt.userIDHeader != "" {
				req.Header.Set("X-User-ID", tt.userIDHeader)
			}

			w := httptest.NewRecorder()
			handler.UpdateUserFavoriteAsset(w, req)

			require.Equal(t, tt.expectedStatusCode, w.Code)

			if tt.expectedResponse != nil {
				var actualResponse interface{}
				if tt.expectedStatusCode == http.StatusOK {
					var favorite domain.Favorite
					err := json.Unmarshal(w.Body.Bytes(), &favorite)
					require.NoError(t, err)
					actualResponse = favorite
				} else {
					var errorResp ErrorResponse
					err := json.Unmarshal(w.Body.Bytes(), &errorResp)
					require.NoError(t, err)
					actualResponse = errorResp
				}
				require.Equal(t, tt.expectedResponse, actualResponse)
			}
		})
	}
}

func TestNewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockFavoriteService(ctrl)
	handler := NewFavoriteHandler(mockService)

	require.NotNil(t, handler)
	require.Equal(t, mockService, handler.favoriteService)
}
