package models

import (
	"encoding/json"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Asset represents an asset in our system.
type Asset struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
	// AssetData contains the asset specific data for Chart/Insight/Audience.
	AssetData json.RawMessage `json:"assetData"`
	Added     time.Time       `json:"added"`
	Modified  time.Time       `json:"modified"`
}

// Chart defines the properties of a chart asset.
type Chart struct {
	Title      string `json:"title"`
	XAxisTitle string `json:"xAxisTitle"`
	YAxisTitle string `json:"yAxisTitle"`
	// Data type is `any` to handle any format of chart data representation.
	Data any `json:"data"`
}

// Insight defines the properties of an insight asset.
type Insight struct {
	Text string `json:"text"`
}

// Audience defines the properties of an audience asset.
type Audience struct {
	Gender       string        `json:"gender" validate:"required,oneof=Male Female"`
	BirthCountry string        `json:"birthCountry" validate:"required"`
	AgeGroup     AgeGroupRange `json:"ageGroup" validate:"required"`
	// Can be converted to float to handle time durations (<hours)
	HoursDailyOnSocialMedia int `json:"hoursDailyOnSocialMedia" validate:"required"`
	Purchases               int `json:"purchases" validate:"required"`
}

// AgeGroupRange defines the range of the age group of an audience.
type AgeGroupRange struct {
	Bottom int `json:"bottom"`
	Top    int `json:"top"`
}

// User defines the properties of a user in the platform.
type User struct {
	ID         string    `json:"id"`
	Password   string    `json:"password"`
	Cookie     string    `json:"-"`
	Privileged bool      `json:"-"`
	Created    time.Time `json:"created"`
}

// Claims defines the user token claim.
type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

// CacheEntry defines a cached response entry.
type CacheEntry struct {
	Token     string
	Response  []byte
	Set       bool
	ExpiresAt time.Time
}
