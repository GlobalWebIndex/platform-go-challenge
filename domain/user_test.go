package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestAddUserFailure(t *testing.T) {
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

func TestAddUserSuccess(t *testing.T) {
	mdb := &MockDB{}
	mdb.userExists = func(ctx context.Context, username string) (bool, error) {
		return false, nil
	}
	mdb.addUser = func(ctx context.Context, user User) (*User, error) {
		user.ID = 1
		return &user, nil
	}
	dom := NewDomain(mdb)
	ctx := context.Background()
	usr, err := dom.CreateUser(ctx, User{
		Username: "manos",
		Password: "secret",
		IsAdmin:  true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, usr)
	assert.Equal(t, usr.ID, uint(1))
	assert.Equal(t, usr.Username, "manos")
	assert.True(t, usr.IsAdmin)
}

func TestLoginUserWrongInputFailure(t *testing.T) {
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

func TestLoginUserDontExistsFailure(t *testing.T) {
	mdb := &MockDB{}
	mdb.findUser = func(ctx context.Context, username string) (*User, error) {
		return nil, gorm.ErrRecordNotFound
	}
	dom := NewDomain(mdb)
	ctx := context.Background()
	_, err := dom.LoginUser(ctx, LoginCredentials{
		Username: "manos",
		Password: "secret",
	})
	assert.ErrorIs(t, err, ErrUserNotFound)
}

func TestLoginUserWrongPasswordFailure(t *testing.T) {
	mdb := &MockDB{}
	mdb.findUser = func(ctx context.Context, username string) (*User, error) {
		hs, _ := hashPassword("wrong-secret")
		return &User{Username: "manos", Password: hs, ID: 1}, nil
	}
	dom := NewDomain(mdb)
	ctx := context.Background()
	_, err := dom.LoginUser(ctx, LoginCredentials{
		Username: "manos",
		Password: "secret",
	})
	assert.ErrorIs(t, err, ErrUnauthorized)
}

func TestLoginUserSuccess(t *testing.T) {
	mdb := &MockDB{}
	mdb.findUser = func(ctx context.Context, username string) (*User, error) {
		hs, _ := hashPassword("secret")
		return &User{Username: "manos", Password: hs, ID: 1}, nil
	}
	dom := NewDomain(mdb)
	ctx := context.Background()
	user, err := dom.LoginUser(ctx, LoginCredentials{
		Username: "manos",
		Password: "secret",
	})
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "manos", user.Username)
	assert.Equal(t, "", user.Password)
	assert.Equal(t, uint(1), user.ID)
}

func TestAuthenticationFailures(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	_, err := dom.AddAsset(ctx, nil, InputAsset{})
	assert.ErrorIs(t, err, ErrUnauthorized)
	_, err = dom.UpdateAsset(ctx, nil, 1, InputAsset{})
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
