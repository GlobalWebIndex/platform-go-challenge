package serializers

import "github.com/alexandros-georgantas/platform-go-challenge/internal/models"

type Chart struct {
	ID                  int       `json:"id"`
	Title               string    `json:"title"`
	HorizontalAxisLabel string    `json:"xAxisLabel"`
	VerticalAxisLabel   string    `json:"yAxisLabel"`
	Data                []float64 `json:"data"`
}

type ChartsResponse struct {
	Charts     *[]models.Chart `json:"items"`
	TotalCount int             `json:"totalCount"`
}
