package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Chart struct {
	gorm.Model
	Title               string
	HorizontalAxisLabel string
	VerticalAxisLabel   string
	Data                pq.Float64Array `gorm:"type:float8[]"`
	Asset               Asset           `gorm:"polymorphic:Related;"`
}

type ChartWithoutAsset struct {
	ID                  uint
	Title               string
	HorizontalAxisLabel string
	VerticalAxisLabel   string
	Data                pq.Float64Array `gorm:"type:float8[]"`
}
