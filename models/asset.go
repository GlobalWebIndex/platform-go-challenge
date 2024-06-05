package models

type Chart struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	AxesTitles  string `json:"axes_titles"`
	Data        string `json:"data"`
	Description string `json:"description"`
}

type Insight struct {
	ID          string `json:"id"`
	Text        string `json:"text"`
	Description string `json:"description"`
}

type Audience struct {
	ID                 string `json:"id"`
	Gender             string `json:"gender"`
	BirthCountry       string `json:"birth_country"`
	AgeGroup           string `json:"age_group"`
	HoursSpentDaily    string `json:"hours_spent_daily"`
	PurchasesLastMonth int    `json:"purchases_last_month"`
	Description        string `json:"description"`
}

type Asset struct {
	Chart    *Chart    `json:"chart,omitempty"`
	Insight  *Insight  `json:"insight,omitempty"`
	Audience *Audience `json:"audience,omitempty"`
}
