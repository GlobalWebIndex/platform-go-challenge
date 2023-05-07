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

func TestDeleteUserFavouriteHandler_DeleteUserFavouriteAssetController(t *testing.T) {
	logger := logger.NewLogger(context.Background())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockUserFavouriteServiceInterface(mockCtrl)

	type args struct {
		userId, assetId uint
		assetType       string
	}

	tests := []struct {
		name                     string
		args                     args
		requestBody              []byte
		mockRequestData          *DeleteUserFavouriteAssetRequest
		mockServiceResponseError error
		expected                 []byte
		expectedStatusCode       int
	}{
		{
			name: "valid",
			args: args{
				userId:    667,
				assetId:   666,
				assetType: "chart",
			},
			requestBody: json.RawMessage(`{
	"assetId":666,
	"assetType":"chart"
}`),
			mockRequestData: &DeleteUserFavouriteAssetRequest{
				AssetId:   666,
				AssetType: "chart",
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
				"POST",
				"/users/"+strconv.Itoa(int(tt.args.userId))+"/favourites",
				requestBodyReader,
			)

			vars := map[string]string{
				"user_id": strconv.Itoa(int(tt.args.userId)),
			}
			mockRequest = mux.SetURLVars(mockRequest, vars)

			mockRequest.Header.Set("Content-Type", "application/json")
			mockResponseRecorder := httptest.NewRecorder()

			mockService.EXPECT().RemoveAsset(
				gomock.Any(),
				tt.args.userId,
				tt.args.assetId,
				tt.args.assetType,
			).Return(tt.mockServiceResponseError)

			assetTypeValidator := helpers.NewAssetTypeValidator()

			handler := &DeleteUserFavouriteHandler{
				UserFavouriteService: mockService,
				AssetTypeValidator:   assetTypeValidator,
				logger:               logger,
			}
			sut := handler.DeleteUserFavouriteController

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

func TestDeleteUserFavouriteHandler_DeleteUserFavouriteAssetControllerHasBadRequestError(t *testing.T) {
	logger := logger.NewLogger(context.Background())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockUserFavouriteServiceInterface(mockCtrl)

	type args struct {
		userId, assetId           uint
		userIdAsString, assetType string
	}

	tests := []struct {
		name               string
		args               args
		requestBody        []byte
		mockRequestData    *DeleteUserFavouriteAssetRequest
		expected           []byte
		expectedStatusCode int
	}{
		{
			name: "empty user id",
			args: args{
				assetId:        1,
				assetType:      "chart",
				userIdAsString: "",
			},
			requestBody: json.RawMessage(`{
	"assetType":"chart"
}`),
			mockRequestData: &DeleteUserFavouriteAssetRequest{
				AssetType: "chart",
			},
			expected: json.RawMessage(`{"errorMessage":"missing user id"}
`),
			expectedStatusCode: 400,
		},
		{
			name: "user id not a number",
			args: args{
				assetId:        1,
				userIdAsString: "fjskdff",
				assetType:      "",
			},
			requestBody: json.RawMessage(`{
	"assetId": fjskdff,
	"assetType":"chart"
}`),
			mockRequestData: &DeleteUserFavouriteAssetRequest{
				AssetType: "chart",
			},
			expected: json.RawMessage(`{"errorMessage":"malformed user id"}
`),
			expectedStatusCode: 400,
		},
		{
			name: "asset id not a number",
			args: args{
				userId:         1,
				userIdAsString: "1",
				assetType:      "chart",
			},
			requestBody: json.RawMessage(`{
	"assetId": fjskdff,
	"assetType":"chart"
}`),
			mockRequestData: &DeleteUserFavouriteAssetRequest{
				AssetType: "chart",
			},
			expected: json.RawMessage(`{"errorMessage":"malformed user favourite assets request"}
`),
			expectedStatusCode: 400,
		},
		{
			name: "unknown asset type",
			args: args{
				userId:         1,
				assetId:        2,
				userIdAsString: "1",
				assetType:      "jflskadf",
			},
			requestBody: json.RawMessage(`{
	"assetId": 2,
	"assetType":"jflskadf"
}`),
			mockRequestData: &DeleteUserFavouriteAssetRequest{
				AssetType: "jflskadf",
				AssetId:   3,
			},
			expected: json.RawMessage(`{"errorMessage":"unknown asset type: jflskadf"}
`),
			expectedStatusCode: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBodyReader := bytes.NewBuffer(tt.requestBody)

			mockRequest := httptest.NewRequest(
				"POST",
				"/users/"+tt.args.userIdAsString+"/favourites",
				requestBodyReader,
			)

			vars := map[string]string{
				"user_id": tt.args.userIdAsString,
			}
			mockRequest = mux.SetURLVars(mockRequest, vars)

			mockRequest.Header.Set("Content-Type", "application/json")
			mockResponseRecorder := httptest.NewRecorder()

			assetTypeValidator := helpers.NewAssetTypeValidator()

			handler := &DeleteUserFavouriteHandler{
				UserFavouriteService: mockService,
				AssetTypeValidator:   assetTypeValidator,
				logger:               logger,
			}
			sut := handler.DeleteUserFavouriteController

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

func TestDeleteUserFavouriteHandler_DeleteUserFavouriteAssetControllerHasServiceError(t *testing.T) {
	logger := logger.NewLogger(context.Background())
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockService := mock_services.NewMockUserFavouriteServiceInterface(mockCtrl)

	type args struct {
		userId, assetId uint
		assetType       string
	}

	tests := []struct {
		name                     string
		args                     args
		requestBody              []byte
		mockRequestData          *DeleteUserFavouriteAssetRequest
		mockServiceResponseError error
		expected                 []byte
		expectedStatusCode       int
	}{
		{
			name: "service random error",
			args: args{
				userId:    667,
				assetId:   666,
				assetType: "chart",
			},
			requestBody: json.RawMessage(`{
	"assetId":666,
	"assetType":"chart"
}`),
			mockRequestData: &DeleteUserFavouriteAssetRequest{
				AssetId:   666,
				AssetType: "chart",
			},
			mockServiceResponseError: errors.New("random error"),
			expected: json.RawMessage(`{"errorMessage":"error in deleting user favourite assets"}
`),
			expectedStatusCode: 500,
		},
		{
			name: "service user not found error",
			args: args{
				userId:    667,
				assetId:   666,
				assetType: "chart",
			},
			requestBody: json.RawMessage(`{
	"assetId":666,
	"assetType":"chart"
}`),
			mockRequestData: &DeleteUserFavouriteAssetRequest{
				AssetId:   666,
				AssetType: "chart",
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
				"POST",
				"/users/"+strconv.Itoa(int(tt.args.userId))+"/favourites",
				requestBodyReader,
			)

			vars := map[string]string{
				"user_id": strconv.Itoa(int(tt.args.userId)),
			}
			mockRequest = mux.SetURLVars(mockRequest, vars)

			mockRequest.Header.Set("Content-Type", "application/json")
			mockResponseRecorder := httptest.NewRecorder()

			mockService.EXPECT().RemoveAsset(
				gomock.Any(),
				tt.args.userId,
				tt.args.assetId,
				tt.args.assetType,
			).Return(tt.mockServiceResponseError)

			assetTypeValidator := helpers.NewAssetTypeValidator()

			handler := &DeleteUserFavouriteHandler{
				UserFavouriteService: mockService,
				AssetTypeValidator:   assetTypeValidator,
				logger:               logger,
			}
			sut := handler.DeleteUserFavouriteController

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
