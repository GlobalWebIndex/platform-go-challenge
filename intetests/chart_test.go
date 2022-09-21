package intetests

import (
	"context"
	"fmt"
	"platform-go-challenge/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListChartsSuccess(t *testing.T) {
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
		asset, err := dom.AddAsset(ctx, admin, domain.InputAsset{
			Data: &domain.Chart{
				Description: desc,
				Title:       "Relationship between tax and GDP",
				XTitle:      "GDP",
				YTitle:      "Tax",
				Data: domain.XYData{
					X: []float64{1, 2, 3, 4, 5},
					Y: []float64{1, 2, 3, 4, 5},
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

func TestCRUDChartSuccess(t *testing.T) {
	dom, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	admin, err := dom.CreateUser(ctx, domain.User{
		Username: "admin",
		Password: "password",
		IsAdmin:  true,
	})
	assert.NoError(t, err)
	asset, err := dom.AddAsset(ctx, admin, domain.InputAsset{
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
	assert.NotNil(t, asset)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), asset.ID)
	assert.Equal(t, "bla bla", asset.Data.(*domain.Chart).Description)

	asset, err = dom.UpdateAsset(ctx, admin, 1, domain.InputAsset{
		Data: &domain.Chart{
			Description: "bla bla 2",
			Title:       "Relationship between tax and GDP",
			XTitle:      "GDP",
			YTitle:      "Tax",
			Data: domain.XYData{
				X: []float64{1, 2, 3, 4, 5},
				Y: []float64{1, 2, 3, 4, 5},
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

func TestListFavouriteChartsSuccess(t *testing.T) {
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
		asset, err := dom.AddAsset(ctx, admin, domain.InputAsset{
			Data: &domain.Chart{
				Description: desc,
				Title:       "Relationship between tax and GDP",
				XTitle:      "GDP",
				YTitle:      "Tax",
				Data: domain.XYData{
					X: []float64{1, 2, 3, 4, 5},
					Y: []float64{1, 2, 3, 4, 5},
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
	assert.True(t, *la.Assets[0].IsFavourite)
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
	assert.True(t, *la.Assets[0].IsFavourite)
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
	assert.False(t, *la.Assets[0].IsFavourite)
	assert.True(t, *la.Assets[1].IsFavourite)
	assert.Equal(t, uint(1), la.FirstID)
	assert.Equal(t, uint(10), la.LastID)

	// Delete a favourite asset, and check the favourite of user
	favQuery = domain.QueryFavouriteAssets{
		FromUserID: user.ID,
		OnlyFav:    true,
	}
	la, err = dom.ListAssets(ctx, user, qa, &favQuery)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(la.Assets))
	assert.Equal(t, uint(2), la.Assets[0].ID)

	err = dom.DeleteAsset(ctx, admin, la.Assets[0].ID, domain.ChartAssetType)
	assert.NoError(t, err)

	favQuery = domain.QueryFavouriteAssets{
		FromUserID: user.ID,
		OnlyFav:    true,
	}
	la, err = dom.ListAssets(ctx, user, qa, &favQuery)
	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.Equal(t, uint(4), la.Assets[0].ID)
}
