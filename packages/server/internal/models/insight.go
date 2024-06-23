package models

import "gorm.io/gorm"

type Insight struct {
	gorm.Model
	Text  string
	Asset Asset `gorm:"polymorphic:Related;"`
}

type InsightWithoutAsset struct {
	ID   uint
	Text string
}
