package serializers

import "github.com/alexandros-georgantas/platform-go-challenge/internal/models"

type Insight struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type InsightResponse struct {
	Insights   *[]models.Insight `json:"items"`
	TotalCount int               `json:"totalCount"`
}
