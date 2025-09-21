package consts

import (
	"platform-go-challenge/domain"
	"time"
)

var Audience = []domain.Asset{
	domain.Asset{
		Type:        "audience",
		Description: "This is the first audience",
		CreatedAt:   time.Now(),
		Data: domain.AssetData{
			AudienceData: &domain.Audience{
				Gender:             "male",
				BirthCountry:       "greece",
				AgeGroup:           "25-35",
				SocialMediaHours:   3,
				PurchasesLastMonth: 5,
			},
		},
	},
	domain.Asset{
		Type:        "audience",
		Description: "This is the second audience",
		CreatedAt:   time.Now(),
		Data: domain.AssetData{
			AudienceData: &domain.Audience{
				Gender:             "female",
				BirthCountry:       "germany",
				AgeGroup:           "18-25",
				SocialMediaHours:   4,
				PurchasesLastMonth: 7,
			},
		},
	},
	domain.Asset{
		Type:        "audience",
		Description: "This is the third audience",
		CreatedAt:   time.Now(),
		Data: domain.AssetData{
			AudienceData: &domain.Audience{
				Gender:             "male",
				BirthCountry:       "italy",
				AgeGroup:           "35-45",
				SocialMediaHours:   2,
				PurchasesLastMonth: 2,
			},
		},
	},
}
