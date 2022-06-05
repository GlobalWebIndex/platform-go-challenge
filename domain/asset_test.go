package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var wrongInputTestAssetData = []interface{}{
	&Insight{
		Text:        "",
		Description: "bla bla",
	},
	&Insight{
		Text:        "40% of millenials spend more than 3hours on social media daily",
		Description: "",
	},
	&Insight{
		Text:        "",
		Description: "",
	},
	&Chart{
		Description: "",
		Title:       "Relationship between tax and GDP",
		XTitle:      "GDP",
		YTitle:      "Tax",
		Data: XYData{
			X: []interface{}{1, 2, 3, 4, 5},
			Y: []interface{}{1, 2, 3, 4, 5},
		},
	},
	&Chart{
		Description: "bla bla",
		Title:       "",
		XTitle:      "GDP",
		YTitle:      "Tax",
		Data: XYData{
			X: []interface{}{1, 2, 3, 4, 5},
			Y: []interface{}{1, 2, 3, 4, 5},
		},
	},
	&Chart{
		Description: "bla bla",
		Title:       "Relationship between tax and GDP",
		XTitle:      "",
		YTitle:      "Tax",
		Data: XYData{
			X: []interface{}{1, 2, 3, 4, 5},
			Y: []interface{}{1, 2, 3, 4, 5},
		},
	},
	&Chart{
		Description: "bla bla",
		Title:       "Relationship between tax and GDP",
		XTitle:      "GDP",
		YTitle:      "",
		Data: XYData{
			X: []interface{}{1, 2, 3, 4, 5},
			Y: []interface{}{1, 2, 3, 4, 5},
		},
	},
	&Chart{
		Description: "bla bla",
		Title:       "Relationship between tax and GDP",
		XTitle:      "GDP",
		YTitle:      "Tax",
		Data: XYData{
			X: []interface{}{},
			Y: []interface{}{},
		},
	},
	&Chart{
		Description: "bla bla",
		Title:       "Relationship between tax and GDP",
		XTitle:      "GDP",
		YTitle:      "Tax",
		Data: XYData{
			X: []interface{}{1, 2, 3, 4},
			Y: []interface{}{1, 2, 3, 4, 5},
		},
	},
	&Audience{
		AgeMax:            200,
		AgeMin:            0,
		Gender:            FemaleGenderType,
		Country:           "Sweden",
		HoursSpent:        3,
		NumberOfPurchases: 3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Gender:            GenderType("lalasdla"),
		Country:           "Sweden",
		HoursSpent:        3,
		NumberOfPurchases: 3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Country:           "Sweden",
		HoursSpent:        3,
		NumberOfPurchases: 3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Gender:            FemaleGenderType,
		Country:           "Mordor",
		HoursSpent:        3,
		NumberOfPurchases: 3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Gender:            FemaleGenderType,
		Country:           "Sweden",
		HoursSpent:        -3,
		NumberOfPurchases: 3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Gender:            FemaleGenderType,
		Country:           "Sweden",
		HoursSpent:        3,
		NumberOfPurchases: -3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Gender:            FemaleGenderType,
		Country:           "Sweden",
		HoursSpent:        3,
		NumberOfPurchases: 3,
	},
}

var wrongInputTestQueryData = []QueryAssets{
	{
		Limit:  0,
		LastID: 10,
	},
	{
		Limit:  10,
		LastID: 0,
	},
}

func TestAddAsset(t *testing.T) {
	dom := NewDomain(&MockDB{})
	asset := Asset{
		Data: &Insight{
			Text:        "40% of millenials spend more than 3hours on social media daily",
			Description: "bla bla",
		},
	}
	ctx := context.Background()
	err := dom.AddAsset(ctx, asset)
	assert.NoError(t, err)
	for _, v := range wrongInputTestAssetData {
		err = dom.AddAsset(ctx, Asset{
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
	err := dom.AddAsset(ctx, asset)
	assert.NoError(t, err)
	for _, v := range wrongInputTestAssetData {
		err = dom.UpdateAsset(ctx, 1, Asset{
			Data: v,
		})
		assert.ErrorIs(t, err, ErrWrongAssetInput)
	}
}

func TestListAsset(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	for _, v := range wrongInputTestQueryData {
		_, err := dom.ListAssets(ctx, 0, v)
		assert.ErrorIs(t, err, ErrWrongQueryInput)
	}
}
