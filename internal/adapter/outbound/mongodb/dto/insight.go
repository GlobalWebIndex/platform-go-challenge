package dto

import (
	"github.com/Kercyn/crud_template/internal/core/domain"
)

type Insight struct {
	AssetBase
	Data string `bson:"data,omitempty"`
}

func (i Insight) ToDomain() (domain.Asset, error) {
	domainAsset, err := i.AssetBase.ToDomain()
	if err != nil {
		return domain.BarChart{}, err
	}

	return domain.Insight{
		Asset: domainAsset,
		Data:  i.Data,
	}, nil
}
