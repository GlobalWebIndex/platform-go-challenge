package favouriteasset_test

import (
	"context"
	"platform-go-challenge/internal/app/favouriteasset"
	rfavouriteasset "platform-go-challenge/internal/infra/repository/mongo/favouriteasset"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestRepository_GetFavouriteAssets(t *testing.T) {
	userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")

	id1 := primitive.NewObjectID()
	id2 := primitive.NewObjectID()
	id3 := primitive.NewObjectID()
	id4 := primitive.NewObjectID()
	id5 := primitive.NewObjectID()

	id1Str := id1.Hex()
	id2Str := id2.Hex()
	id3Str := id3.Hex()
	id4Str := id4.Hex()
	id5Str := id5.Hex()

	chartData := bson.A{
		bson.D{
			{"_id", id5Str},
			{"type", "Chart"},
			{"data",
				bson.D{
					{"title", "chart title"},
					{"x_title", "chart x_title"},
					{"y_title", "chart y_title"},
					{"data", "chart data"},
				}},
		},
	}
	insightData := bson.A{
		bson.D{
			{"_id", id3Str},
			{"type", "Insight"},
			{"data", bson.D{{"text", "test insight text"}}},
		},
	}
	audienceData := bson.A{
		bson.D{
			{"_id", id4Str},
			{"type", "Audience"},
			{"data", bson.D{
				{"gender", "male"},
				{"birth_country", "Greece"},
				{"age_groups", "18-24"},
				{"hours_spent_daily_on_social_media", "100"},
				{"number_of_purchases_last_month", "2"},
			},
			},
		},
	}

	tests := map[string]struct {
		mongoRes *[]bson.D
		userID   primitive.ObjectID
		expRes   *[]favouriteasset.FavouriteAsset
		expErr   *mongo.CommandError
	}{
		"should get favourite assets": {
			mongoRes: &[]bson.D{
				mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
					{"_id", id1},
					{"description", "desc 1"},
					{"user_favourite_assets", insightData},
				}),
				mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
					{"_id", id2},
					{"description", "desc 2"},
					{"user_favourite_assets", audienceData},
				}),
				mtest.CreateCursorResponse(1, "foo.bar", mtest.NextBatch, bson.D{
					{"_id", id3},
					{"description", "desc 3"},
					{"user_favourite_assets", chartData},
				}),
				mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch),
			},
			userID: userID,
			expRes: &[]favouriteasset.FavouriteAsset{
				{
					ID:          &id1Str,
					Description: "desc 1",
					Asset: &favouriteasset.Asset{
						ID:   id3Str,
						Type: "Insight",
						Data: rfavouriteasset.InsightAsset{
							Text: "test insight text",
						},
					},
				},
				{
					ID:          &id2Str,
					Description: "desc 2",
					Asset: &favouriteasset.Asset{
						ID:   id4Str,
						Type: "Audience",
						Data: rfavouriteasset.AudienceAsset{
							Gender:                       "male",
							BirthCountry:                 "Greece",
							AgeGroups:                    "18-24",
							HoursSpentDailyOnSocialMedia: "100",
							NumberOfPurchasesLastMonth:   "2",
						},
					},
				},
				{
					ID:          &id3Str,
					Description: "desc 3",
					Asset: &favouriteasset.Asset{
						ID:   id5Str,
						Type: "Chart",
						Data: rfavouriteasset.ChartAsset{
							Title:  "chart title",
							XTitle: "chart x_title",
							YTitle: "chart y_title",
							Data:   "chart data",
						},
					},
				},
			},
		},
		"should handle error on mongo error": {
			mongoRes: &[]bson.D{
				mtest.CreateCommandErrorResponse(mtest.CommandError{
					Code:    1,
					Message: "random error",
					Name:    "Error",
				}),
			},
			userID: userID,
			expErr: &mongo.CommandError{
				Code:    1,
				Message: "random error",
				Name:    "Error",
			},
		},
	}
	for name, tt := range tests {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		defer mt.Close()
		mt.Run(name, func(mt *mtest.T) {
			rep, err := rfavouriteasset.NewRepository(mt.DB)
			assert.Nil(t, err)

			mt.AddMockResponses(*tt.mongoRes...)

			fAssets, err := rep.GetFavouriteAssets(context.TODO(), tt.userID)
			assert.Equal(t, tt.expRes, fAssets)
			if tt.expErr != nil {
				assert.ErrorContains(t, *tt.expErr, err.Error())
			}
		})
	}
}

