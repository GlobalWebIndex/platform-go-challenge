package models

type Audience struct {
	ID            uint   `json:"id" gorm:"primaryKey"`
	Gender        string `json:"gender"`
	BirthCountry  string `json:"birthcountry"`
	AgeGroup      string `json:"agegroup"`
	DailyHours    int    `json:"dailyhours"`
	NoOfPurchases int    `json:"noofpurchases"`
}
