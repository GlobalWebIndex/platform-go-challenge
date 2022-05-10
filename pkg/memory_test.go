package gwi

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddAssetToFavs(t *testing.T) {

	assert := assert.New(t)

	type input struct {
		userid string
		asset  Asset
	}

	tests := []struct {
		Description    string
		Input          input
		ExpectedOutput error
		ExpectedResult bool
	}{
		{
			Description: "when the asset does not previously exist",
			Input: input{
				userid: uuid.New().String(),
				asset: Asset{
					AssetType: "chart",
					Chart: Chart{
						Id: uuid.New().String(),
					},
				},
			},
			ExpectedOutput: nil,
			ExpectedResult: true,
		},
	}

	mer, err := NewMemoryRepository()
	assert.NoError(err)

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			output, err := mer.AddAssetToFavs(ctx, test.Input.userid, &test.Input.asset)
			assert.NoError(err)
			assert.Equal(output, test.ExpectedResult)
		})
	}
}

func TestUpdateFav(t *testing.T) {

	assert := assert.New(t)

	type input struct {
		userid string
		asset  Asset
	}
	tests := []struct {
		Description    string
		Input          input
		ExpectedOutput error
		ExpectedResult bool
	}{
		{
			Description: "when asset does not previously exists",
			Input: input{
				userid: uuid.New().String(),
				asset: Asset{
					AssetType: "chart",
					Chart: Chart{
						Id: uuid.New().String(),
					},
				},
			},
			ExpectedOutput: errors.New(""),
			ExpectedResult: false,
		},
		{
			Description: "when asset previously exists",
			Input: input{
				userid: uuid.New().String(),
				asset: Asset{
					AssetType: "chart",
					Chart: Chart{
						Id: uuid.New().String(),
					},
				},
			},
			ExpectedOutput: nil,
			ExpectedResult: true,
		},
	}

	mer, err := NewMemoryRepository()
	assert.NoError(err)

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			if test.ExpectedResult {
				mer.AddAssetToFavs(ctx, test.Input.userid, &test.Input.asset)
			}
			output, err := mer.UpdateFav(ctx, test.Input.userid, &test.Input.asset)
			if test.ExpectedResult {
				assert.NoError(err)
			} else {
				assert.NotNil(err)
			}
			assert.Equal(output, test.ExpectedResult)
		})
	}
}

func TestRemoveFav(t *testing.T) {

	assert := assert.New(t)

	type input struct {
		userid string
		asset  Asset
	}
	tests := []struct {
		Description    string
		Input          input
		ExpectedOutput error
		ExpectedResult bool
	}{
		{
			Description: "when asset does not previously exists",
			Input: input{
				userid: uuid.New().String(),
				asset: Asset{
					AssetType: "chart",
					Chart: Chart{
						Id: uuid.New().String(),
					},
				},
			},
			ExpectedOutput: errors.New(""),
			ExpectedResult: false,
		},
		{
			Description: "when asset previously exists",
			Input: input{
				userid: uuid.New().String(),
				asset: Asset{
					AssetType: "chart",
					Chart: Chart{
						Id: uuid.New().String(),
					},
				},
			},
			ExpectedOutput: nil,
			ExpectedResult: true,
		},
	}

	mer, err := NewMemoryRepository()
	assert.NoError(err)

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			if test.ExpectedResult {
				mer.AddAssetToFavs(ctx, test.Input.userid, &test.Input.asset)
			}
			output, err := mer.removeFav(ctx, test.Input.userid, &test.Input.asset)
			if test.ExpectedResult {
				assert.NoError(err)
			} else {
				assert.NotNil(err)
			}
			assert.Equal(output, test.ExpectedResult)
		})
	}
}

func TestRetrieveFavs(t *testing.T) {

	assert := assert.New(t)

	type input struct {
		userid string
		asset  Asset
	}

	tests := []struct {
		Description    string
		Input          input
		ExpectedOutput error
		ExpectedResult map[AssetType][]Asseter
	}{
		{
			Description: "retrieve empty list of favorites",
			Input: input{
				userid: uuid.New().String(),
			},
			ExpectedOutput: errors.New("no results"),
			ExpectedResult: make(map[AssetType][]Asseter),
		},
	}

	mer, err := NewMemoryRepository()
	assert.NoError(err)

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {
			if test.ExpectedOutput == nil {
				mer.AddAssetToFavs(ctx, test.Input.userid, &test.Input.asset)
			}
			output, err := mer.RetrieveFavs(ctx, test.Input.userid)
			assert.Equal(err, test.ExpectedOutput)
			assert.Equal(output, test.ExpectedResult)
		})
	}
}
