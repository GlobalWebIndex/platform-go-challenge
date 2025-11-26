package models

type Insight struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Text string `json:"text"`
}
