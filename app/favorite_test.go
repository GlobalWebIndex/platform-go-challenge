package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"pgc/domain"
	"pgc/persistence/mocks"
)

func TestFavoriteService_Add(t *testing.T) {
	tests := []struct {
		name               string
		favoriteRepository func(mock *gomock.Controller) domain.FavoriteRepository
		inputFavorite      domain.Favorite
		expectedFavorite   domain.Favorite
		expectedError      error
	}{
		{
			name: "Valid Favorite Addition",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				repo.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(
						domain.Favorite{ID: 1, UserID: 1, AssetID: 1, Description: "Test"},
						nil,
					)
				return repo
			},
			inputFavorite: domain.Favorite{
				UserID:      1,
				AssetID:     1,
				Description: "Test",
			},
			expectedFavorite: domain.Favorite{
				ID:          1,
				UserID:      1,
				AssetID:     1,
				Description: "Test",
			},
			expectedError: nil,
		},
		{
			name: "Invalid UserID",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				return repo
			},
			inputFavorite: domain.Favorite{
				UserID:  0, // Invalid UserID
				AssetID: 1,
			},
			expectedFavorite: domain.Favorite{},
			expectedError:    domain.ErrInvalidUserID,
		},
		{
			name: "Invalid AssetID",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				return repo
			},
			inputFavorite: domain.Favorite{
				UserID:  1,
				AssetID: 0, // Invalid AssetID
			},
			expectedFavorite: domain.Favorite{},
			expectedError:    domain.ErrInvalidAssetID,
		},
		{
			name: "Duplicate Favorite",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				repo.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(
						domain.Favorite{},
						domain.ErrDuplicateFavorite,
					)
				return repo
			},
			inputFavorite: domain.Favorite{
				UserID:      1,
				AssetID:     1,
				Description: "Test",
			},
			expectedFavorite: domain.Favorite{},
			expectedError:    domain.ErrDuplicateFavorite,
		},
		{
			name: "Internal Server Error",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				repo.EXPECT().
					Insert(gomock.Any(), gomock.Any()).
					Return(
						domain.Favorite{},
						domain.ErrInternalServerError,
					)
				return repo
			},
			inputFavorite: domain.Favorite{
				UserID:      1,
				AssetID:     1,
				Description: "Test",
			},
			expectedFavorite: domain.Favorite{},
			expectedError:    domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			favoriteService := NewFavoriteService(tt.favoriteRepository(ctrl))

			favorite, err := favoriteService.Add(context.Background(), tt.inputFavorite)

			if err != nil {
				require.ErrorIs(t, err, tt.expectedError)
				require.Equal(t, tt.expectedFavorite, favorite)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedFavorite, favorite)
			}
		})
	}
}

func TestFavoriteService_Edit(t *testing.T) {
	tests := []struct {
		name               string
		favoriteRepository func(mock *gomock.Controller) domain.FavoriteRepository
		favoriteID         int
		userID             int
		description        string
		expectedFavorite   domain.Favorite
		expectedError      error
	}{
		{
			name: "Valid Favorite Edit",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				repo.EXPECT().
					UpdateDescription(gomock.Any(), 1, 1, "Updated Description").
					Return(
						domain.Favorite{ID: 1, UserID: 1, AssetID: 1, Description: "Updated Description"},
						nil,
					)
				return repo
			},
			favoriteID:       1,
			userID:           1,
			description:      "Updated Description",
			expectedFavorite: domain.Favorite{ID: 1, UserID: 1, AssetID: 1, Description: "Updated Description"},
			expectedError:    nil,
		},
		{
			name: "Invalid FavoriteID",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				return repo
			},
			favoriteID:       0, // Invalid FavoriteID
			userID:           1,
			description:      "Updated Description",
			expectedFavorite: domain.Favorite{},
			expectedError:    domain.ErrInvalidFavoriteID,
		},
		{
			name: "Invalid UserID",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				return repo
			},
			favoriteID:       1,
			userID:           0, // Invalid UserID
			description:      "Updated Description",
			expectedFavorite: domain.Favorite{},
			expectedError:    domain.ErrInvalidUserID,
		},
		{
			name: "Favorite Not Found",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				repo.EXPECT().
					UpdateDescription(gomock.Any(), 1, 1, "Updated Description").
					Return(
						domain.Favorite{},
						domain.ErrFavoriteNotFound,
					)
				return repo
			},
			favoriteID:       1,
			userID:           1,
			description:      "Updated Description",
			expectedFavorite: domain.Favorite{},
			expectedError:    domain.ErrFavoriteNotFound,
		},
		{
			name: "Internal Server Error",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				repo.EXPECT().
					UpdateDescription(gomock.Any(), 1, 1, "Updated Description").
					Return(
						domain.Favorite{},
						domain.ErrInternalServerError,
					)
				return repo
			},
			favoriteID:       1,
			userID:           1,
			description:      "Updated Description",
			expectedFavorite: domain.Favorite{},
			expectedError:    domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			favoriteService := NewFavoriteService(tt.favoriteRepository(ctrl))

			favorite, err := favoriteService.Edit(context.Background(), tt.favoriteID, tt.userID, tt.description)

			if err != nil {
				require.ErrorIs(t, err, tt.expectedError)
				require.Equal(t, tt.expectedFavorite, favorite)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedFavorite, favorite)
			}
		})
	}

}

