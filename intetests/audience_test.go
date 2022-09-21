package intetests

import (
	"context"
	"fmt"
	"platform-go-challenge/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRUDAudienceSuccess(t *testing.T) {
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

	asset, err = dom.UpdateAsset(ctx, admin, 1, domain.InputAsset{
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

func TestListAudiencesSuccess(t *testing.T) {
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

	assert.NoError(t, err)
	assert.NotNil(t, la)
	assert.Equal(t, 10, len(la.Assets))
	assert.Equal(t, uint(100), la.FirstID)
	assert.Equal(t, uint(91), la.LastID)
}

func TestListFavouriteAudiencesSuccess(t *testing.T) {
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
	assert.True(t, *la.Assets[0].IsFavourite)
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
	assert.True(t, *la.Assets[0].IsFavourite)
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

	err = dom.DeleteAsset(ctx, admin, la.Assets[0].ID, domain.AudienceAssetType)
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
