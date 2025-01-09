package models

// Chart represents a chart with title, axis titles, and data points
// @Description Chart represents a chart with title, axis titles, and data points
type Chart struct {
	Title      string   `json:"title" example:"Sales Chart"`
	AxesTitles []string `json:"axes_titles" example:"['Time', 'Sales']"`
	Data       []int    `json:"data" example:"[10, 20, 30, 40]"`
}
