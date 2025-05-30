package postgres

import (
	"github.com/uptrace/bun"
)

// Asset represents a general asset entity in the database.
// It includes relations to specialized asset types such as Chart, Insight, and Audience.
type Asset struct {
	bun.BaseModel `bun:"table:assets,alias:a"` // Specifies the table name and alias for the Asset model.

	ID   int    `bun:"id,autoincrement,pk"` // Primary key for the Asset table.
	Type string `bun:"type,notnull"`        // Type of the asset (e.g., chart, insight, audience).

	// Relations to specialized asset tables
	Chart    *Chart    `bun:"rel:has-one,join:id=asset_id"` // One-to-one relation with the Chart table.
	Insight  *Insight  `bun:"rel:has-one,join:id=asset_id"` // One-to-one relation with the Insight table.
	Audience *Audience `bun:"rel:has-one,join:id=asset_id"` // One-to-one relation with the Audience table.

	Favorites []*Favorite `bun:"rel:has-many,join:id=asset_id"` // One-to-many relation with the Favorite table.
}

// Chart represents a chart asset with metadata and data.
// It is linked to the Asset table via a one-to-one relationship.
type Chart struct {
	bun.BaseModel `bun:"table:charts,alias:c"` // Specifies the table name and alias for the Chart model.

	AssetID    int    `bun:"asset_id,pk"`           // Primary key and foreign key linking to the Asset table.
	Title      string `bun:"title,notnull"`         // Title of the chart.
	TitleXAxis string `bun:"title_x_axis,nullzero"` // Label for the X-axis of the chart.
	TitleYAxis string `bun:"title_y_axis,nullzero"` // Label for the Y-axis of the chart.
	Data       any    `bun:"data,notnull"`          // Chart data stored in a flexible format (e.g., JSON).

	Asset *Asset `bun:"rel:belongs-to,join:asset_id=id"` // Belongs-to relation linking back to the Asset table.
}

// Audience represents an audience asset with demographic and behavioral data.
// It is linked to the Asset table via a one-to-one relationship.
type Audience struct {
	bun.BaseModel `bun:"table:audiences,alias:au"` // Specifies the table name and alias for the Audience model.

	AssetID                      int    `bun:"asset_id,pk"`                                // Primary key and foreign key linking to the Asset table.
	Gender                       string `bun:"gender,nullzero"`                            // Gender of the audience.
	BirthCountry                 string `bun:"birth_country,nullzero"`                     // Birth country of the audience.
	AgeGroup                     string `bun:"age_group,nullzero"`                         // Age group of the audience.
	HoursSpentDailyOnSocialMedia int    `bun:"hours_spent_daily_on_social_media,nullzero"` // Daily social media usage in hours.
	LastMonthPurchases           int    `bun:"last_month_purchases,nullzero"`              // Number of purchases made last month.

	Asset *Asset `bun:"rel:belongs-to,join:asset_id=id"` // Belongs-to relation linking back to the Asset table.
}

// Insight represents an insight asset containing textual content.
// It is linked to the Asset table via a one-to-one relationship.
type Insight struct {
	bun.BaseModel `bun:"table:insights,alias:i"` // Specifies the table name and alias for the Insight model.

	AssetID int    `bun:"asset_id,pk"`     // Primary key and foreign key linking to the Asset table.
	Content string `bun:"content,notnull"` // Textual content of the insight.

	Asset *Asset `bun:"rel:belongs-to,join:asset_id=id"` // Belongs-to relation linking back to the Asset table.
}
