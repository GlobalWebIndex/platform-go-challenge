package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	_, err := dom.CreateUser(ctx, User{
		Username: "manos",
		Password: "",
	})
	assert.ErrorIs(t, err, ErrWrongUserInput)
	_, err = dom.CreateUser(ctx, User{
		Username: "",
		Password: "secret",
	})
	assert.ErrorIs(t, err, ErrWrongUserInput)
}

func TestLoginUser(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	_, err := dom.LoginUser(ctx, LoginCredentials{
		Username: "manos",
		Password: "",
	})
	assert.ErrorIs(t, err, ErrWrongLoginInput)
	_, err = dom.LoginUser(ctx, LoginCredentials{
		Username: "",
		Password: "secret",
	})
	assert.ErrorIs(t, err, ErrWrongLoginInput)
}

func TestAuthentication(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	_, err := dom.AddAsset(ctx, nil, Asset{})
	assert.ErrorIs(t, err, ErrUnauthorized)
	_, err = dom.UpdateAsset(ctx, nil, Asset{})
	assert.ErrorIs(t, err, ErrUnauthorized)
	err = dom.DeleteAsset(ctx, nil, 1, AudienceAssetType)
	assert.ErrorIs(t, err, ErrUnauthorized)
	_, err = dom.GetAsset(ctx, nil, 1, AudienceAssetType)
	assert.ErrorIs(t, err, ErrUnauthorized)
	_, err = dom.ListAssets(ctx, nil, QueryAssets{}, nil)
	assert.ErrorIs(t, err, ErrUnauthorized)
	err = dom.FavouriteAsset(ctx, nil, 1, AudienceAssetType, false)
	assert.ErrorIs(t, err, ErrUnauthorized)
}
