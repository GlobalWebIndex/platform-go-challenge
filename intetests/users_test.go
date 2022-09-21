package intetests

import (
	"context"
	"fmt"
	"platform-go-challenge/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserSuccess(t *testing.T) {
	dom, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	user, err := dom.CreateUser(ctx, domain.User{
		Username: "manos",
		Password: "password",
		IsAdmin:  true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, user.IsAdmin, true)

	user2, err := dom.LoginUser(ctx, domain.LoginCredentials{
		Username: "manos",
		Password: "password",
	})
	assert.NoError(t, err)
	assert.NotNil(t, user2)
	assert.EqualValues(t, user, user2)

	_, err = dom.CreateUser(ctx, domain.User{
		Username: "manos",
		Password: "password",
		IsAdmin:  true,
	})
	assert.Error(t, err)
}

func TestAuthenticationToListFavOfOtherUsersSuccess(t *testing.T) {
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

	asset, err := dom.AddAsset(ctx, admin, domain.InputAsset{
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
