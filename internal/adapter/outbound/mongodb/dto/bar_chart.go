package dto

import "github.com/Kercyn/crud_template/internal/core/domain"

type BarChart struct {
	AssetBase
	Title      string         `bson:"title,omitempty"`
	YAxisTitle string         `bson:"y_axis_title,omitempty"`
	Data       []BarChartData `bson:"data,omitempty"`
}

type BarChartData struct {
	Label string  `bson:"label,omitempty"`
	Value float64 `bson:"value,omitempty"`
}

func (b BarChart) ToDomain() (domain.Asset, error) {
	var data []domain.BarData

	for _, d := range b.Data {
		data = append(data, domain.BarData{
			Label: d.Label,
			Value: d.Value,
		})
	}

	domainAsset, err := b.AssetBase.ToDomain()
	if err != nil {
		return domain.BarChart{}, err
	}

	return domain.BarChart{
		Asset:      domainAsset,
		Title:      b.Title,
		YAxisTitle: b.YAxisTitle,
		Data:       data,
	}, nil
}
