package models

// Asset interface
type Asset interface {
	GetID() string
	SetDescription(desc string)
}

/// ChartData represents one data point in a chart
type ChartData struct {
	DatapointCode string  `json:"datapoint_code"` // e.g. "SM_AGE_18_24"
	Value         float64 `json:"value"`          // numeric value
}

// Chart asset
type Chart struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	XAxisTitle  string      `json:"x_axis_title"`
	YAxisTitle  string      `json:"y_axis_title"`
	Data        []ChartData `json:"data"` // data points
}

func (c *Chart) GetID() string              { return c.ID }
func (c *Chart) SetDescription(desc string) { c.Description = desc }

// Insight asset
type Insight struct {
	ID          string `json:"id"`
	Description string `json:"description"` // short insight text
}

func (i *Insight) GetID() string              { return i.ID }
func (i *Insight) SetDescription(desc string) { i.Description = desc }

// Audience asset
type Audience struct {
	ID          string `json:"id"`
	Gender      string `json:"gender"`       // Male / Female
	Country     string `json:"country"`      // birth country
	AgeGroup    string `json:"age_group"`    // e.g. "24-35"
	SocialHours int    `json:"social_hours"` // hours spent on social media daily
	Purchases   int    `json:"purchases"`    // number of purchases last month
	Description string `json:"description"`  // short description
}

func (a *Audience) GetID() string              { return a.ID }
func (a *Audience) SetDescription(desc string) { a.Description = desc }
