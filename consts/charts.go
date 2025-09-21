package consts

import (
	"platform-go-challenge/domain"
	"time"
)

var Charts = []domain.Asset{
	domain.Asset{
		Type:        "chart",
		Description: "This is the first chart",
		CreatedAt:   time.Now(),
		Data: domain.AssetData{
			ChartData: &domain.Chart{
				Title:      "First chart",
				XAxisTitle: "Space",
				YAxisTitle: "Time",
				Data: []domain.Point{
					domain.Point{
						X: "0",
						Y: "1",
					},
					domain.Point{
						X: "1",
						Y: "2",
					},
					domain.Point{
						X: "2",
						Y: "0",
					},
					domain.Point{
						X: "3",
						Y: "-1",
					},
					domain.Point{
						X: "4",
						Y: "3",
					},
				},
			},
		},
	},
	domain.Asset{
		Type:        "chart",
		Description: "This is the second chart",
		CreatedAt:   time.Now(),
		Data: domain.AssetData{
			ChartData: &domain.Chart{
				Title:      "Second chart",
				XAxisTitle: "Space",
				YAxisTitle: "Time",
				Data: []domain.Point{
					domain.Point{
						X: "0",
						Y: "1",
					},
					domain.Point{
						X: "1",
						Y: "2",
					},
					domain.Point{
						X: "2",
						Y: "0",
					},
					domain.Point{
						X: "3",
						Y: "-1",
					},
					domain.Point{
						X: "4",
						Y: "3",
					},
				},
			},
		},
	},
	domain.Asset{
		Type:        "chart",
		Description: "This is the third chart",
		CreatedAt:   time.Now(),
		Data: domain.AssetData{
			ChartData: &domain.Chart{
				Title:      "Third chart",
				XAxisTitle: "Space",
				YAxisTitle: "Time",
				Data: []domain.Point{
					domain.Point{
						X: "0",
						Y: "1",
					},
					domain.Point{
						X: "1",
						Y: "2",
					},
					domain.Point{
						X: "2",
						Y: "0",
					},
					domain.Point{
						X: "3",
						Y: "-1",
					},
					domain.Point{
						X: "4",
						Y: "3",
					},
				},
			},
		},
	},
}
