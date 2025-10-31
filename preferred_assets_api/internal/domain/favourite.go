package domain

import (
	"errors"
	"time"
)

type Favourite struct {
	UserID    string
	AssetID   string
	CreatedAt time.Time
	AssetType AssetType

	Audience *Audience
	Chart    *Chart
	Insight  *Insight
}

func (f *Favourite) GetAsset() Asset {
	switch f.AssetType {
	case AssetTypeAudience:
		return f.Audience
	case AssetTypeChart:
		return f.Chart
	case AssetTypeInsight:
		return f.Insight
	default:
		return nil
	}
}

func (f *Favourite) SetAsset(asset Asset) error {
	f.AssetID = asset.GetID()
	f.AssetType = asset.GetType()

	switch a := asset.(type) {
	case *Audience:
		f.Audience = a
		f.Chart = nil
		f.Insight = nil
	case *Chart:
		f.Audience = nil
		f.Chart = a
		f.Insight = nil
	case *Insight:
		f.Audience = nil
		f.Chart = nil
		f.Insight = a
	default:
		return errors.New("unknown asset type")
	}
	return nil
}
