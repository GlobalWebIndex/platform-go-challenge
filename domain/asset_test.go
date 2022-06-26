package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddAsset(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	usr := &User{
		ID:       1,
		Username: "manos",
		Password: "",
		IsAdmin:  false,
	}
	asset := Asset{
		Data: CorrectInputTestAssetData[0],
	}
	_, err := dom.AddAsset(ctx, usr, asset)
	assert.ErrorIs(t, err, ErrUnauthorized)

	usr = &User{
		ID:       1,
		Username: "manos",
		Password: "",
		IsAdmin:  true,
	}
	for _, v := range CorrectInputTestAssetData {
		asset := Asset{
			Data: v,
		}
		_, err := dom.AddAsset(ctx, usr, asset)
		assert.NoError(t, err)
	}
	for _, v := range WrongInputTestAssetData {
		_, err := dom.AddAsset(ctx, usr, Asset{
			Data: v,
		})
		assert.ErrorIs(t, err, ErrWrongAssetInput)
	}
}

func TestUpdateAsset(t *testing.T) {
	dom := NewDomain(&MockDB{})
	asset := Asset{
		Data: &Insight{
			Text:        "40% of millenials spend more than 3hours on social media daily",
			Description: "bla bla",
		},
	}
	ctx := context.Background()
	usr := &User{
		ID:       1,
		Username: "manos",
		Password: "",
		IsAdmin:  false,
	}
	_, err := dom.AddAsset(ctx, usr, asset)
	assert.ErrorIs(t, err, ErrUnauthorized)

	usr = &User{
		ID:       1,
		Username: "manos",
		Password: "",
		IsAdmin:  true,
	}
	_, err = dom.AddAsset(ctx, usr, asset)
	assert.NoError(t, err)
	for _, v := range WrongInputTestAssetData {
		_, err = dom.UpdateAsset(ctx, usr, Asset{
			ID:   asset.ID,
			Data: v,
		})
		assert.ErrorIs(t, err, ErrWrongAssetInput)
	}
}

func TestListAsset(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	usr := &User{
		ID:       1,
		Username: "manos",
		Password: "",
		IsAdmin:  false,
	}
	for _, v := range WrongInputTestQueryData {
		_, err := dom.ListAssets(ctx, usr, v, nil)
		assert.ErrorIs(t, err, ErrWrongQueryInput)
	}
}
