package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAssetUnauthorizedFailure(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	usr := &User{
		ID:       1,
		Username: "manos",
		IsAdmin:  false,
	}
	asset := Asset{
		Data: CorrectInputTestAssetData[0],
	}
	_, err := dom.AddAsset(ctx, usr, asset)
	assert.ErrorIs(t, err, ErrUnauthorized)

}

func TestAddAssetWrongInputFailure(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	usr := &User{
		ID:       1,
		Username: "manos",
		IsAdmin:  true,
	}
	for _, v := range WrongInputTestAssetData {
		_, err := dom.AddAsset(ctx, usr, Asset{
			Data: v,
		})
		assert.ErrorIs(t, err, ErrWrongAssetInput)
	}

}

func TestAddAssetSuccess(t *testing.T) {
	mdb := &MockDB{}
	mdb.addAsset = func(ctx context.Context, asset Asset) (*Asset, error) {
		asset.ID = 1
		return &asset, nil
	}
	dom := NewDomain(mdb)
	ctx := context.Background()
	usr := &User{
		ID:       1,
		Username: "manos",
		IsAdmin:  true,
	}
	for _, v := range CorrectInputTestAssetData {
		asset := Asset{
			Data: v,
		}
		newAsset, err := dom.AddAsset(ctx, usr, asset)
		assert.NoError(t, err)
		assert.NotNil(t, newAsset)
		assert.EqualValues(t, newAsset.Data, asset.Data)
		assert.Equal(t, newAsset.ID, uint(1))
	}
}

func TestUpdateAssetUnauthorizedFailure(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	asset := Asset{
		Data: &Insight{
			Text:        "40% of millenials spend more than 3hours on social media daily",
			Description: "bla bla",
		},
	}

	usr := &User{
		ID:       1,
		Username: "manos",
		IsAdmin:  false,
	}
	_, err := dom.AddAsset(ctx, usr, asset)
	assert.ErrorIs(t, err, ErrUnauthorized)

}

func TestUpdateAssetWrongInputFailure(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	usr := &User{
		ID:       1,
		Username: "manos",
		IsAdmin:  true,
	}
	for _, v := range WrongInputTestAssetData {
		_, err := dom.UpdateAsset(ctx, usr, Asset{
			ID:   1,
			Data: v,
		})
		assert.ErrorIs(t, err, ErrWrongAssetInput)
	}
}

func TestUpdateAssetSuccess(t *testing.T) {
	mdb := &MockDB{}
	mdb.updateAsset = func(ctx context.Context, asset Asset) (*Asset, error) {
		asset.ID = 1
		return &asset, nil
	}
	dom := NewDomain(mdb)
	ctx := context.Background()
	usr := &User{
		ID:       1,
		Username: "manos",
		IsAdmin:  true,
	}
	for _, v := range CorrectInputTestAssetData {
		newAsset, err := dom.UpdateAsset(ctx, usr, Asset{
			ID:   1,
			Data: v,
		})
		assert.NoError(t, err)
		assert.NotNil(t, newAsset)
		assert.EqualValues(t, newAsset.Data, v)
		assert.Equal(t, newAsset.ID, uint(1))
	}
}

func TestListAssetFailure(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	usr := &User{
		ID:       1,
		Username: "manos",
		IsAdmin:  false,
	}
	for _, v := range WrongInputTestQueryData {
		_, err := dom.ListAssets(ctx, usr, v, nil)
		assert.ErrorIs(t, err, ErrWrongQueryInput)
	}
}

func TestListAssetSuccess(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	usr := &User{
		ID:       1,
		Username: "manos",
		IsAdmin:  false,
	}
	ls, err := dom.ListAssets(ctx, usr, QueryAssets{Limit: 10, LastID: 1, Type: AudienceAssetType}, nil)
	assert.NoError(t, err)
	assert.Empty(t, ls)
}
