package services

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	mock_ports "github.com/loukaspe/platform-go-challenge/mocks/mock_internal/core/ports"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"github.com/loukaspe/platform-go-challenge/pkg/logger"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserFavouriteService_AddAsset(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockUserRepositoryInterface(mockCtrl)
	logger := logger.NewLogger(context.Background())

	type args struct {
		userId    uint
		assetId   uint
		assetType string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid",
			args: args{
				userId:    1,
				assetId:   2,
				assetType: "chart",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.EXPECT().AddFavouriteAsset(
				context.Background(),
				uint32(tt.args.userId),
				uint32(tt.args.assetId),
				tt.args.assetType,
			).Return(nil)

			u := UserFavouriteService{
				logger:     logger,
				repository: mockRepository,
			}
			err := u.AddAsset(context.Background(), tt.args.userId, tt.args.assetId, tt.args.assetType)

			assert.NoError(t, err)
		})
	}
}

func TestUserFavouriteService_AddAssetRepositoryError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockUserRepositoryInterface(mockCtrl)
	logger := logger.NewLogger(context.Background())

	type args struct {
		userId    uint
		assetId   uint
		assetType string
	}
	tests := []struct {
		name                  string
		args                  args
		mockRepoErrorReturned error
		expectedError         error
	}{
		{
			name: "repo error",
			args: args{
				userId:    1,
				assetId:   2,
				assetType: "chart",
			},
			mockRepoErrorReturned: errors.New("random error"),
			expectedError:         errors.New("random error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.EXPECT().AddFavouriteAsset(
				context.Background(),
				uint32(tt.args.userId),
				uint32(tt.args.assetId),
				tt.args.assetType,
			).Return(tt.mockRepoErrorReturned)

			u := UserFavouriteService{
				logger:     logger,
				repository: mockRepository,
			}
			actualError := u.AddAsset(context.Background(), tt.args.userId, tt.args.assetId, tt.args.assetType)

			assert.Equal(t, tt.expectedError, actualError)
		})
	}
}

func TestUserFavouriteService_RemoveAsset(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockUserRepositoryInterface(mockCtrl)
	logger := logger.NewLogger(context.Background())

	type args struct {
		userId    uint
		assetId   uint
		assetType string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid",
			args: args{
				userId:    1,
				assetId:   2,
				assetType: "chart",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.EXPECT().RemoveFavouriteAsset(
				context.Background(),
				uint32(tt.args.userId),
				uint32(tt.args.assetId),
				tt.args.assetType,
			).Return(nil)

			u := UserFavouriteService{
				logger:     logger,
				repository: mockRepository,
			}
			err := u.RemoveAsset(context.Background(), tt.args.userId, tt.args.assetId, tt.args.assetType)

			assert.NoError(t, err)
		})
	}
}

func TestUserFavouriteService_RemoveAssetRepositoryError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockUserRepositoryInterface(mockCtrl)
	logger := logger.NewLogger(context.Background())

	type args struct {
		userId    uint
		assetId   uint
		assetType string
	}
	tests := []struct {
		name                  string
		args                  args
		mockRepoErrorReturned error
		expectedError         error
	}{
		{
			name: "repo error",
			args: args{
				userId:    1,
				assetId:   2,
				assetType: "chart",
			},
			mockRepoErrorReturned: errors.New("random error"),
			expectedError:         errors.New("random error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.EXPECT().RemoveFavouriteAsset(
				context.Background(),
				uint32(tt.args.userId),
				uint32(tt.args.assetId),
				tt.args.assetType,
			).Return(tt.mockRepoErrorReturned)

			u := UserFavouriteService{
				logger:     logger,
				repository: mockRepository,
			}
			actualError := u.RemoveAsset(context.Background(), tt.args.userId, tt.args.assetId, tt.args.assetType)

			assert.Equal(t, tt.expectedError, actualError)
		})
	}
}

func TestUserFavouriteService_EditAssetDescription(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockUserRepositoryInterface(mockCtrl)
	logger := logger.NewLogger(context.Background())

	type args struct {
		userId      uint
		assetId     uint
		assetType   string
		description string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "valid",
			args: args{
				userId:      1,
				assetId:     2,
				assetType:   "chart",
				description: "description",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.EXPECT().EditFavouriteAssetDescription(
				context.Background(),
				uint32(tt.args.userId),
				uint32(tt.args.assetId),
				tt.args.assetType,
				tt.args.description,
			).Return(nil)

			u := UserFavouriteService{
				logger:     logger,
				repository: mockRepository,
			}
			err := u.EditAssetDescription(
				context.Background(),
				tt.args.userId,
				tt.args.assetId,
				tt.args.assetType,
				tt.args.description,
			)

			assert.NoError(t, err)
		})
	}
}

