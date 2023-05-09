package domain

type Audience struct {
	Gender             string
	BirthCountry       string
	AgeGroup           string
	HoursSpentDaily    uint
	PurchasesLastMonth uint
	Description        string
	AssetType          string
}

func (a Audience) GetDescription() string {
	return a.Description
}

func (a Audience) GetType() string {
	return string(AudienceType)
}
