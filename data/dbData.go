package data

type Users struct {
	UserId int    `json:"UserId"`
	Name   string `json:"name"`
}

var User = []Users{
	{UserId: 1, Name: "John Doe"},
	{UserId: 2, Name: "Jane Doe"},
	{UserId: 3, Name: "Alan Doe"},
	{UserId: 4, Name: "Bob Doe"},
	{UserId: 5, Name: "Alice Doe"},
}

// All assets

type AllAssets struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

var Assets = []AllAssets{
	{
		ID:          1,
		Type:        "Insight",
		Description: "Insights for youngs",
	},
	{
		ID:          2,
		Type:        "Chart",
		Description: "Annual inflation rate",
	},
	{
		ID:          3,
		Type:        "Audience",
		Description: "Number of hours young people spend on their phones",
	},
	{
		ID:          4,
		Type:        "Chart",
		Description: "Growth",
	},
	{
		ID:          5,
		Type:        "Insight",
		Description: "Random Info of the day",
	},
}

// Favorite asset table

type FavAssets struct {
	ID      int `json:"id"`
	UserId  int `json:"UserId"`
	AssetId int `json:"AssetId"`
}

var FavAsset = []FavAssets{
	{
		ID:      1,
		UserId:  1,
		AssetId: 2,
	},
	{
		ID:      2,
		UserId:  1,
		AssetId: 1,
	},
	{
		ID:      3,
		UserId:  2,
		AssetId: 3,
	},
	{
		ID:      4,
		UserId:  3,
		AssetId: 5,
	},
}

// Chart asset details

type ChartAsset struct {
	ChartId int       `json:"chartId"`
	AssetId int       `json:"assetId"`
	Title   string    `json:"title"`
	XAxis   string    `json:"x_axis"`
	YAxis   string    `json:"y_axis"`
	Values  []float64 `json:"values"`
}

var ChartAssets = []ChartAsset{
	{
		ChartId: 1,
		AssetId: 2,
		Title:   "Monthly Sales",
		XAxis:   "Month",
		YAxis:   "Revenue",
		Values:  []float64{100, 200, 400},
	},
	{
		ChartId: 2,
		AssetId: 4,
		Title:   "Annual inflation rate",
		XAxis:   "Month",
		YAxis:   "percentage",
		Values:  []float64{1, 2, 1.5},
	},
}

// Insight Asset details

type InsightAsset struct {
	InsightId int    `json:"insightId"`
	AssetId   int    `json:"asset_id"`
	Title     string `json:"title"`
}

var InsightAssets = []InsightAsset{
	{
		InsightId: 1,
		AssetId:   1,
		Title:     "Insights for youngs",
	},
	{
		InsightId: 2,
		AssetId:   5,
		Title:     "Random Information of each day",
	},
}

// Audience asset details

type AudienceAsset struct {
	AudienceId     int    `json:"audienceId"`
	AssetId        int    `json:"asset_id"`
	Gender         string `json:"gender"`
	AgeGroup       string `json:"age_group"`
	BirthCountry   string `json:"country"`
	HoursOnSocial  int    `json:"hours_social_media"`
	PurchasesLastM int    `json:"purchases_last_month"`
}

var AudienceAssets = []AudienceAsset{
	{
		AudienceId:     1,
		AssetId:        3,
		Gender:         "Male",
		AgeGroup:       "24-35",
		BirthCountry:   "USA",
		HoursOnSocial:  3,
		PurchasesLastM: 5,
	},
}
