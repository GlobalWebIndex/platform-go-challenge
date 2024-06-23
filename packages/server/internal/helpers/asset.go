package helpers

import (
	"fmt"

	"github.com/alexandros-georgantas/platform-go-challenge/internal/models"
	"gorm.io/gorm"
)

func AggregateAssetType(a models.Asset, db gorm.DB) (*models.AssetResponse, error) {

	var assetResponse models.AssetResponse
	assetResponse.ID = a.ID
	assetResponse.RelatedID = a.RelatedID
	assetResponse.RelatedType = a.RelatedType
	assetResponse.Description = a.Description

	switch a.RelatedType {
	case "charts":
		var chart models.Chart
		if err := db.Where("id = ?", a.RelatedID).Find(&chart).Error; err != nil {
			return nil, fmt.Errorf("something went wrong while fetching chart with id %v of asset with id %v", a.RelatedID, a.ID)
		}
		assetResponse.Chart = &models.ChartWithoutAsset{ID: chart.ID, Title: chart.Title, HorizontalAxisLabel: chart.HorizontalAxisLabel, VerticalAxisLabel: chart.VerticalAxisLabel, Data: chart.Data}
	case "audiences":
		var audience models.Audience
		if err := db.Where("id = ?", a.RelatedID).Find(&audience).Error; err != nil {
			return nil, fmt.Errorf("something went wrong while fetching audience with id %v of asset with id %v", a.RelatedID, a.ID)
		}
		assetResponse.Audience = &models.AudienceWithoutAsset{ID: audience.ID, Gender: audience.Gender, CountryOfBirth: audience.CountryOfBirth, AgeGroup: audience.AgeGroup, DailyHoursOnSocialMedia: audience.DailyHoursOnSocialMedia, LastMonthNumberOfPurchases: audience.LastMonthNumberOfPurchases}
	case "insights":
		var insight models.Insight
		if err := db.Where("id = ?", a.RelatedID).Find(&insight).Error; err != nil {
			return nil, fmt.Errorf("something went wrong while fetching insight with id %v of asset with id %v", a.RelatedID, a.ID)
		}
		assetResponse.Insight = &models.InsightWithoutAsset{ID: insight.ID, Text: insight.Text}
	}

	return &assetResponse, nil
}
