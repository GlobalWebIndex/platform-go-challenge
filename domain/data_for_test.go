package domain

var CorrectInputTestAssetData = []interface{}{
	&Insight{
		Text:        "40% of millenials spend more than 3hours on social media daily",
		Description: "bla bla",
	},
	&Chart{
		Description: "bla bla",
		Title:       "Something about GDP and Tax",
		XTitle:      "GDP",
		YTitle:      "Tax",
		Data: XYData{
			X: []float64{1, 2, 3, 4, 5},
			Y: []float64{1, 2, 3, 4, 5},
		},
	},
	&Audience{
		AgeMax:            29,
		AgeMin:            20,
		Gender:            FemaleGenderType,
		Country:           "Sweden",
		HoursSpent:        3,
		NumberOfPurchases: 3,
		Description:       "bla bla",
	},
}

var WrongInputTestAssetData = []interface{}{
	&Insight{
		Text:        "",
		Description: "bla bla",
	},
	&Insight{
		Text:        "40% of millenials spend more than 3hours on social media daily",
		Description: "",
	},
	&Insight{
		Text:        "",
		Description: "",
	},
	&Chart{
		Description: "",
		Title:       "Relationship between tax and GDP",
		XTitle:      "GDP",
		YTitle:      "Tax",
		Data: XYData{
			X: []float64{1, 2, 3, 4, 5},
			Y: []float64{1, 2, 3, 4, 5},
		},
	},
	&Chart{
		Description: "bla bla",
		Title:       "",
		XTitle:      "GDP",
		YTitle:      "Tax",
		Data: XYData{
			X: []float64{1, 2, 3, 4, 5},
			Y: []float64{1, 2, 3, 4, 5},
		},
	},
	&Chart{
		Description: "bla bla",
		Title:       "Relationship between tax and GDP",
		XTitle:      "",
		YTitle:      "Tax",
		Data: XYData{
			X: []float64{1, 2, 3, 4, 5},
			Y: []float64{1, 2, 3, 4, 5},
		},
	},
	&Chart{
		Description: "bla bla",
		Title:       "Relationship between tax and GDP",
		XTitle:      "GDP",
		YTitle:      "",
		Data: XYData{
			X: []float64{1, 2, 3, 4, 5},
			Y: []float64{1, 2, 3, 4, 5},
		},
	},
	&Chart{
		Description: "bla bla",
		Title:       "Relationship between tax and GDP",
		XTitle:      "GDP",
		YTitle:      "Tax",
		Data: XYData{
			X: []float64{},
			Y: []float64{},
		},
	},
	&Chart{
		Description: "bla bla",
		Title:       "Relationship between tax and GDP",
		XTitle:      "GDP",
		YTitle:      "Tax",
		Data: XYData{
			X: []float64{1, 2, 3, 4},
			Y: []float64{1, 2, 3, 4, 5},
		},
	},
	&Audience{
		AgeMax:            200,
		AgeMin:            0,
		Gender:            FemaleGenderType,
		Country:           "Sweden",
		HoursSpent:        3,
		NumberOfPurchases: 3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Gender:            GenderType("lalasdla"),
		Country:           "Sweden",
		HoursSpent:        3,
		NumberOfPurchases: 3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Country:           "Sweden",
		HoursSpent:        3,
		NumberOfPurchases: 3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Gender:            FemaleGenderType,
		Country:           "Mordor",
		HoursSpent:        3,
		NumberOfPurchases: 3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Gender:            FemaleGenderType,
		Country:           "Sweden",
		HoursSpent:        -3,
		NumberOfPurchases: 3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Gender:            FemaleGenderType,
		Country:           "Sweden",
		HoursSpent:        3,
		NumberOfPurchases: -3,
		Description:       "bla bla",
	},
	&Audience{
		AgeMax:            30,
		AgeMin:            20,
		Gender:            FemaleGenderType,
		Country:           "Sweden",
		HoursSpent:        3,
		NumberOfPurchases: 3,
	},
}

var WrongInputTestQueryData = []QueryAssets{
	{
		Limit:  0,
		LastID: 10,
		Type:   AudienceAssetType,
	},
	{
		Limit:  10,
		LastID: 1,
	},
}
