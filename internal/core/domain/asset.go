package domain

import (
	"errors"
	"fmt"
)

type AssetType int16

func NewAssetTypeFromString(s string) (AssetType, error) {
	switch s {
	case "audience":
		return AssetTypeAudience, nil
	case "bar_chart":
		return AssetTypeBarChart, nil
	case "insight":
		return AssetTypeInsight, nil
	default:
		return AssetTypeInvalid, errors.New(fmt.Sprintf("Invalid asset type: %s", s))
	}
}

const (
	AssetTypeInvalid            = -1
	AssetTypeAudience AssetType = iota
	AssetTypeBarChart
	AssetTypeInsight
)

type Asset interface {
	GetType() AssetType
	GetDescription() string
	GetIsFavourite() bool
}

type AssetBase struct {
	Type        AssetType
	Description string
	IsFavourite bool
}

func (a AssetBase) GetType() AssetType {
	return a.Type
}

func (a AssetBase) GetDescription() string {
	return a.Description
}

func (a AssetBase) GetIsFavourite() bool {
	return a.IsFavourite
}
