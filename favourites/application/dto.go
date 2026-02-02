package application

import "github.com/manos/favourites/favourites/domain"

type AssetDTO struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	OwnerUserID string               `json:"ownerUserId"`
	Type        domain.FavouriteType `json:"type"`

	// Insight
	Text string `json:"text,omitempty"`

	// Chart
	XAxisTitle string    `json:"xAxisTitle,omitempty"`
	YAxisTitle string    `json:"yAxisTitle,omitempty"`
	Data       []float64 `json:"data,omitempty"`

	// Audience (minimal)
	SampleSize        int64   `json:"sampleSize,omitempty"`
	TotalRespondents  int64   `json:"totalRespondents,omitempty"`
	EstimatedReach    int64   `json:"estimatedReach,omitempty"`
	PopulationPercent float64 `json:"populationPercent,omitempty"`
}

type FavouriteDTO struct {
	ID          string               `json:"id"`
	UserID      string               `json:"userId"`
	Type        domain.FavouriteType `json:"type"`
	Description string               `json:"description"`
	Asset       AssetDTO             `json:"asset"`
}

type AddFavouriteRequest struct {
	Type        string `json:"type"`
	AssetID     string `json:"assetId"`
	Description string `json:"description"`
}

type FavouritePageDTO struct {
	Items      []FavouriteDTO `json:"items"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	Total      int            `json:"total"`
	TotalPages int            `json:"totalPages"`
}
