package models

import (
	"gorm.io/gorm"
)

type GenderType string
type AgeType string

const (
	MALE   GenderType = "MALE"
	FEMALE GenderType = "FEMALE"
)

const (
	CHILDREN AgeType = "0-14"
	YOUTH    AgeType = "15-24"
	ADULTS   AgeType = "25-64"
	SENIORS  AgeType = "65-94"
)

type Audience struct {
	gorm.Model
	Gender                     GenderType
	CountryOfBirth             string
	AgeGroup                   AgeType
	DailyHoursOnSocialMedia    int
	LastMonthNumberOfPurchases int
	Asset                      Asset `gorm:"polymorphic:Related;"`
}

type AudienceWithoutAsset struct {
	ID                         uint
	Gender                     GenderType
	CountryOfBirth             string
	AgeGroup                   AgeType
	DailyHoursOnSocialMedia    int
	LastMonthNumberOfPurchases int
}
