package models

type Audience struct {
	Asset
	Gender      string `gorm:"not null;"`
	Country     string `gorm:"not null;"`
	AgeFrom     int    `gorm:"not null;"`
	AgeTo       int    `gorm:"not null;"`
	SocialHours int    `gorm:"not null;"`
	Purchases   int    `gorm:"not null;"`
}
