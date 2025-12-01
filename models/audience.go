package models

type Audience struct {
	ID int `json:"id"`
	Gender string `json:"gender"`
	BirthCountry string `json:"birth_country"`
	AgeGroup string `json:"age_group"`
	HoursSocialMediaDaily int `json:"hours_social_media_daily"`
	NumberOfPurchasesLastMonth int `json:"number_of_purchases_last_month"`
	Description string `json:"description"`
}

func (a *Audience) GetID() int {
	return a.ID
}

func (a *Audience) GetType() AssetType {
	return AssetTypeAudience
}

func (a *Audience) GetDescription() string {
	return a.Description
}

func (a *Audience) SetDescription(desc string) {
	a.Description = desc
}