package repositories

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username           string     `gorm:"not null;"`
	Password           string     `gorm:"not null;"`
	FavouriteCharts    []Chart    `gorm:"many2many:users_charts;"`
	FavouriteInsights  []Insight  `gorm:"many2many:users_insights;"`
	FavouriteAudiences []Audience `gorm:"many2many:users_audiences;"`
}
