package models

// Audience represents a group of people with certain characteristics
// @Description Audience represents a group of people with certain characteristics
type Audience struct {
	Gender            string `json:"gender" example:"Male"`
	BirthCountry      string `json:"birth_country" example:"USA"`
	AgeGroup          string `json:"age_group" example:"24-35"`
	HoursSpentDaily   int    `json:"hours_spent_daily" example:"3"`
	NumberOfPurchases int    `json:"number_of_purchases" example:"5"`
}
