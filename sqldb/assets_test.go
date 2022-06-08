package sqldb

import (
	"context"
	"log"
	"platform-go-challenge/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupSuite(tb testing.TB) (*DB, func(tb testing.TB)) {
	log.Println("setup suite")
	db, err := NewDB("user", "user", "127.0.0.1:3306", "mydb")
	if err != nil {
		tb.Fatal(err)
	}
	db.dropTablesIfExist()
	sqldb, _ := db.db.DB()
	db.createTables()
	// Return a function to teardown the test
	return db, func(tb testing.TB) {
		db.dropTablesIfExist()
		sqldb.Close()
	}
}

func TestCRUDInsight(t *testing.T) {
	db, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	asset, err := db.AddAsset(ctx, domain.Asset{
		Data: &domain.Insight{
			Text:        "40% of millenials spend more than 3hours on social media daily",
			Description: "example",
		}})
	assert.NotNil(t, asset)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), asset.ID)
	assert.Equal(t, "example", asset.Data.(*domain.Insight).Description)

	asset, err = db.UpdateAsset(ctx, domain.Asset{
		ID: 1,
		Data: &domain.Insight{
			Text:        "100% of millenials spend more than 3hours on social media daily",
			Description: "updated example",
		}})
	assert.NotNil(t, asset)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), asset.ID)
	assert.Equal(t, "updated example", asset.Data.(*domain.Insight).Description)

	gottenAsset, err := db.GetAsset(ctx, domain.InsightAssetType, asset.ID)
	assert.NotNil(t, gottenAsset)
	assert.NoError(t, err)
	assert.EqualValues(t, asset, gottenAsset)
	err = db.DeleteAsset(ctx, domain.InsightAssetType, asset.ID)
	assert.Nil(t, err)

	_, err = db.GetAsset(ctx, domain.InsightAssetType, asset.ID)
	assert.NotNil(t, err)
}

func TestCRUDChart(t *testing.T) {
	db, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	asset, err := db.AddAsset(ctx, domain.Asset{
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
	assert.NotNil(t, asset)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), asset.ID)
	assert.Equal(t, "bla bla", asset.Data.(*domain.Chart).Description)

	asset, err = db.UpdateAsset(ctx, domain.Asset{
		ID: 1,
		Data: &domain.Chart{
			Description: "bla bla 2",
			Title:       "Relationship between tax and GDP",
			XTitle:      "GDP",
			YTitle:      "Tax",
			Data: domain.XYData{
				X: []interface{}{1, 2, 3, 4, 5},
				Y: []interface{}{1, 2, 3, 4, 5},
			},
		}})
	assert.NotNil(t, asset)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), asset.ID)
	assert.Equal(t, "bla bla 2", asset.Data.(*domain.Chart).Description)

	gottenAsset, err := db.GetAsset(ctx, domain.ChartAssetType, asset.ID)
	assert.NotNil(t, gottenAsset)
	assert.NoError(t, err)
	assert.EqualValues(t, asset, gottenAsset)
	err = db.DeleteAsset(ctx, domain.ChartAssetType, asset.ID)
	assert.Nil(t, err)

	_, err = db.GetAsset(ctx, domain.ChartAssetType, asset.ID)
	assert.NotNil(t, err)
}

func TestCRUDAudience(t *testing.T) {
	db, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
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
	assert.NotNil(t, asset)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), asset.ID)
	assert.Equal(t, "bla bla", asset.Data.(*domain.Audience).Description)

	asset, err = db.UpdateAsset(ctx, domain.Asset{
		ID: 1,
		Data: &domain.Audience{
			AgeMax:            30,
			AgeMin:            20,
			Gender:            domain.FemaleGenderType,
			Country:           "Sweden",
			HoursSpent:        3,
			NumberOfPurchases: 3,
			Description:       "bla bla 2",
		}})
	assert.NotNil(t, asset)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), asset.ID)
	assert.Equal(t, "bla bla 2", asset.Data.(*domain.Audience).Description)

	gottenAsset, err := db.GetAsset(ctx, domain.AudienceAssetType, asset.ID)
	assert.NotNil(t, gottenAsset)
	assert.NoError(t, err)
	assert.EqualValues(t, asset, gottenAsset)
	err = db.DeleteAsset(ctx, domain.AudienceAssetType, asset.ID)
	assert.Nil(t, err)

	_, err = db.GetAsset(ctx, domain.AudienceAssetType, asset.ID)
	assert.NotNil(t, err)
}
