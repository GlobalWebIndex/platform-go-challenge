package domain

import "time"

type Gender int8

const (
	GenderMale Gender = iota
	GenderFemale
)

type AgeGroup int8

const (
	AgeGroup0to18 AgeGroup = iota
	AgeGroup19to24
	AgeGroup25to29
	AgeGroup30to34
	AgeGroup35to39
	AgeGroup40to44
	AgeGroup45to49
	AgeGroup50to54
	AgeGroup55to59
	AgeGroup60to64
	AgeGroup65to69
	AgeGroup70plus
)

type Audience struct {
	Asset
	Gender                    Gender
	AgeGroup                  AgeGroup
	NumberOfPurchasesPerMonth int
	HoursSpendDailyOnSocial   time.Duration
}

func (a Audience) GetType() AssetType {
	return a.Asset.GetType()
}

func (a Audience) GetDescription() string {
	return a.Asset.GetDescription()
}

func (a Audience) GetIsFavourite() bool {
	return a.Asset.GetIsFavourite()
}
