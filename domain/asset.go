package domain

// AssetType represents the type of asset as a string.
type AssetType string

const (
	// AssetTypeChart represents a chart asset type.
	AssetTypeChart string = "chart"
	// AssetTypeInsight represents an insight asset type.
	AssetTypeInsight string = "insight"
	// AssetTypeAudience represents an audience asset type.
	AssetTypeAudience string = "audience"
)

// Asset represents a generic asset in the domain layer.
// It contains information about the asset's ID, description, type, and optional related entities.
type Asset struct {
	ID          int       `json:"id"`                 // Unique identifier for the asset.
	Description string    `json:"description"`        // Description of the asset.
	Type        string    `json:"type"`               // Type of the asset (e.g., chart, insight, audience).
	Chart       *Chart    `json:"chart,omitempty"`    // Optional chart associated with the asset.
	Insight     *Insight  `json:"insight,omitempty"`  // Optional insight associated with the asset.
	Audience    *Audience `json:"audience,omitempty"` // Optional audience associated with the asset.
}

// Chart represents a data visualization chart.
// It contains information about the chart's title, axis titles, and data.
type Chart struct {
	Title      string // Title of the chart.
	TitleXAxis string // Title of the X-axis.
	TitleYAxis string // Title of the Y-axis.
	Data       any    // Data used to populate the chart (can be of any type).
}

// Insight represents a piece of analytical information or interpretation.
// It contains content that provides meaningful information derived from data analysis.
type Insight struct {
	Content string // Textual content of the insight, containing the analytical information.
}

// Audience represents demographic and behavioral information about a group of people.
type Audience struct {
	Gender                       string // Gender of the audience (e.g., male, female, non-binary).
	BirthCountry                 string // Country where the audience was born.
	AgeGroup                     string // Age group of the audience (e.g., 18-24, 25-34).
	HoursSpentDailyOnSocialMedia int    // Average number of hours the audience spends daily on social media.
	LastMonthPurchases           int    // Number of purchases made by the audience in the last month.
}