func TestRepository_AddToFavourites(t *testing.T) {
	userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")
	assetID, _ := primitive.ObjectIDFromHex("6378a00d5893f7f7cc2770a3")

	tests := map[string]struct {
		mongoRes        *[]bson.D
		userID, assetID primitive.ObjectID
		expErr          *mongo.CommandError
	}{
		"should add asset to favourites": {
			mongoRes: &[]bson.D{
				mtest.CreateSuccessResponse(),
			},
			userID:  userID,
			assetID: assetID,
		},
	}
	for name, tt := range tests {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		defer mt.Close()
		mt.Run(name, func(mt *mtest.T) {
			rep, err := rfavouriteasset.NewRepository(mt.DB)
			assert.Nil(t, err)

			mt.AddMockResponses(*tt.mongoRes...)

			fAssets, err := rep.AddToFavourites(context.TODO(), tt.userID, tt.assetID)
			if tt.expErr == nil {
				assert.NotNil(t, fAssets)
			} else {
				assert.ErrorContains(t, *tt.expErr, err.Error())
			}
		})
	}
}

func TestRepository_UpdateFavouriteAsset(t *testing.T) {
	userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")
	assetID, _ := primitive.ObjectIDFromHex("6378a00d5893f7f7cc2770a3")

	tests := map[string]struct {
		mongoRes        *[]bson.D
		userID, assetID primitive.ObjectID
		fAsset          favouriteasset.EditFavouriteAsset
		expErr          *mongo.CommandError
	}{
		"should update favourite asset": {
			mongoRes: &[]bson.D{
				mtest.CreateSuccessResponse(),
			},
			userID:  userID,
			assetID: assetID,
			fAsset: favouriteasset.EditFavouriteAsset{
				Description: "test description",
			},
		},
	}
	for name, tt := range tests {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		defer mt.Close()
		mt.Run(name, func(mt *mtest.T) {
			rep, err := rfavouriteasset.NewRepository(mt.DB)
			assert.Nil(t, err)

			mt.AddMockResponses(*tt.mongoRes...)

			err = rep.UpdateFavouriteAsset(context.TODO(), tt.userID, tt.assetID, tt.fAsset)
			assert.Nil(t, err)
		})
	}
}

func TestRepository_RemoveAssetFromFavourites(t *testing.T) {
	userID, _ := primitive.ObjectIDFromHex("6377c9c68c9f75c3d2a8a5ba")
	assetID, _ := primitive.ObjectIDFromHex("6378a00d5893f7f7cc2770a3")

	tests := map[string]struct {
		mongoRes        *[]bson.D
		userID, assetID primitive.ObjectID
		expErr          *mongo.CommandError
	}{
		"should update favourite asset": {
			mongoRes: &[]bson.D{
				mtest.CreateSuccessResponse(),
			},
			userID:  userID,
			assetID: assetID,
		},
	}
	for name, tt := range tests {
		mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
		defer mt.Close()
		mt.Run(name, func(mt *mtest.T) {
			rep, err := rfavouriteasset.NewRepository(mt.DB)
			assert.Nil(t, err)

			mt.AddMockResponses(*tt.mongoRes...)

			err = rep.RemoveAssetFromFavourites(context.TODO(), tt.userID, tt.assetID)
			assert.Nil(t, err)
		})
	}
}
