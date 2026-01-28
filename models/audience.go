package models

type Audience struct {
	Gender                string   `json:"gender"`
	BirthCountry          string   `json:"birth_country"`
	AgeGroups             []string `json:"age_groups"`
	HoursOnSocialMediaMin int      `json:"hours_on_social_media_min"`
	PurchasesLastMonthMin int      `json:"purchases_last_month_min"`
}
