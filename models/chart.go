package models

type Chart struct {
	Asset
	Title string `gorm:"not null;"`
	XAxes string `gorm:"not null;"`
	YAxes string `gorm:"not null;"`
	Data  string `gorm:"not null;"`
}
