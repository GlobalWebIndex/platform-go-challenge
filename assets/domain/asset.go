// Package domain defines the asset catalog model.
//
// The catalog holds three kinds of items: insights (narrative text), charts (a single series
// with axis labels), and audiences (reach and share). All items are owner-scoped and share
// a common BaseAsset (id, name, summary, owner, timestamps).
package domain

import "time"

const (
	InsightType  = "insight"
	ChartType    = "chart"
	AudienceType = "audience"
)

// BaseAsset is the shared identity and metadata for every catalog asset.
type BaseAsset struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Summary   string    `json:"summary,omitempty"`
	OwnerID   string    `json:"owner_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Insight is a narrative insight: BaseAsset plus body text.
type Insight struct {
	BaseAsset
	Body string `json:"body"`
}

// Chart is a single-series chart: BaseAsset plus axis labels and the data series.
type Chart struct {
	BaseAsset
	AxisX  string    `json:"axis_x"`
	AxisY  string    `json:"axis_y"`
	Series []float64 `json:"series"`
}

// Audience is a segment: BaseAsset plus reach metrics and optional demographic/behavior criteria.
type Audience struct {
	BaseAsset
	Reach           int      `json:"reach"`
	SharePct        float64  `json:"share_pct"`
	SampleSize      int      `json:"sample_size,omitempty"`
	Gender          []string `json:"gender,omitempty"`
	BirthCountry    []string `json:"birthCountry,omitempty"`
	AgeGroups       []string `json:"ageGroups,omitempty"`
	HoursSocialMedia []string `json:"hoursSocialMedia,omitempty"`
	Purchases       []string `json:"purchases,omitempty"`
}
