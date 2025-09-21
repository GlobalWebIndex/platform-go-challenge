package gwiexercise

import (
	"encoding/json"
	"errors"
	"fmt"
)

type AddAssetRequest struct {
	ID          string          `json:"id"`
	Type        AssetType       `json:"type"`
	Description string          `json:"description"`
	Payload     json.RawMessage `json:"payload"`
}

func (chart *ChartAsset) MarshalJSON() ([]byte, error) {
	type alias ChartAsset
	return json.Marshal((*alias)(chart))
}

func (insight *InsightAsset) MarshalJSON() ([]byte, error) {
	type alias InsightAsset
	return json.Marshal((*alias)(insight))
}

func (audience *AudienceAsset) MarshalJSON() ([]byte, error) {
	type alias AudienceAsset
	return json.Marshal((*alias)(audience))
}

func (req AddAssetRequest) ValidateBasic() error {
	if req.Type == "" {
		return errors.New("type is required")
	}
	if req.Description == "" {
		return errors.New("description is required")
	}
	if len(req.Payload) == 0 {
		return errors.New("payload is required")
	}
	return nil
}

func DecodeAssetFromAddRequest(req AddAssetRequest) (Asset, error) {
	switch req.Type {
	case AssetChart:
		return ExtractChart(req)
	case AssetInsight:
		return ExtractInsight(req)
	case AssetAudience:
		return ExtractAudience(req)
	default:
		return nil, fmt.Errorf("unsupported asset type: %s", req.Type)
	}
}

func ExtractAudience(req AddAssetRequest) (Asset, error) {
	var body struct {
		Gender             string `json:"gender"`
		BirthCountry       string `json:"birthCountry"`
		AgeGroup           string `json:"ageGroup"`
		SocialHoursPerDay  string `json:"socialHoursPerDay"`
		PurchasesLastMonth int    `json:"purchasesLastMonth"`
	}
	if err := json.Unmarshal(req.Payload, &body); err != nil {
		return nil, fmt.Errorf("invalid audience payload: %w", err)
	}
	asset := &AudienceAsset{BaseAsset: BaseAsset{ID: req.ID, Type: AssetAudience, Description: req.Description},
		Gender: body.Gender, BirthCountry: body.BirthCountry, AgeGroup: body.AgeGroup, SocialHoursPerDay: body.SocialHoursPerDay, PurchasesLastMonth: body.PurchasesLastMonth}
	return asset, nil
}

func ExtractInsight(req AddAssetRequest) (Asset, error) {
	var body struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal(req.Payload, &body); err != nil {
		return nil, fmt.Errorf("invalid insight payload: %w", err)
	}
	if body.Text == "" {
		return nil, fmt.Errorf("invalid insight payload: missing text")
	}
	asset := &InsightAsset{BaseAsset: BaseAsset{ID: req.ID, Type: AssetInsight, Description: req.Description}, Text: body.Text}
	return asset, nil
}

func ExtractChart(req AddAssetRequest) (Asset, error) {
	var body struct {
		Title      string    `json:"title"`
		XAxisTitle string    `json:"xAxisTitle"`
		YAxisTitle string    `json:"yAxisTitle"`
		Data       []float64 `json:"data"`
	}
	if err := json.Unmarshal(req.Payload, &body); err != nil {
		return nil, fmt.Errorf("invalid chart payload: %w", err)
	}
	if body.Title == "" {
		return nil, fmt.Errorf("invalid chart payload: missing title")
	}
	asset := &ChartAsset{BaseAsset: BaseAsset{ID: req.ID, Type: AssetChart, Description: req.Description},
		Title: body.Title, XAxisTitle: body.XAxisTitle, YAxisTitle: body.YAxisTitle, Data: body.Data}
	return asset, nil
}

func DecodeAssetFromJSON(data []byte) (Asset, error) {
	var probe struct {
		Type AssetType `json:"type"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return nil, fmt.Errorf("invalid asset json: %w", err)
	}
	var asset Asset
	switch probe.Type {
	case AssetChart:
		asset = &ChartAsset{}
	case AssetInsight:
		asset = &InsightAsset{}
	case AssetAudience:
		asset = &AudienceAsset{}
	default:
		return nil, fmt.Errorf("unsupported asset type: %s", probe.Type)
	}
	if err := json.Unmarshal(data, asset); err != nil {
		return nil, err
	}
	return asset, nil
}
