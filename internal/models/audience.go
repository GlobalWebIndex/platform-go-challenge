package models

// Audience represents an audience asset
type Audience struct {
	ID                int       `json:"id"`
	Type              AssetType `json:"type"`
	Description       string    `json:"description"`
	Age               uint      `json:"age"`
	AgeGroup          string    `json:"ageGroup"`
	Gender            string    `json:"gender"`
	BirthCountry      string    `json:"birthCountry"`
	HoursSpentOnMedia uint      `json:"hoursSpentOnMedia"`
	NumberOfPurchases uint      `json:"numberOfPurchases"`
}

// Implement the Asset interface for Audience
func (a Audience) GetID() int {
	return a.ID
}

func (a Audience) GetType() AssetType {
	return a.Type
}

func (a Audience) GetDescription() string {
	return a.Description
}
