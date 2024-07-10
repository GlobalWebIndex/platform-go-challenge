package models

type AssetType int

type Chart struct {
	Title string
	XAxis string
	YAxis string
	Data  map[string]float64
}

type Insight struct {
	Text string
}

type Audience struct {
	Gender             string
	BirthCountry       string
	AgeGroup           string
	HoursOnSocialMedia int
	PurchasesLastMonth int
}

type Asset struct {
	ID          string
	Type        AssetType
	Description string
	Chart       *Chart
	Insight     *Insight
	Audience    *Audience
}

type User struct {
	ID        uint
	Favorites []Asset
}