func TestUserFavouriteService_EditAssetDescriptionRepositoryError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockUserRepositoryInterface(mockCtrl)
	logger := logger.NewLogger(context.Background())

	type args struct {
		userId      uint
		assetId     uint
		assetType   string
		description string
	}
	tests := []struct {
		name                  string
		args                  args
		mockRepoErrorReturned error
		expectedError         error
	}{
		{
			name: "repo error",
			args: args{
				userId:      1,
				assetId:     2,
				assetType:   "chart",
				description: "description",
			},
			mockRepoErrorReturned: errors.New("random error"),
			expectedError:         errors.New("random error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.EXPECT().EditFavouriteAssetDescription(
				context.Background(),
				uint32(tt.args.userId),
				uint32(tt.args.assetId),
				tt.args.assetType,
				tt.args.description,
			).Return(tt.mockRepoErrorReturned)

			u := UserFavouriteService{
				logger:     logger,
				repository: mockRepository,
			}
			actualError := u.EditAssetDescription(
				context.Background(),
				tt.args.userId,
				tt.args.assetId,
				tt.args.assetType,
				tt.args.description,
			)

			assert.Equal(t, tt.expectedError, actualError)
		})
	}
}

func TestUserFavouriteService_GetAssets(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockUserRepositoryInterface(mockCtrl)
	logger := logger.NewLogger(context.Background())

	type args struct {
		userId uint
	}
	tests := []struct {
		name                   string
		args                   args
		mockDomainUserReturned *domain.User
		expected               []domain.Asset
	}{
		{
			name: "valid",
			args: args{
				userId: 1,
			},
			mockDomainUserReturned: &domain.User{
				Username: "",
				Password: "",
				FavouriteAssets: []domain.Asset{
					domain.Chart{},
					domain.Audience{},
					domain.Chart{},
					domain.Insight{},
					domain.Audience{},
					domain.Chart{},
					domain.Insight{},
					domain.Chart{},
					domain.Audience{},
					domain.Chart{},
				},
			},
			expected: []domain.Asset{
				domain.Chart{},
				domain.Audience{},
				domain.Chart{},
				domain.Insight{},
				domain.Audience{},
				domain.Chart{},
				domain.Insight{},
				domain.Chart{},
				domain.Audience{},
				domain.Chart{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.EXPECT().GetUser(
				context.Background(),
				uint32(tt.args.userId),
			).Return(tt.mockDomainUserReturned, nil)

			u := UserFavouriteService{
				logger:     logger,
				repository: mockRepository,
			}
			actual, err := u.GetAssets(
				context.Background(),
				tt.args.userId,
			)

			if err != nil {
				t.Error("unexpected error: " + err.Error())
			}

			assert.ElementsMatch(t, tt.expected, actual)
		})
	}
}

func TestUserFavouriteService_GetAssetsUserHasNoFavourites(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockUserRepositoryInterface(mockCtrl)
	logger := logger.NewLogger(context.Background())

	type args struct {
		userId uint
	}
	tests := []struct {
		name                   string
		args                   args
		mockDomainUserReturned *domain.User
		expected               []domain.Asset
		expectedError          error
	}{
		{
			name: "valid",
			args: args{
				userId: 1,
			},
			mockDomainUserReturned: &domain.User{
				Username:        "",
				Password:        "",
				FavouriteAssets: nil,
			},
			expected: nil,
			expectedError: apierrors.NoFavouriteAssetsErrorWrapper{
				ReturnedStatusCode: 204,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.EXPECT().GetUser(
				context.Background(),
				uint32(tt.args.userId),
			).Return(tt.mockDomainUserReturned, nil)

			u := UserFavouriteService{
				logger:     logger,
				repository: mockRepository,
			}
			actual, actualError := u.GetAssets(
				context.Background(),
				tt.args.userId,
			)

			assert.ElementsMatch(t, tt.expected, actual)
			assert.Equal(t, tt.expectedError, actualError)
		})
	}
}

func TestUserFavouriteService_GetAssetsRepositoryError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRepository := mock_ports.NewMockUserRepositoryInterface(mockCtrl)
	logger := logger.NewLogger(context.Background())

	type args struct {
		userId uint
	}
	tests := []struct {
		name                  string
		args                  args
		mockRepoErrorReturned error
		expectedError         error
	}{
		{
			name: "repo error",
			args: args{
				userId: 1,
			},
			mockRepoErrorReturned: errors.New("random error"),
			expectedError:         errors.New("random error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepository.EXPECT().GetUser(
				context.Background(),
				uint32(tt.args.userId),
			).Return(nil, tt.mockRepoErrorReturned)

			u := UserFavouriteService{
				logger:     logger,
				repository: mockRepository,
			}
			_, actualError := u.GetAssets(context.Background(), tt.args.userId)

			assert.Equal(t, tt.expectedError, actualError)
		})
	}
}
