package models

type Chart struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	Title      string `json:"title"`
	XAxisTitle string `json:"xaxistitle"`
	YAxisTitle string `json:"yaxistitle"`
}
