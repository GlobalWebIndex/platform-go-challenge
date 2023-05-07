package repositories

import "gorm.io/gorm"

type Chart struct {
	gorm.Model
	Title       string `gorm:"not null;"`
	XTitle      string `gorm:"not null;"`
	YTitle      string `gorm:"not null;"`
	Data        string `gorm:"not null;"`
	Description string `gorm:"not null;"`
	Users       []User `gorm:"many2many:users_charts;"`
}

type Insight struct {
	gorm.Model
	Text        string `gorm:"not null"`
	Description string `gorm:"not null;"`
	Users       []User `gorm:"many2many:users_insights;"`
}

type Audience struct {
	gorm.Model
	Gender             string `gorm:"not null"`
	BirthCountry       string `gorm:"not null"`
	AgeGroup           string `gorm:"not null"`
	HoursSpentOnSocial uint   `gorm:"not null"`
	PurchasesLastMonth uint   `gorm:"not null"`
	Description        string `gorm:"not null;"`
	Users              []User `gorm:"many2many:users_audiences;"`
}
