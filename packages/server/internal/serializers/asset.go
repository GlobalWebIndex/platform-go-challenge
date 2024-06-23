package serializers

import "github.com/alexandros-georgantas/platform-go-challenge/internal/models"

type Asset struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	RelatedType string `json:"relatedType"`
	RelatedID   int    `json:"relatedId"`
	Chart       *Chart
	Audience    *Audience
	Insight     *Insight
}

type AddToFavorites struct {
	ID int `json:"id"`
}

type UpdateDescription struct {
	Description string `json:"description"`
}

type AssetsResponse struct {
	Assets     *[]models.AssetResponse `json:"items"`
	TotalCount int                     `json:"totalCount"`
}

type FavoritesResponse struct {
	Favorites  *[]models.FavoritesResponse `json:"items"`
	TotalCount int                         `json:"totalCount"`
}
