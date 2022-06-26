package intetests

import (
	"context"
	"fmt"
	"platform-go-challenge/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRUDInsight(t *testing.T) {
	dom, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()

	admin, err := dom.CreateUser(ctx, domain.User{
		Username: "admin",
		Password: "password",
		IsAdmin:  true,
	})
	assert.NoError(t, err)
	fmt.Println(admin)
	asset, err := dom.AddAsset(ctx, admin, domain.Asset{
		Data: &domain.Insight{
			Text:        "40% of millenials spend more than 3hours on social media daily",
			Description: "example",
		}})
	assert.NotNil(t, asset)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), asset.ID)
	assert.Equal(t, "example", asset.Data.(*domain.Insight).Description)

	asset, err = dom.UpdateAsset(ctx, admin, domain.Asset{
		ID: 1,
		Data: &domain.Insight{
			Text:        "100% of millenials spend more than 3hours on social media daily",
			Description: "updated example",
		}})
	assert.NotNil(t, asset)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), asset.ID)
	assert.Equal(t, "updated example", asset.Data.(*domain.Insight).Description)

	gottenAsset, err := dom.GetAsset(ctx, admin, asset.ID, domain.InsightAssetType)
	assert.NotNil(t, gottenAsset)
	assert.NoError(t, err)
	assert.EqualValues(t, asset, gottenAsset)

	err = dom.DeleteAsset(ctx, admin, asset.ID, domain.InsightAssetType)
	assert.Nil(t, err)

	_, err = dom.GetAsset(ctx, admin, asset.ID, domain.InsightAssetType)
	assert.NotNil(t, err)
}

func TestCRUDChart(t *testing.T) {
	dom, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	admin, err := dom.CreateUser(ctx, domain.User{
		Username: "admin",
		Password: "password",
		IsAdmin:  true,
	})
	assert.NoError(t, err)
	asset, err := dom.AddAsset(ctx, admin, domain.Asset{
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

	asset, err = dom.UpdateAsset(ctx, admin, domain.Asset{
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

	gottenAsset, err := dom.GetAsset(ctx, admin, asset.ID, domain.ChartAssetType)
	assert.NotNil(t, gottenAsset)
	assert.NoError(t, err)
	assert.EqualValues(t, asset, gottenAsset)
	err = dom.DeleteAsset(ctx, admin, asset.ID, domain.ChartAssetType)
	assert.Nil(t, err)

	_, err = dom.GetAsset(ctx, admin, asset.ID, domain.ChartAssetType)
	assert.NotNil(t, err)
}

func TestCRUDAudience(t *testing.T) {
	dom, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	admin, err := dom.CreateUser(ctx, domain.User{
		Username: "admin",
		Password: "password",
		IsAdmin:  true,
	})
	assert.NoError(t, err)
	asset, err := dom.AddAsset(ctx, admin, domain.Asset{
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

	asset, err = dom.UpdateAsset(ctx, admin, domain.Asset{
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

	gottenAsset, err := dom.GetAsset(ctx, admin, asset.ID, domain.AudienceAssetType)
	assert.NotNil(t, gottenAsset)
	assert.NoError(t, err)
	assert.EqualValues(t, asset, gottenAsset)
	err = dom.DeleteAsset(ctx, admin, asset.ID, domain.AudienceAssetType)
	assert.Nil(t, err)

	_, err = dom.GetAsset(ctx, admin, asset.ID, domain.AudienceAssetType)
	assert.NotNil(t, err)
}
