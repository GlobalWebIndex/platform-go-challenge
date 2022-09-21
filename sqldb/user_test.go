package sqldb

import (
	"context"
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

	nuser, err := db.FindUser(ctx, user.Username)
	assert.NoError(t, err)
	assert.NotNil(t, nuser)
	assert.Equal(t, user, nuser)

	guser, err := db.GetUser(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, nuser)
	assert.Equal(t, user, guser)

	exists, err := db.UserExists(ctx, user.Username)
	assert.NoError(t, err)
	assert.NotNil(t, nuser)
	assert.True(t, exists)

	exists, err = db.UserExists(ctx, "none")
	assert.NoError(t, err)
	assert.NotNil(t, nuser)
	assert.False(t, exists)

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

	asset, err := db.AddAsset(ctx, domain.InputAsset{
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

	asset, err := db.AddAsset(ctx, domain.InputAsset{
		Data: &domain.Chart{
			Description: "bla bla",
			Title:       "Relationship between tax and GDP",
			XTitle:      "GDP",
			YTitle:      "Tax",
			Data: domain.XYData{
				X: []float64{1, 2, 3, 4, 5},
				Y: []float64{1, 2, 3, 4, 5},
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
