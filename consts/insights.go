package consts

import (
	"platform-go-challenge/domain"
	"time"
)

var Insights = []domain.Asset{
	domain.Asset{
		Type:        "insight",
		Description: "This is the first insight",
		CreatedAt:   time.Now(),
		Data: domain.AssetData{
			InsightData: &domain.Insight{
				Text: "This is the first insight",
			},
		},
	},
	domain.Asset{
		Type:        "insight",
		Description: "This is the second insight",
		CreatedAt:   time.Now(),
		Data: domain.AssetData{
			InsightData: &domain.Insight{
				Text: "This is the second insight",
			},
		},
	},
	domain.Asset{
		Type:        "insight",
		Description: "This is the third insight",
		CreatedAt:   time.Now(),
		Data: domain.AssetData{
			InsightData: &domain.Insight{
				Text: "This is the third insight",
			},
		},
	},
}
