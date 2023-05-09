package userFavourite

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	mock_services "github.com/loukaspe/platform-go-challenge/mocks/mock_internal/core/services"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"github.com/loukaspe/platform-go-challenge/pkg/helpers"
	"github.com/loukaspe/platform-go-challenge/pkg/logger"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestUpdateUserFavouriteHandler_UpdateUserFavouriteController(t *testing.T) {
	logger := logger.NewLogger(context.Background())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockUserFavouriteServiceInterface(mockCtrl)

	type args struct {
		userId, assetId        uint
		assetType, description string
	}

	tests := []struct {
		name                     string
		args                     args
		requestBody              []byte
		requestedUuid            string
		mockRequestData          *UpdateUserFavouriteAssetDescriptionRequest
		mockServiceResponseError error
		expected                 []byte
		expectedStatusCode       int
	}{
		{
			name: "valid",
			args: args{
				userId:      667,
				assetId:     666,
				assetType:   "chart",
				description: "description",
			},
			requestBody: json.RawMessage(`{
	"assetId":666,
	"assetType":"chart",
	"description":"description"
}`),
			mockRequestData: &UpdateUserFavouriteAssetDescriptionRequest{
				AssetType:   "chart",
				Description: "description",
			},
			mockServiceResponseError: nil,
			expected:                 json.RawMessage(``),
			expectedStatusCode:       200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBodyReader := bytes.NewBuffer(tt.requestBody)

			mockRequest := httptest.NewRequest(
				"PATCH",
				"/users/"+strconv.Itoa(int(tt.args.userId))+
					"/favourites/"+strconv.Itoa(int(tt.args.assetId)),
				requestBodyReader,
			)

			vars := map[string]string{
				"user_id":  strconv.Itoa(int(tt.args.userId)),
				"asset_id": strconv.Itoa(int(tt.args.assetId)),
			}
			mockRequest = mux.SetURLVars(mockRequest, vars)

			mockRequest.Header.Set("Content-Type", "application/json")
			mockResponseRecorder := httptest.NewRecorder()

			mockService.EXPECT().
				EditAssetDescription(
					gomock.Any(),
					tt.args.userId,
					tt.args.assetId,
					tt.args.assetType,
					tt.args.description,
				).
				Return(tt.mockServiceResponseError)

			assetTypeValidator := helpers.NewAssetTypeValidator()

			handler := &UpdateUserFavouriteHandler{
				UserFavouriteService: mockService,
				AssetTypeValidator:   assetTypeValidator,
				logger:               logger,
			}
			sut := handler.UpdateUserFavouriteController

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

func TestUpdateUserFavouriteHandler_UpdateUserFavouriteControllerHasBadRequestError(t *testing.T) {
	logger := logger.NewLogger(context.Background())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockUserFavouriteServiceInterface(mockCtrl)

	type args struct {
		userId, assetId                 uint
		userIdAsString, assetIdAsString string
		assetType, description          string
	}

	tests := []struct {
		name               string
		args               args
		requestBody        []byte
		requestedUuid      string
		mockRequestData    *UpdateUserFavouriteAssetDescriptionRequest
		expected           []byte
		expectedStatusCode int
	}{
		{
			name: "empty user id",
			args: args{
				assetId:         1,
				assetIdAsString: "1",
				assetType:       "chart",
				userIdAsString:  "",
				description:     "description",
			},
			requestBody: json.RawMessage(`{
	"assetType":"chart",
	"description":"description"
}`),
			mockRequestData: &UpdateUserFavouriteAssetDescriptionRequest{
				AssetType:   "chart",
				Description: "description",
			},
			expected: json.RawMessage(`{"errorMessage":"missing user id"}
`),
			expectedStatusCode: 400,
		},
		{
			name: "user id not a number",
			args: args{
				assetId:         1,
				assetIdAsString: "1",
				userIdAsString:  "fjskdff",
				assetType:       "chart",
				description:     "description",
			},
			requestBody: json.RawMessage(`{
	"description": "description",
	"assetType":"chart"
}`),
			mockRequestData: &UpdateUserFavouriteAssetDescriptionRequest{
				AssetType:   "chart",
				Description: "description",
			},
			expected: json.RawMessage(`{"errorMessage":"malformed user id"}
`),
			expectedStatusCode: 400,
		},
		{
			name: "empty asset id",
			args: args{
				userId:          1,
				userIdAsString:  "1",
				assetType:       "chart",
				assetIdAsString: "",
				description:     "description",
			},
			requestBody: json.RawMessage(`{
	"assetType":"chart",
	"description":"description"
}`),
			mockRequestData: &UpdateUserFavouriteAssetDescriptionRequest{
				AssetType:   "chart",
				Description: "description",
			},
			expected: json.RawMessage(`{"errorMessage":"missing asset id"}
`),
			expectedStatusCode: 400,
		},
		{
			name: "asset id not a number",
			args: args{
				userId:          1,
				userIdAsString:  "1",
				assetIdAsString: "fjskdff",
				assetType:       "chart",
				description:     "description",
			},
			requestBody: json.RawMessage(`{
	"description": "description",
	"assetType":"chart"
}`),
			mockRequestData: &UpdateUserFavouriteAssetDescriptionRequest{
				AssetType:   "chart",
				Description: "description",
			},
			expected: json.RawMessage(`{"errorMessage":"malformed asset id"}
`),
			expectedStatusCode: 400,
		},
		{
			name: "unknown asset type",
			args: args{
				userId:          1,
				assetId:         2,
				userIdAsString:  "1",
				assetIdAsString: "2",
				assetType:       "jflskadf",
				description:     "description",
			},
			requestBody: json.RawMessage(`{
	"description": "description",
	"assetId": 2,
	"assetType":"jflskadf"
}`),
			mockRequestData: &UpdateUserFavouriteAssetDescriptionRequest{
				AssetType:   "jflskadf",
				Description: "description",
			},
			expected: json.RawMessage(`{"errorMessage":"unknown asset type: jflskadf"}
`),
			expectedStatusCode: 400,
		},
		{
			name: "empty description",
			args: args{
				userId:          1,
				assetId:         2,
				userIdAsString:  "1",
				assetIdAsString: "2",
				assetType:       "chart",
				description:     "",
			},
			requestBody: json.RawMessage(`{
	"description": "",
	"assetId": 2,
	"assetType":"chart"
}`),
			mockRequestData: &UpdateUserFavouriteAssetDescriptionRequest{
				AssetType:   "chart",
				Description: "",
			},
			expected: json.RawMessage(`{"errorMessage":"empty description"}
`),
			expectedStatusCode: 400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBodyReader := bytes.NewBuffer(tt.requestBody)

			mockRequest := httptest.NewRequest(
				"PATCH",
				"/users/"+tt.args.userIdAsString+
					"/favourites/"+tt.args.assetIdAsString,
				requestBodyReader,
			)

			vars := map[string]string{
				"user_id":  tt.args.userIdAsString,
				"asset_id": tt.args.assetIdAsString,
			}
			mockRequest = mux.SetURLVars(mockRequest, vars)

			mockRequest.Header.Set("Content-Type", "application/json")
			mockResponseRecorder := httptest.NewRecorder()

			assetTypeValidator := helpers.NewAssetTypeValidator()

			handler := &UpdateUserFavouriteHandler{
				UserFavouriteService: mockService,
				AssetTypeValidator:   assetTypeValidator,
				logger:               logger,
			}
			sut := handler.UpdateUserFavouriteController

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

func TestUpdateUserFavouriteHandler_UpdateUserFavouriteControllerHasServiceError(t *testing.T) {
	logger := logger.NewLogger(context.Background())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockUserFavouriteServiceInterface(mockCtrl)

	type args struct {
		userId, assetId        uint
		assetType, description string
	}

	tests := []struct {
		name                     string
		args                     args
		requestBody              []byte
		requestedUuid            string
		mockRequestData          *UpdateUserFavouriteAssetDescriptionRequest
		mockServiceResponseError error
		expected                 []byte
		expectedStatusCode       int
	}{
		{
			name: "service random error",
			args: args{
				userId:      667,
				assetId:     666,
				assetType:   "chart",
				description: "description",
			},
			requestBody: json.RawMessage(`{
	"description":"description",
	"assetType":"chart"
}`),
			mockRequestData: &UpdateUserFavouriteAssetDescriptionRequest{
				Description: "description",
				AssetType:   "chart",
			},
			mockServiceResponseError: errors.New("random error"),
			expected: json.RawMessage(`{"errorMessage":"error in updating user favourite assets description"}
`),
			expectedStatusCode: 500,
		},
		{
			name: "service user not found error",
			args: args{
				userId:      667,
				assetId:     666,
				assetType:   "chart",
				description: "description",
			},
			requestBody: json.RawMessage(`{
	"description":"description",
	"assetType":"chart"
}`),
			mockRequestData: &UpdateUserFavouriteAssetDescriptionRequest{
				Description: "description",
				AssetType:   "chart",
			},
			mockServiceResponseError: apierrors.UserNotFoundErrorWrapper{
				ReturnedStatusCode: 404,
				OriginalError:      errors.New("user id 667 not found"),
			},
			expected: json.RawMessage(`{}
`),
			expectedStatusCode: 404,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBodyReader := bytes.NewBuffer(tt.requestBody)

			mockRequest := httptest.NewRequest(
				"PATCH",
				"/users/"+strconv.Itoa(int(tt.args.userId))+
					"/favourites/"+strconv.Itoa(int(tt.args.assetId)),
				requestBodyReader,
			)

			vars := map[string]string{
				"user_id":  strconv.Itoa(int(tt.args.userId)),
				"asset_id": strconv.Itoa(int(tt.args.assetId)),
			}
			mockRequest = mux.SetURLVars(mockRequest, vars)

			mockRequest.Header.Set("Content-Type", "application/json")
			mockResponseRecorder := httptest.NewRecorder()

			mockService.EXPECT().
				EditAssetDescription(
					gomock.Any(),
					tt.args.userId,
					tt.args.assetId,
					tt.args.assetType,
					tt.args.description,
				).
				Return(tt.mockServiceResponseError)

			assetTypeValidator := helpers.NewAssetTypeValidator()

			handler := &UpdateUserFavouriteHandler{
				UserFavouriteService: mockService,
				AssetTypeValidator:   assetTypeValidator,
				logger:               logger,
			}
			sut := handler.UpdateUserFavouriteController

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
