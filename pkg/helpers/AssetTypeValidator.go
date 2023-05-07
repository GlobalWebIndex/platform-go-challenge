package helpers

import (
	"errors"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
)

type AssetTypeValidatorInterface interface {
	ValidateAssetType(string) error
}

type AssetTypeValidator struct{}

func NewAssetTypeValidator() *AssetTypeValidator {
	return &AssetTypeValidator{}
}

func (validator AssetTypeValidator) ValidateAssetType(assetType string) error {
	if assetType != string(domain.ChartType) &&
		assetType != string(domain.AudienceType) &&
		assetType != string(domain.InsightType) {
		return errors.New("unknown asset type: " + string(assetType))
	}

	return nil
}
