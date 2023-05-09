package userFavourite

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	mock_services "github.com/loukaspe/platform-go-challenge/mocks/mock_internal/core/services"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"github.com/loukaspe/platform-go-challenge/pkg/logger"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"testing"
)

func TestGetUserFavouriteHandler_GetUserFavouriteController(t *testing.T) {
	logger := logger.NewLogger(context.Background())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockUserFavouriteServiceInterface(mockCtrl)

	type args struct {
		userIdAsString string
		userId         uint
	}

	tests := []struct {
		name                     string
		args                     args
		mockServiceResponseData  []domain.Asset
		mockServiceResponseError error
		expected                 string
		expectedStatusCode       int
	}{
		{
			name: "valid",
			args: args{
				userIdAsString: "1",
				userId:         1,
			},
			mockServiceResponseData: []domain.Asset{
				domain.Chart{
					Title:       "mockTitle1",
					XAxisTitle:  "mockXTitle1",
					YAxisTitle:  "mockYTitle1",
					Data:        "mockData",
					Description: "mockDescription1",
					AssetType:   "chart",
				},
				domain.Audience{
					Gender:             "mockGender",
					BirthCountry:       "mockBirthCountry",
					AgeGroup:           "mockAgeGroup",
					HoursSpentDaily:    1,
					PurchasesLastMonth: 2,
					Description:        "mockDescription3",
					AssetType:          "audience",
				},
				domain.Insight{
					Text:        "mockText",
					Description: "mockDescription4",
					AssetType:   "insight",
				},
				domain.Chart{
					Title:       "mockTitle2",
					XAxisTitle:  "mockXTitle2",
					YAxisTitle:  "mockYTitle2",
					Data:        "mockData2",
					Description: "mockDescription2",
					AssetType:   "chart",
				},
			},
			mockServiceResponseError: nil,
			expected: `{"favouritesAssets":[{"Title":"mockTitle1","XAxisTitle":"mockXTitle1","YAxisTitle":"mockYTitle1","Data":"mockData","Description":"mockDescription1","AssetType":"chart"},{"Gender":"mockGender","BirthCountry":"mockBirthCountry","AgeGroup":"mockAgeGroup","HoursSpentDaily":1,"PurchasesLastMonth":2,"Description":"mockDescription3","AssetType":"audience"},{"Text":"mockText","Description":"mockDescription4","AssetType":"insight"},{"Title":"mockTitle2","XAxisTitle":"mockXTitle2","YAxisTitle":"mockYTitle2","Data":"mockData2","Description":"mockDescription2","AssetType":"chart"}]}
`,
			expectedStatusCode: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRequest := httptest.NewRequest(
				"GET",
				"/user/"+tt.args.userIdAsString+"/favourites",
				nil,
			)
			vars := map[string]string{
				"user_id": tt.args.userIdAsString,
			}
			mockRequest = mux.SetURLVars(mockRequest, vars)
			mockResponseRecorder := httptest.NewRecorder()

			mockService.EXPECT().GetAssets(
				gomock.Any(),
				tt.args.userId,
			).Return(tt.mockServiceResponseData, tt.mockServiceResponseError)

			handler := &GetUserFavouriteHandler{
				UserFavouriteService: mockService,
				logger:               logger,
			}
			sut := handler.GetUserFavouriteController

			sut(mockResponseRecorder, mockRequest)

			mockResponse := mockResponseRecorder.Result()
			actual, err := io.ReadAll(mockResponse.Body)
			if err != nil {
				t.Errorf("error with response reading: %v", err)
				return
			}
			actualStatusCode := mockResponse.StatusCode

			assert.Equal(t, tt.expected, string(actual))
			assert.Equal(t, tt.expectedStatusCode, actualStatusCode)
		})
	}
}

