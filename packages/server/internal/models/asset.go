package models

import (
	"gorm.io/gorm"
)

type Asset struct {
	gorm.Model
	Description string
	RelatedID   uint
	RelatedType string
}

type AssetResponse struct {
	ID          uint
	Description string
	RelatedID   uint
	RelatedType string
	Chart       *ChartWithoutAsset
	Audience    *AudienceWithoutAsset
	Insight     *InsightWithoutAsset
}
