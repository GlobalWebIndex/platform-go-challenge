package intetests

import (
	"context"
	"fmt"
	"platform-go-challenge/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListInsights(t *testing.T) {
	dom, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	admin, err := dom.CreateUser(ctx, domain.User{
		Username: "admin",
		Password: "password",
		IsAdmin:  true,
	})
	assert.NoError(t, err)

	user, err := dom.CreateUser(ctx, domain.User{
		Username: "user",
		Password: "password",
		IsAdmin:  false,
	})
	assert.NoError(t, err)
	for i := 1; i <= 100; i++ {
		desc := fmt.Sprintf("example %d", i)
		asset, err := dom.AddAsset(ctx, admin, domain.Asset{
			Data: &domain.Insight{
				Text:        "40% of millenials spend more than 3hours on social media daily",
				Description: desc,
			}})
		assert.NotNil(t, asset)
		assert.NoError(t, err)
		assert.Equal(t, uint(i), asset.ID)
		assert.Equal(t, desc, asset.Data.(*domain.Insight).Description)
	}
	qa := domain.QueryAssets{
		Limit:  10,
		LastID: 0,
		Type:   domain.InsightAssetType,
		IsDesc: false,
	}
	la, err := dom.ListAssets(ctx, user, qa, nil)
	fmt.Println(la)

	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.Equal(t, uint(1), la.FirstID)
	assert.Equal(t, uint(10), la.LastID)

	qa = domain.QueryAssets{
		Limit:  10,
		LastID: 101,
		Type:   domain.InsightAssetType,
		IsDesc: true,
	}
	la, err = dom.ListAssets(ctx, user, qa, nil)
	fmt.Println(la)

	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.Equal(t, uint(100), la.FirstID)
	assert.Equal(t, uint(91), la.LastID)
}

func TestListCharts(t *testing.T) {
	dom, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	admin, err := dom.CreateUser(ctx, domain.User{
		Username: "admin",
		Password: "password",
		IsAdmin:  true,
	})
	assert.NoError(t, err)

	user, err := dom.CreateUser(ctx, domain.User{
		Username: "user",
		Password: "password",
		IsAdmin:  false,
	})
	assert.NoError(t, err)
	for i := 1; i <= 100; i++ {
		desc := fmt.Sprintf("example %d", i)
		asset, err := dom.AddAsset(ctx, admin, domain.Asset{
			Data: &domain.Chart{
				Description: desc,
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
		assert.Equal(t, uint(i), asset.ID)
		assert.Equal(t, desc, asset.Data.(*domain.Chart).Description)
	}
	qa := domain.QueryAssets{
		Limit:  10,
		LastID: 0,
		Type:   domain.ChartAssetType,
		IsDesc: false,
	}
	la, err := dom.ListAssets(ctx, user, qa, nil)
	fmt.Println(la)

	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.Equal(t, uint(1), la.FirstID)
	assert.Equal(t, uint(10), la.LastID)

	qa = domain.QueryAssets{
		Limit:  10,
		LastID: 101,
		Type:   domain.ChartAssetType,
		IsDesc: true,
	}
	la, err = dom.ListAssets(ctx, user, qa, nil)
	fmt.Println(la)

	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.Equal(t, uint(100), la.FirstID)
	assert.Equal(t, uint(91), la.LastID)
}

func TestListAudiences(t *testing.T) {
	dom, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	admin, err := dom.CreateUser(ctx, domain.User{
		Username: "admin",
		Password: "password",
		IsAdmin:  true,
	})
	assert.NoError(t, err)

	user, err := dom.CreateUser(ctx, domain.User{
		Username: "user",
		Password: "password",
		IsAdmin:  false,
	})
	assert.NoError(t, err)
	for i := 1; i <= 100; i++ {
		desc := fmt.Sprintf("example %d", i)
		asset, err := dom.AddAsset(ctx, admin, domain.Asset{
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
	}
	qa := domain.QueryAssets{
		Limit:  10,
		LastID: 0,
		Type:   domain.AudienceAssetType,
		IsDesc: false,
	}
	la, err := dom.ListAssets(ctx, user, qa, nil)
	fmt.Println(la)

	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.Equal(t, uint(1), la.FirstID)
	assert.Equal(t, uint(10), la.LastID)

	qa = domain.QueryAssets{
		Limit:  10,
		LastID: 101,
		Type:   domain.AudienceAssetType,
		IsDesc: true,
	}
	la, err = dom.ListAssets(ctx, user, qa, nil)
	fmt.Println(la)

	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.Equal(t, uint(100), la.FirstID)
	assert.Equal(t, uint(91), la.LastID)
}
