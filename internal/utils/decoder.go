package utils

import (
	"encoding/json"
	"fmt"

	"github.com/ceciivanov/platform-go-challenge/internal/models"
)

// DecodeAsset decodes JSON into the correct asset type
func DecodeAsset(data []byte) (models.Asset, error) {

	var base struct {
		Type models.AssetType `json:"type"`
	}

	// Unmarshal the JSON into a struct that only contains the type field to determine the asset type
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	// Unmarshal the JSON into the correct asset type
	switch base.Type {
	case models.ChartType:
		var chart models.Chart
		if err := json.Unmarshal(data, &chart); err != nil {
			return nil, err
		}
		return &chart, nil
	case models.InsightType:
		var insight models.Insight
		if err := json.Unmarshal(data, &insight); err != nil {
			return nil, err
		}
		return &insight, nil
	case models.AudienceType:
		var audience models.Audience
		if err := json.Unmarshal(data, &audience); err != nil {
			return nil, err
		}
		return &audience, nil
	default:
		return nil, fmt.Errorf("invalid asset type")
	}
}
