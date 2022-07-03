package intetests

import (
	"context"
	"fmt"
	"platform-go-challenge/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticationToListFavOfOtherUsers(t *testing.T) {
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

	user2, err := dom.CreateUser(ctx, domain.User{
		Username: "user2",
		Password: "password2",
		IsAdmin:  false,
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
			Description:       "example",
		}})
	assert.NotNil(t, asset)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), asset.ID)
	assert.Equal(t, "example", asset.Data.(*domain.Audience).Description)

	err = dom.FavouriteAsset(ctx, admin, asset.ID, domain.AudienceAssetType, true)
	assert.NoError(t, err)

	err = dom.FavouriteAsset(ctx, user, asset.ID, domain.AudienceAssetType, true)
	assert.NoError(t, err)

	err = dom.FavouriteAsset(ctx, user2, asset.ID, domain.AudienceAssetType, true)
	assert.NoError(t, err)

	qa := domain.QueryAssets{
		Limit:  10,
		LastID: 0,
		Type:   domain.AudienceAssetType,
		IsDesc: false,
	}

	favQuery := domain.QueryFavouriteAssets{
		FromUserID: user.ID,
		OnlyFav:    true,
	}
	la, err := dom.ListAssets(ctx, user, qa, &favQuery)
	fmt.Println(la)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(la.Assets))

	favQuery = domain.QueryFavouriteAssets{
		FromUserID: admin.ID,
		OnlyFav:    true,
	}
	la, err = dom.ListAssets(ctx, user, qa, &favQuery)
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
	assert.Nil(t, la)

	favQuery = domain.QueryFavouriteAssets{
		FromUserID: user2.ID,
		OnlyFav:    true,
	}
	la, err = dom.ListAssets(ctx, user, qa, &favQuery)
	assert.ErrorIs(t, err, domain.ErrUnauthorized)
	assert.Nil(t, la)

	favQuery = domain.QueryFavouriteAssets{
		FromUserID: admin.ID,
		OnlyFav:    true,
	}
	la, err = dom.ListAssets(ctx, admin, qa, &favQuery)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 1, len(la.Assets))

	favQuery = domain.QueryFavouriteAssets{
		FromUserID: user2.ID,
		OnlyFav:    true,
	}
	la, err = dom.ListAssets(ctx, admin, qa, &favQuery)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 1, len(la.Assets))

}

func TestListFavouriteAudiences(t *testing.T) {
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
		if i%2 == 0 {
			err := dom.FavouriteAsset(ctx, user, asset.ID, domain.AudienceAssetType, true)
			assert.NoError(t, err)
		}
	}
	qa := domain.QueryAssets{
		Limit:  10,
		LastID: 0,
		Type:   domain.AudienceAssetType,
		IsDesc: false,
	}
	favQuery := domain.QueryFavouriteAssets{
		FromUserID: user.ID,
		OnlyFav:    true,
	}
	la, err := dom.ListAssets(ctx, user, qa, &favQuery)
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
	la, err = dom.ListAssets(ctx, user, qa, &favQuery)
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
	favQuery = domain.QueryFavouriteAssets{
		FromUserID: user.ID,
		OnlyFav:    false,
	}

	la, err = dom.ListAssets(ctx, user, qa, &favQuery)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.False(t, la.Assets[0].IsFavourite)
	assert.True(t, la.Assets[1].IsFavourite)
	assert.Equal(t, uint(1), la.FirstID)
	assert.Equal(t, uint(10), la.LastID)
}

func TestListFavouriteInsights(t *testing.T) {
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
		if i%2 == 0 {
			err := dom.FavouriteAsset(ctx, user, asset.ID, domain.InsightAssetType, true)
			assert.NoError(t, err)
		}
	}
	qa := domain.QueryAssets{
		Limit:  10,
		LastID: 0,
		Type:   domain.InsightAssetType,
		IsDesc: false,
	}
	favQuery := domain.QueryFavouriteAssets{
		FromUserID: user.ID,
		OnlyFav:    true,
	}
	la, err := dom.ListAssets(ctx, user, qa, &favQuery)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.True(t, la.Assets[0].IsFavourite)
	assert.Equal(t, uint(2), la.FirstID)
	assert.Equal(t, uint(20), la.LastID)

	qa = domain.QueryAssets{
		Limit:  10,
		LastID: 101,
		Type:   domain.InsightAssetType,
		IsDesc: true,
	}
	la, err = dom.ListAssets(ctx, user, qa, &favQuery)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.True(t, la.Assets[0].IsFavourite)
	assert.Equal(t, uint(100), la.FirstID)
	assert.Equal(t, uint(82), la.LastID)

	qa = domain.QueryAssets{
		Limit:  10,
		LastID: 0,
		Type:   domain.InsightAssetType,
		IsDesc: false,
	}
	favQuery = domain.QueryFavouriteAssets{
		FromUserID: user.ID,
		OnlyFav:    false,
	}

	la, err = dom.ListAssets(ctx, user, qa, &favQuery)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.False(t, la.Assets[0].IsFavourite)
	assert.True(t, la.Assets[1].IsFavourite)
	assert.Equal(t, uint(1), la.FirstID)
	assert.Equal(t, uint(10), la.LastID)
}

func TestListFavouriteCharts(t *testing.T) {
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
		if i%2 == 0 {
			err := dom.FavouriteAsset(ctx, user, asset.ID, domain.ChartAssetType, true)
			assert.NoError(t, err)
		}
	}
	qa := domain.QueryAssets{
		Limit:  10,
		LastID: 0,
		Type:   domain.ChartAssetType,
		IsDesc: false,
	}
	favQuery := domain.QueryFavouriteAssets{
		FromUserID: user.ID,
		OnlyFav:    true,
	}
	la, err := dom.ListAssets(ctx, user, qa, &favQuery)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.True(t, la.Assets[0].IsFavourite)
	assert.Equal(t, uint(2), la.FirstID)
	assert.Equal(t, uint(20), la.LastID)

	qa = domain.QueryAssets{
		Limit:  10,
		LastID: 101,
		Type:   domain.ChartAssetType,
		IsDesc: true,
	}
	la, err = dom.ListAssets(ctx, user, qa, &favQuery)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.True(t, la.Assets[0].IsFavourite)
	assert.Equal(t, uint(100), la.FirstID)
	assert.Equal(t, uint(82), la.LastID)

	qa = domain.QueryAssets{
		Limit:  10,
		LastID: 0,
		Type:   domain.ChartAssetType,
		IsDesc: false,
	}
	favQuery = domain.QueryFavouriteAssets{
		FromUserID: user.ID,
		OnlyFav:    false,
	}

	la, err = dom.ListAssets(ctx, user, qa, &favQuery)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.False(t, la.Assets[0].IsFavourite)
	assert.True(t, la.Assets[1].IsFavourite)
	assert.Equal(t, uint(1), la.FirstID)
	assert.Equal(t, uint(10), la.LastID)
}