func TestGetUserFavouriteHandler_GetUserFavouriteControllerHasBadRequestError(t *testing.T) {
	logger := logger.NewLogger(context.Background())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockUserFavouriteServiceInterface(mockCtrl)

	type args struct {
		userIdAsString string
		userId         uint
	}

	tests := []struct {
		name               string
		args               args
		expected           string
		expectedStatusCode int
	}{
		{
			name: "empty user id",
			args: args{
				userIdAsString: "",
			},
			expected: `{"errorMessage":"missing user id"}
`,
			expectedStatusCode: 400,
		},
		{
			name: "user id not a number",
			args: args{
				userIdAsString: "fjskdff",
			},
			expected: `{"errorMessage":"malformed user id"}
`,
			expectedStatusCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRequest := httptest.NewRequest(
				"GET",
				"/user/"+tt.args.userIdAsString+"/favourites",
				nil,
			)
			vars := map[string]string{
				"user_id": tt.args.userIdAsString,
			}
			mockRequest = mux.SetURLVars(mockRequest, vars)
			mockResponseRecorder := httptest.NewRecorder()

			handler := &GetUserFavouriteHandler{
				UserFavouriteService: mockService,
				logger:               logger,
			}
			sut := handler.GetUserFavouriteController

			sut(mockResponseRecorder, mockRequest)

			mockResponse := mockResponseRecorder.Result()
			actual, err := io.ReadAll(mockResponse.Body)
			if err != nil {
				t.Errorf("error with response reading: %v", err)
				return
			}
			actualStatusCode := mockResponse.StatusCode

			assert.Equal(t, tt.expected, string(actual))
			assert.Equal(t, tt.expectedStatusCode, actualStatusCode)
		})
	}
}

func TestGetUserFavouriteHandler_GetUserFavouriteControllerHasServiceError(t *testing.T) {
	logger := logger.NewLogger(context.Background())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockUserFavouriteServiceInterface(mockCtrl)

	type args struct {
		userIdAsString string
		userId         uint
	}

	tests := []struct {
		name                     string
		args                     args
		mockServiceResponseError error
		expected                 []byte
		expectedStatusCode       int
	}{
		{
			name: "random error",
			args: args{
				userIdAsString: "1",
				userId:         1,
			},
			mockServiceResponseError: errors.New("random error"),
			expected: json.RawMessage(`{"errorMessage":"error in getting user favourite assets"}
`),
			expectedStatusCode: 500,
		},
		{
			name: "service user not found error",
			args: args{
				userId:         667,
				userIdAsString: "667",
			},
			mockServiceResponseError: apierrors.UserNotFoundErrorWrapper{
				ReturnedStatusCode: 404,
				OriginalError:      errors.New("user id 667 not found"),
			},
			expected: json.RawMessage(`{}
`),
			expectedStatusCode: 404,
		},
		{
			name: "service user has no favourites",
			args: args{
				userId:         667,
				userIdAsString: "667",
			},
			mockServiceResponseError: apierrors.NoFavouriteAssetsErrorWrapper{
				ReturnedStatusCode: 204,
			},
			expected: json.RawMessage(`{}
`),
			expectedStatusCode: 204,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRequest := httptest.NewRequest(
				"GET",
				"/user/"+tt.args.userIdAsString+"/favourites",
				nil,
			)
			vars := map[string]string{
				"user_id": tt.args.userIdAsString,
			}
			mockRequest = mux.SetURLVars(mockRequest, vars)
			mockResponseRecorder := httptest.NewRecorder()

			mockService.EXPECT().GetAssets(
				gomock.Any(),
				tt.args.userId,
			).Return(nil, tt.mockServiceResponseError)

			handler := &GetUserFavouriteHandler{
				UserFavouriteService: mockService,
				logger:               logger,
			}
			sut := handler.GetUserFavouriteController

			sut(mockResponseRecorder, mockRequest)

			mockResponse := mockResponseRecorder.Result()
			actual, err := io.ReadAll(mockResponse.Body)
			if err != nil {
				t.Errorf("error with response reading: %v", err)
				return
			}
			actualStatusCode := mockResponse.StatusCode

			assert.Equal(t, string(tt.expected), string(actual))
			assert.Equal(t, tt.expectedStatusCode, actualStatusCode)
		})
	}
}
