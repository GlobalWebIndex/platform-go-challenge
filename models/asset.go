package models

type Asset struct {
	ID        uint `gorm:"primary_key;"`
	UserId    int  `gorm:"not null;"`
	Favourite bool `gorm:"not null;default:true"`
}
