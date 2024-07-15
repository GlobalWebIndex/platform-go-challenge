package models

// Point represents a data point in the chart
type Point struct {
	X float32 `json:"X"`
	Y float32 `json:"Y"`
}

// Chart represents a chart asset
type Chart struct {
	ID          int       `json:"id"`
	Type        AssetType `json:"type"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
	XAxesTitle  string    `json:"xAxesTitle"`
	YAxesTitle  string    `json:"yAxesTitle"`
	DataPoints  []Point   `json:"dataPoints"`
}

// Implement the Asset interface for Chart
func (c Chart) GetID() int {
	return c.ID
}

func (c Chart) GetType() AssetType {
	return c.Type
}

func (c Chart) GetDescription() string {
	return c.Description
}