func TestFavoriteService_Delete(t *testing.T) {
	tests := []struct {
		name               string
		favoriteRepository func(mock *gomock.Controller) domain.FavoriteRepository
		favoriteID         int
		userID             int
		expectedError      error
	}{
		{
			name: "Valid Favorite Deletion",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				repo.EXPECT().
					Remove(gomock.Any(), 1, 1).
					Return(nil)
				return repo
			},
			favoriteID:    1,
			userID:        1,
			expectedError: nil,
		},
		{
			name: "Invalid FavoriteID",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				return repo
			},
			favoriteID:    0, // Invalid FavoriteID
			userID:        1,
			expectedError: domain.ErrInvalidFavoriteID,
		},
		{
			name: "Invalid UserID",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				return repo
			},
			favoriteID:    1,
			userID:        0, // Invalid UserID
			expectedError: domain.ErrInvalidUserID,
		},
		{
			name: "Favorite Not Found",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				repo.EXPECT().
					Remove(gomock.Any(), 1, 1).
					Return(domain.ErrFavoriteNotFound)
				return repo
			},
			favoriteID:    1,
			userID:        1,
			expectedError: domain.ErrFavoriteNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			favoriteService := NewFavoriteService(tt.favoriteRepository(ctrl))

			err := favoriteService.Delete(context.Background(), tt.favoriteID, tt.userID)

			if err != nil {
				require.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestFavoriteService_GetUserFavorites(t *testing.T) {
	tests := []struct {
		name               string
		favoriteRepository func(mock *gomock.Controller) domain.FavoriteRepository
		userID             int
		expectedFavorites  []domain.Favorite
		expectedError      error
	}{
		{
			name: "Valid User Favorites Retrieval",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				repo.EXPECT().
					SearchByUserID(gomock.Any(), 1).
					Return(
						[]domain.Favorite{
							{ID: 1, UserID: 1, AssetID: 1, Description: "Test"},
						},
						nil,
					)
				return repo
			},
			userID: 1,
			expectedFavorites: []domain.Favorite{
				{ID: 1, UserID: 1, AssetID: 1, Description: "Test"},
			},
			expectedError: nil,
		},
		{
			name: "Invalid UserID",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				return repo
			},
			userID:            0, // Invalid UserID
			expectedFavorites: nil,
			expectedError:     domain.ErrInvalidUserID,
		},
		{
			name: "No Favorites Found",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				repo.EXPECT().
					SearchByUserID(gomock.Any(), 2).
					Return([]domain.Favorite{}, nil)
				return repo
			},
			userID:            2,
			expectedFavorites: []domain.Favorite{},
			expectedError:     nil,
		},
		{
			name: "Internal Server Error",
			favoriteRepository: func(mock *gomock.Controller) domain.FavoriteRepository {
				repo := mocks.NewMockFavoriteRepository(mock)
				repo.EXPECT().
					SearchByUserID(gomock.Any(), 1).
					Return(nil, domain.ErrInternalServerError)
				return repo
			},
			userID:            1,
			expectedFavorites: nil,
			expectedError:     domain.ErrInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			favoriteService := NewFavoriteService(tt.favoriteRepository(ctrl))

			favorites, err := favoriteService.GetUserFavorites(context.Background(), tt.userID)

			if err != nil {
				require.ErrorIs(t, err, tt.expectedError)
				require.Equal(t, tt.expectedFavorites, favorites)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedFavorites, favorites)
			}
		})
	}
}
