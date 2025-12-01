package models

import (
	"encoding/json"
	"errors"
	"time"

	"platform-go-challenge/internal/common"
)

// ---------- Asset interface ----------
type Asset interface {
	GetID() string
	GetDescription() string
	GetCreatedTimestamp() time.Time
	GetUpdatedTimestamp() time.Time
	SetUpdatedTimestamp(time.Time)
}

type AssetConstructor func() Asset

var assetRegistry = map[string]AssetConstructor{
	"chart":    func() Asset { return &Chart{} },
	"insight":  func() Asset { return &Insight{} },
	"audience": func() Asset { return &Audience{} },
}

// ---------- Asset types ----------

// Base Asset to be included in all asset types
type BaseAsset struct {
	// Asset ID is represented as a string to be more user-friendly; typically a user is asked to give a unique-per-user name to their favorite and I assume this is what is happening here.
	// Alternative representations could be UUID or int64.
	ID          string    `json:"id" merge:"ignore"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at" merge:"ignore"`
	UpdatedAt   time.Time `json:"updated_at" merge:"ignore"`
}

func (b *BaseAsset) GetID() string {
	return b.ID
}

func (b *BaseAsset) GetDescription() string {
	return b.Description
}

func (b *BaseAsset) SetDescription(desc string) {
	b.Description = desc
}

func (b *BaseAsset) GetCreatedTimestamp() time.Time {
	return b.CreatedAt
}

func (b *BaseAsset) GetUpdatedTimestamp() time.Time {
	return b.UpdatedAt
}

func (b *BaseAsset) SetUpdatedTimestamp(t time.Time) {
	b.UpdatedAt = t
}

func (b *BaseAsset) setCreatedTimestamp(t time.Time) {
	if b.CreatedAt.IsZero() {
		b.CreatedAt = t
	}
}

// ChartDataPoint represents a single data point in a Chart
type ChartDataPoint struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

// Chart asset
type Chart struct {
	BaseAsset
	Title      string           `json:"title"`
	XAxisTitle string           `json:"x_axis_title"`
	YAxisTitle string           `json:"y_axis_title"`
	Data       []ChartDataPoint `json:"data"`
}

// Insight asset
type Insight struct {
	BaseAsset
	InsightText string `json:"insight_text"`
}

// Audience asset
type Audience struct {
	BaseAsset
	Gender           string `json:"gender"`
	BirthCountry     string `json:"birth_country"`
	AgeGroup         string `json:"age_group"`
	DailySMHours     int    `json:"daily_sm_hours"`
	MonthlyPurchases int    `json:"monthly_purchases"`
}

func CreateAsset(assetData json.RawMessage) (Asset, error) {
	// Unmarshal to get the type
	var typeWrapper struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(assetData, &typeWrapper); err != nil {
		return nil, err
	}

	constructor, exists := assetRegistry[typeWrapper.Type]
	if !exists {
		return nil, errors.New(common.ErrInvalidAssetType)
	}

	// Unmarshal into concrete type
	asset := constructor()
	if err := json.Unmarshal(assetData, asset); err != nil {
		return nil, err
	}

	// Initialize timestamps
	now := time.Now()
	if b, ok := asset.(interface{ setCreatedTimestamp(time.Time) }); ok {
		b.setCreatedTimestamp(now)
	}
	asset.SetUpdatedTimestamp(now)

	return asset, nil
}
