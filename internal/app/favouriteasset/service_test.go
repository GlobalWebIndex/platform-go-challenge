package favouriteasset_test

import (
	"context"
	"errors"
	"platform-go-challenge/internal/app/favouriteasset"
	rfavouriteasset "platform-go-challenge/internal/infra/repository/mongo/favouriteasset"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestNewService(t *testing.T) {
	repoMock := &rfavouriteasset.RepoMock{}
	s := favouriteasset.NewService(repoMock)

	assert.IsType(t, new(favouriteasset.Service), &s)
}

func Test_assetService_GetFavouriteAssets(t *testing.T) {
	favouriteAssetID := "favourite asset ID"

	tests := map[string]struct {
		repoMock *rfavouriteasset.RepoMock
		ctx      context.Context
		userID   string
		expRes   *favouriteasset.GetFavouriteAssetsRes
		expErr   error
	}{
		"Should get favourite assets": {
			repoMock: func() *rfavouriteasset.RepoMock {
				userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")

				repoMock := &rfavouriteasset.RepoMock{}
				repoMock.On("GetFavouriteAssets", context.TODO(), userID).
					Return(&[]favouriteasset.FavouriteAsset{
						{
							ID:          &favouriteAssetID,
							Description: "favourite asset description",
							Asset: &favouriteasset.Asset{
								ID:   "asset id",
								Type: "asset type",
								Data: nil,
							},
						},
					}, nil)

				return repoMock
			}(),
			ctx:    context.TODO(),
			userID: "6377c9c68c9f75c3d2a8a5ba",
			expRes: &favouriteasset.GetFavouriteAssetsRes{
				FavouriteAssets: &[]favouriteasset.FavouriteAsset{
					{
						ID:          &favouriteAssetID,
						Description: "favourite asset description",
						Asset: &favouriteasset.Asset{
							ID:   "asset id",
							Type: "asset type",
							Data: nil,
						},
					},
				},
			},
		},
		"Should return error on GetFavouriteAssets error": {
			repoMock: func() *rfavouriteasset.RepoMock {
				userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")

				repoMock := &rfavouriteasset.RepoMock{}
				repoMock.On("GetFavouriteAssets", context.TODO(), userID).
					Return(&[]favouriteasset.FavouriteAsset{}, errors.New("random error"))

				return repoMock
			}(),
			ctx:    context.TODO(),
			userID: "6377c9c68c9f75c3d2a8a5ba",
			expErr: errors.New("random error"),
		},
		"Should return error on malformed user_id error": {
			ctx:    context.TODO(),
			userID: "malformed id",
			expErr: errors.New("the provided hex string is not a valid ObjectID"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aSvc := favouriteasset.NewService(tt.repoMock)
			res, err := aSvc.GetFavouriteAssets(tt.ctx, tt.userID)
			assert.Equal(t, tt.expRes, res)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func Test_assetService_AddAssetToFavourites(t *testing.T) {
	favouriteAssetID := "6378a0ae5893f7f7cc2770a9"

	tests := map[string]struct {
		repoMock *rfavouriteasset.RepoMock
		ctx      context.Context
		userID   string
		assetID  string
		expRes   *favouriteasset.AddFavouriteAssetRes
		expErr   error
	}{
		"Should add asset to favourites": {
			repoMock: func() *rfavouriteasset.RepoMock {
				userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")
				assetID, _ := primitive.ObjectIDFromHex("6378a00d5893f7f7cc2770a3")
				pFavouriteAssetID, _ := primitive.ObjectIDFromHex(favouriteAssetID)

				repoMock := &rfavouriteasset.RepoMock{}
				repoMock.On("AddToFavourites", context.TODO(), userID, assetID).
					Return(&pFavouriteAssetID, nil)

				return repoMock
			}(),
			ctx:     context.TODO(),
			userID:  "6377c9c68c9f75c3d2a8a5ba",
			assetID: "6378a00d5893f7f7cc2770a3",
			expRes: &favouriteasset.AddFavouriteAssetRes{
				ID: favouriteAssetID,
			},
		},
		"Should return error on AddToFavourites error": {
			repoMock: func() *rfavouriteasset.RepoMock {
				userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")
				assetID, _ := primitive.ObjectIDFromHex("6378a00d5893f7f7cc2770a3")

				repoMock := &rfavouriteasset.RepoMock{}
				repoMock.On("AddToFavourites", context.TODO(), userID, assetID).
					Return(&primitive.ObjectID{}, errors.New("AddToFavourites error"))

				return repoMock
			}(),
			ctx:     context.TODO(),
			userID:  "6377c9c68c9f75c3d2a8a5ba",
			assetID: "6378a00d5893f7f7cc2770a3",
			expErr:  errors.New("AddToFavourites error"),
		},
		"Should return error on malformed user_id error": {
			ctx:     context.TODO(),
			userID:  "malformed id",
			assetID: "6378a00d5893f7f7cc2770a3",
			expErr:  errors.New("the provided hex string is not a valid ObjectID"),
		},
		"Should return error on malformed asset_id error": {
			ctx:     context.TODO(),
			userID:  "6377c9c68c9f75c3d2a8a5ba",
			assetID: "malformed id",
			expErr:  errors.New("the provided hex string is not a valid ObjectID"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aSvc := favouriteasset.NewService(tt.repoMock)
			res, err := aSvc.AddAssetToFavourites(tt.ctx, tt.userID, tt.assetID)
			assert.Equal(t, tt.expRes, res)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func Test_assetService_EditFavouriteAsset(t *testing.T) {
	favouriteAssetID := "6378a0ae5893f7f7cc2770a9"

	tests := map[string]struct {
		repoMock *rfavouriteasset.RepoMock
		ctx      context.Context
		userID   string
		fAssetID string
		fAsset   favouriteasset.EditFavouriteAsset
		expErr   error
	}{
		"Should edit favourite asset": {
			repoMock: func() *rfavouriteasset.RepoMock {
				userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")
				pFavouriteAssetID, _ := primitive.ObjectIDFromHex(favouriteAssetID)
				fAsset := favouriteasset.EditFavouriteAsset{
					Description: "test description",
				}

				repoMock := &rfavouriteasset.RepoMock{}
				repoMock.On("UpdateFavouriteAsset", context.TODO(), userID, pFavouriteAssetID, fAsset).
					Return(nil)

				return repoMock
			}(),
			ctx:      context.TODO(),
			userID:   "6377c9c68c9f75c3d2a8a5ba",
			fAssetID: favouriteAssetID,
			fAsset: favouriteasset.EditFavouriteAsset{
				Description: "test description",
			},
		},
		"Should return error on EditFavouriteAsset error": {
			repoMock: func() *rfavouriteasset.RepoMock {
				userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")
				fassetID, _ := primitive.ObjectIDFromHex(favouriteAssetID)
				fAsset := favouriteasset.EditFavouriteAsset{}

				repoMock := &rfavouriteasset.RepoMock{}
				repoMock.On("UpdateFavouriteAsset", context.TODO(), userID, fassetID, fAsset).
					Return(errors.New("UpdateFavouriteAsset error"))

				return repoMock
			}(),
			ctx:      context.TODO(),
			userID:   "6377c9c68c9f75c3d2a8a5ba",
			fAssetID: favouriteAssetID,
			expErr:   errors.New("UpdateFavouriteAsset error"),
		},
		"Should return error on malformed user_id error": {
			ctx:      context.TODO(),
			userID:   "malformed id",
			fAssetID: favouriteAssetID,
			expErr:   errors.New("the provided hex string is not a valid ObjectID"),
		},
		"Should return error on malformed asset_id error": {
			ctx:      context.TODO(),
			userID:   "6377c9c68c9f75c3d2a8a5ba",
			fAssetID: "malformed id",
			expErr:   errors.New("the provided hex string is not a valid ObjectID"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aSvc := favouriteasset.NewService(tt.repoMock)
			err := aSvc.EditFavouriteAsset(tt.ctx, tt.userID, tt.fAssetID, tt.fAsset)
			assert.Equal(t, tt.expErr, err)
		})
	}
}

func Test_assetService_RemoveAssetFromFavourites(t *testing.T) {
	favouriteAssetID := "6378a0ae5893f7f7cc2770a9"

	tests := map[string]struct {
		repoMock *rfavouriteasset.RepoMock
		ctx      context.Context
		userID   string
		fAssetID string
		expErr   error
	}{
		"Should remove favourite asset": {
			repoMock: func() *rfavouriteasset.RepoMock {
				userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")
				pFavouriteAssetID, _ := primitive.ObjectIDFromHex(favouriteAssetID)

				repoMock := &rfavouriteasset.RepoMock{}
				repoMock.On("RemoveAssetFromFavourites", context.TODO(), pFavouriteAssetID, userID).
					Return(nil)

				return repoMock
			}(),
			ctx:      context.TODO(),
			userID:   "6377c9c68c9f75c3d2a8a5ba",
			fAssetID: favouriteAssetID,
		},
		"Should return error on RemoveAssetFromFavourites error": {
			repoMock: func() *rfavouriteasset.RepoMock {
				userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")
				pFavouriteAssetID, _ := primitive.ObjectIDFromHex(favouriteAssetID)

				repoMock := &rfavouriteasset.RepoMock{}
				repoMock.On("RemoveAssetFromFavourites", context.TODO(), pFavouriteAssetID, userID).
					Return(errors.New("RemoveAssetFromFavourites error"))

				return repoMock
			}(),
			ctx:      context.TODO(),
			userID:   "6377c9c68c9f75c3d2a8a5ba",
			fAssetID: favouriteAssetID,
			expErr:   errors.New("RemoveAssetFromFavourites error"),
		},
		"Should return error on malformed user_id error": {
			ctx:      context.TODO(),
			userID:   "malformed id",
			fAssetID: favouriteAssetID,
			expErr:   errors.New("the provided hex string is not a valid ObjectID"),
		},
		"Should return error on malformed asset_id error": {
			ctx:      context.TODO(),
			userID:   "6377c9c68c9f75c3d2a8a5ba",
			fAssetID: "malformed id",
			expErr:   errors.New("the provided hex string is not a valid ObjectID"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			aSvc := favouriteasset.NewService(tt.repoMock)
			err := aSvc.RemoveAssetFromFavourites(tt.ctx, tt.fAssetID, tt.userID)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
