package models

type Chart struct {
	Title      string    `json:"title"`
	XAxisTitle string    `json:"x_axis_title"`
	YAxisTitle string    `json:"y_axis_title"`
	DataPoints []float64 `json:"data"`
}
