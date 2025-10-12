package model

// AssetType defines the type of asset
type AssetType string

// Supported asset types
// Use them in the request body
const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

// BaseAsset contains common fields for all assets
type BaseAsset struct {
	ID          string    `json:"id"`
	Type        AssetType `json:"type"`
	Description string    `json:"description"`
	CreatedAt   int64     `json:"created_at"`
	UpdatedAt   int64     `json:"updated_at"`
}

// Chart represents a chart asset
type Chart struct {
	BaseAsset `json:",inline"`
	Title     string        `json:"title"`
	XAxisKey  string        `json:"x_axis_key"`
	YAxisKey  string        `json:"y_axis_key"`
	Data      []interface{} `json:"data"`
}

// Insight represents an insight asset
type Insight struct {
	BaseAsset `json:",inline"`
	Title     string `json:"title"`
	Content   string `json:"content"`
}

// Audience represents an audience segment
type Audience struct {
	BaseAsset `json:",inline"`
	Gender    string   `json:"gender"` // "Male", "Female", or ""
	Countries []string `json:"countries"`
	AgeMin    int      `json:"age_min"`
	AgeMax    int      `json:"age_max"`
	MinHours  int      `json:"min_hours_daily_social_media"`
	MinBuys   int      `json:"min_purchases_last_month"`
}

// Asset is a union type that holds any asset
type Asset struct {
	Chart    *Chart    `json:"chart,omitempty"`
	Insight  *Insight  `json:"insight,omitempty"`
	Audience *Audience `json:"audience,omitempty"`
}

// GetBaseAsset returns the base asset from any asset type
func (a *Asset) GetBaseAsset() *BaseAsset {
	switch {
	case a.Chart != nil:
		return &a.Chart.BaseAsset
	case a.Insight != nil:
		return &a.Insight.BaseAsset
	case a.Audience != nil:
		return &a.Audience.BaseAsset
	default:
		return nil
	}
}

// Favorite represents a favorited asset for a user
type Favorite struct {
	Asset   Asset `json:"asset"`
	AddedAt int64 `json:"added_at"`
}

// User contains a user's favorite assets
type User struct {
	ID        string
	Favorites map[string]*Favorite // key is asset ID
}