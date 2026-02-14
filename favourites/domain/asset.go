package domain

import "time"

const (
	AssetTypeChart    = "chart"
	AssetTypeInsight  = "insight"
	AssetTypeAudience = "audience"
)

type Asset struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Description   string        `json:"description"`
	SourceAssetID string        `json:"source_asset_id,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Chart         *ChartData    `json:"chart,omitempty"`
	Insight       *InsightData  `json:"insight,omitempty"`
	Audience      *AudienceData `json:"audience,omitempty"`
}

type ChartData struct {
	Title      string      `json:"title"`
	XAxisTitle string      `json:"x_axis_title"`
	YAxisTitle string      `json:"y_axis_title"`
	Data       interface{} `json:"data"`
}

type InsightData struct {
	Text string `json:"text"`
}

type AudienceData struct {
	Gender                string  `json:"gender"`
	BirthCountry          string  `json:"birth_country"`
	AgeGroups             string  `json:"age_groups"`
	HoursSocialMediaDaily float64 `json:"hours_social_media_daily"`
	PurchasesLastMonth    int     `json:"purchases_last_month"`
	HoursSocialMedia      string  `json:"hours_social_media,omitempty"`
	Purchases             string  `json:"purchases,omitempty"`
	SampleSize            int     `json:"sample_size,omitempty"`
	TotalRespondents      int     `json:"total_respondents,omitempty"`
	EstimatedReach        int     `json:"estimated_reach,omitempty"`
	PopulationPercent     float64 `json:"population_percent,omitempty"`
}
