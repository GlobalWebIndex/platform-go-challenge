package models

type Insight struct {
	Asset
	Text string `gorm:"not null;"`
}
