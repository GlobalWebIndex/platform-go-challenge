package models

import (
	"time"
)

type Asset struct {
	ID          string                 `json:"id"`
	TypeID      string                 `json:"type"`
	Description string                 `json:"description"`
	CreatedAt   time.Time              `json:"created_at"`
	Data        map[string]interface{} `json:"data"`
}

type Chart struct {
	Asset
	Title string `json:"title"`
	XAxis string `json:"x_axis"`
	YAxis string `json:"y_axis"`
	Data  []int  `json:"data"`
}

type Insight struct {
	Asset
	Text string `json:"text"`
}

type Audience struct {
	Asset
	Gender         string `json:"gender"`
	BirthCountry   string `json:"birth_country"`
	AgeGroup       string `json:"age_group"`
	SocialHours    string `json:"social_hours"`
	PurchasesLast4 string `json:"purchases_last"`
}
