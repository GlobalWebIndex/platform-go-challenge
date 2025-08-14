package domain

import (
	"encoding/json"
	"time"
)

type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

// Asset interface defines common behavior for all asset types
type Asset interface {
	GetID() string
	GetType() AssetType
	GetDescription() string
	SetDescription(string)
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetUpdatedAt(time.Time)
	Validate() error
}

// BaseAsset contains common fields for all assets
type BaseAsset struct {
	ID          string    `json:"id"`
	Type        AssetType `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (b *BaseAsset) GetID() string          { return b.ID }
func (b *BaseAsset) GetType() AssetType     { return b.Type }
func (b *BaseAsset) GetDescription() string { return b.Description }
func (b *BaseAsset) SetDescription(desc string) {
	b.Description = desc
	b.UpdatedAt = time.Now()
}
func (b *BaseAsset) GetCreatedAt() time.Time  { return b.CreatedAt }
func (b *BaseAsset) GetUpdatedAt() time.Time  { return b.UpdatedAt }
func (b *BaseAsset) SetUpdatedAt(t time.Time) { b.UpdatedAt = t }

// Chart represents a chart asset
type Chart struct {
	BaseAsset
	Title      string           `json:"title"`
	XAxisTitle string           `json:"x_axis_title"`
	YAxisTitle string           `json:"y_axis_title"`
	Data       []ChartDataPoint `json:"data"`
}

type ChartDataPoint struct {
	X interface{} `json:"x"`
	Y interface{} `json:"y"`
}

func (c *Chart) Validate() error {
	if c.ID == "" {
		return ErrMissingRequiredField
	}
	if c.Title == "" {
		return ErrMissingRequiredField
	}
	return nil
}

// Insight represents an insight asset
type Insight struct {
	BaseAsset
	Content  string   `json:"content"`
	Tags     []string `json:"tags,omitempty"`
	Category string   `json:"category,omitempty"`
}

func (i *Insight) Validate() error {
	if i.ID == "" {
		return ErrMissingRequiredField
	}
	if i.Content == "" {
		return ErrMissingRequiredField
	}
	return nil
}

// Audience represents an audience asset
type Audience struct {
	BaseAsset
	Gender             []string `json:"gender,omitempty"`
	BirthCountries     []string `json:"birth_countries,omitempty"`
	AgeGroups          []string `json:"age_groups,omitempty"`
	SocialMediaHours   string   `json:"social_media_hours,omitempty"`
	PurchasesLastMonth int      `json:"purchases_last_month,omitempty"`
}

func (a *Audience) Validate() error {
	if a.ID == "" {
		return ErrMissingRequiredField
	}
	return nil
}

// AssetFromJSON creates assets from JSON
func AssetFromJSON(data []byte) (Asset, error) {
	var base struct {
		Type AssetType `json:"type"`
	}

	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	switch base.Type {
	case AssetTypeChart:
		var chart Chart
		if err := json.Unmarshal(data, &chart); err != nil {
			return nil, err
		}
		return &chart, nil
	case AssetTypeInsight:
		var insight Insight
		if err := json.Unmarshal(data, &insight); err != nil {
			return nil, err
		}
		return &insight, nil
	case AssetTypeAudience:
		var audience Audience
		if err := json.Unmarshal(data, &audience); err != nil {
			return nil, err
		}
		return &audience, nil
	default:
		return nil, ErrInvalidAssetType
	}
}

// NewChart creates a new chart asset
func NewChart(id, title, xAxis, yAxis, description string, data []ChartDataPoint) *Chart {
	now := time.Now()
	return &Chart{
		BaseAsset: BaseAsset{
			ID:          id,
			Type:        AssetTypeChart,
			Description: description,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Title:      title,
		XAxisTitle: xAxis,
		YAxisTitle: yAxis,
		Data:       data,
	}
}

// NewInsight creates a new insight asset
func NewInsight(id, content, description string, tags []string, category string) *Insight {
	now := time.Now()
	return &Insight{
		BaseAsset: BaseAsset{
			ID:          id,
			Type:        AssetTypeInsight,
			Description: description,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Content:  content,
		Tags:     tags,
		Category: category,
	}
}

// NewAudience creates a new audience asset
func NewAudience(id, description string) *Audience {
	now := time.Now()
	return &Audience{
		BaseAsset: BaseAsset{
			ID:          id,
			Type:        AssetTypeAudience,
			Description: description,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}
}
