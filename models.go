package main

import "time"

// Data structures. All data is in memory. No connection to persistent storage on this version of the API.
type Chart struct {
	Title string `json:"title"`
	Xaxis string `json:"x_axis"`
	Yaxis string `json:"y_axis"`
	Data []float64 `json:"data"`
}

type Insight struct {
	Text string `json:"text"`
}

type Audience struct {
	Gender string `json:"gender"`
	BirthCountry string `json:"birth_country"`
	AgeGroups []string `json:"age_groups"`
	DailyHoursOnSocialMedia float64 `json:"daily_hours_on_social_media"`
	PurchasesLastMonth int `json:"purchases_last_month"`
}

// Each Asset objects wraps either a Chart, an Insight, or an Audience payload in the Data field.
type Asset struct {
	ID string `json:"id"`
	Description string `json:"description"`
	Type string `json:"type"`
	Data interface{} `json:"data"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}