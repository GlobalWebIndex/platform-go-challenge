package domain

import "time"

type Asset struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UserId      string    `json:"userId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type InsightAsset struct {
	Asset
	Text string `json:"text"`
}

type AudienceAsset struct {
	Asset
	SampleSize        int64   `json:"sampleSize"`
	TotalRespondents  int64   `json:"totalRespondents"`
	EstimatedReach    int64   `json:"estimatedReach"`
	PopulationPercent float64 `json:"populationPercent"`
	// Waves intentionally skipped
}

type ChartAsset struct {
	Asset
	XAxisTitle string    `json:"xAxisTitle"`
	YAxisTitle string    `json:"yAxisTitle"`
	Data       []float64 `json:"data"`
}
