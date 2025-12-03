package models

type Chart struct {
	Title      string      `json:"title"`
	XAxisTitle string      `json:"x_axis_title"`
	YAxisTitle string      `json:"y_axis_title"`
	Data       []DataPoint `json:"data"`
}

type DataPoint struct {
	X any     `json:"x"`
	Y float64 `json:"y"`
}
