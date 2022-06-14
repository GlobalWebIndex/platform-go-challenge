package sqldb

import (
	"context"
	"fmt"
	"platform-go-challenge/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRUser(t *testing.T) {
	db, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	du := domain.User{
		Username: "manos",
		Password: "hashed",
		IsAdmin:  false,
	}
	user, err := db.AddUser(ctx, du)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)

	nuser, err := db.FindUser(ctx, domain.LoginCredentials{Username: user.Username, Password: user.Password})
	assert.NoError(t, err)
	assert.NotNil(t, nuser)
	assert.Equal(t, user, nuser)
}

func TestFavourInsight(t *testing.T) {
	db, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	du := domain.User{
		Username: "manos",
		Password: "hashed",
		IsAdmin:  false,
	}
	user, err := db.AddUser(ctx, du)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)

	asset, err := db.AddAsset(ctx, domain.Asset{
		Data: &domain.Insight{
			Text:        "40% of millenials spend more than 3hours on social media daily",
			Description: "example",
		}})
	assert.NoError(t, err)
	assert.NotNil(t, asset)
	fid, err := db.FavouriteAsset(ctx, user.ID, asset.ID, domain.InsightAssetType, true)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), fid)

	fid, err = db.FavouriteAsset(ctx, user.ID, asset.ID, domain.InsightAssetType, true)
	assert.Error(t, err)
	assert.Equal(t, uint(0), fid)

	fid, err = db.FavouriteAsset(ctx, user.ID, asset.ID, domain.InsightAssetType, false)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), fid)

}

func TestFavourChart(t *testing.T) {
	db, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	du := domain.User{
		Username: "manos",
		Password: "hashed",
		IsAdmin:  false,
	}
	user, err := db.AddUser(ctx, du)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)

	asset, err := db.AddAsset(ctx, domain.Asset{
		ID: 1,
		Data: &domain.Chart{
			Description: "bla bla",
			Title:       "Relationship between tax and GDP",
			XTitle:      "GDP",
			YTitle:      "Tax",
			Data: domain.XYData{
				X: []interface{}{1, 2, 3, 4, 5},
				Y: []interface{}{1, 2, 3, 4, 5},
			},
		}})
	assert.NoError(t, err)
	assert.NotNil(t, asset)
	fid, err := db.FavouriteAsset(ctx, user.ID, asset.ID, domain.ChartAssetType, true)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), fid)

	fid, err = db.FavouriteAsset(ctx, user.ID, asset.ID, domain.ChartAssetType, true)
	assert.Error(t, err)
	assert.Equal(t, uint(0), fid)

	fid, err = db.FavouriteAsset(ctx, user.ID, asset.ID, domain.ChartAssetType, false)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), fid)
}

func TestFavourAudience(t *testing.T) {
	db, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	du := domain.User{
		Username: "manos",
		Password: "hashed",
		IsAdmin:  false,
	}
	user, err := db.AddUser(ctx, du)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)

	asset, err := db.AddAsset(ctx, domain.Asset{
		Data: &domain.Audience{
			AgeMax:            30,
			AgeMin:            20,
			Gender:            domain.FemaleGenderType,
			Country:           "Sweden",
			HoursSpent:        3,
			NumberOfPurchases: 3,
			Description:       "bla bla",
		}})
	assert.NoError(t, err)
	assert.NotNil(t, asset)
	fid, err := db.FavouriteAsset(ctx, user.ID, asset.ID, domain.AudienceAssetType, true)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), fid)

	fid, err = db.FavouriteAsset(ctx, user.ID, asset.ID, domain.AudienceAssetType, true)
	assert.Error(t, err)
	assert.Equal(t, uint(0), fid)

	fid, err = db.FavouriteAsset(ctx, user.ID, asset.ID, domain.AudienceAssetType, false)
	assert.NoError(t, err)
	assert.Equal(t, uint(0), fid)
}

func TestListFavouriteAudiences(t *testing.T) {
	db, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	du := domain.User{
		Username: "manos",
		Password: "hashed",
		IsAdmin:  false,
	}
	user, err := db.AddUser(ctx, du)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	counter := 1
	for i := 1; i <= 100; i++ {
		desc := fmt.Sprintf("example %d", i)
		asset, err := db.AddAsset(ctx, domain.Asset{
			Data: &domain.Audience{
				AgeMax:            30,
				AgeMin:            20,
				Gender:            domain.FemaleGenderType,
				Country:           "Sweden",
				HoursSpent:        3,
				NumberOfPurchases: 3,
				Description:       desc,
			}})
		assert.NotNil(t, asset)
		assert.NoError(t, err)
		assert.Equal(t, uint(i), asset.ID)
		assert.Equal(t, desc, asset.Data.(*domain.Audience).Description)

		if i%2 == 0 {
			fid, err := db.FavouriteAsset(ctx, user.ID, asset.ID, domain.AudienceAssetType, true)
			assert.NoError(t, err)
			assert.Equal(t, uint(counter), fid)
			counter += 1
		}
	}
	qa := domain.QueryAssets{
		Limit:  10,
		LastID: 0,
		Type:   domain.AudienceAssetType,
		IsDesc: false,
	}
	la, err := db.ListFavouriteAssets(ctx, user.ID, true, qa)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.True(t, la.Assets[0].IsFavourite)
	assert.Equal(t, uint(2), la.FirstID)
	assert.Equal(t, uint(20), la.LastID)

	qa = domain.QueryAssets{
		Limit:  10,
		LastID: 101,
		Type:   domain.AudienceAssetType,
		IsDesc: true,
	}
	la, err = db.ListFavouriteAssets(ctx, user.ID, true, qa)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.True(t, la.Assets[0].IsFavourite)
	assert.Equal(t, uint(100), la.FirstID)
	assert.Equal(t, uint(82), la.LastID)

	qa = domain.QueryAssets{
		Limit:  10,
		LastID: 0,
		Type:   domain.AudienceAssetType,
		IsDesc: false,
	}

	la, err = db.ListFavouriteAssets(ctx, user.ID, false, qa)
	fmt.Println(la)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.False(t, la.Assets[0].IsFavourite)
	assert.True(t, la.Assets[1].IsFavourite)
	assert.Equal(t, uint(1), la.FirstID)
	assert.Equal(t, uint(10), la.LastID)
}
