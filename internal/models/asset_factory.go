package models

import (
	"encoding/json"
	"errors"
)

// AssetCreator defines the interface for creating assets.
type AssetCreator interface {
	Create(data map[string]interface{}) (Asset, error)
}

type ChartCreator struct{}

func (c *ChartCreator) Create(data map[string]interface{}) (Asset, error) {
	var chart Chart
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &chart)
	return &chart, nil
}

type InsightCreator struct{}

func (c *InsightCreator) Create(data map[string]interface{}) (Asset, error) {
	var insight Insight
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &insight)
	return &insight, nil
}

type AudienceCreator struct{}

func (c *AudienceCreator) Create(data map[string]interface{}) (Asset, error) {
	var audience Audience
	bytes, _ := json.Marshal(data)
	json.Unmarshal(bytes, &audience)
	return &audience, nil
}

var AssetFactory = map[string]AssetCreator{
	"chart":    &ChartCreator{},
	"insight":  &InsightCreator{},
	"audience": &AudienceCreator{},
}

func CreateAsset(assetType string, data map[string]interface{}) (Asset, error) {
	creator, ok := AssetFactory[assetType]
	if !ok {
		return nil, errors.New("unknown asset type")
	}
	return creator.Create(data)
}
