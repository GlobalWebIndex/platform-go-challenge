package models

// Example: "Males from 24-35 that spent more than 3 hours on social media daily"
// Audience is a collection of filters & criteria.
// Each characteristic filters the audience by a specific attribute.
type Audience struct {
	Characteristics []Characteristic `json:"characteristics"`
}


type Characteristic struct {
	Type  CharacteristicType `json:"type"`.
	Value any                `json:"value"` // Can be string, Gender, AgeGroup etc.
}

type CharacteristicType string

const (
	CharacteristicGender             CharacteristicType = "gender"
	CharacteristicBirthCountry       CharacteristicType = "birth_country"
	CharacteristicAgeGroup           CharacteristicType = "age_group"
	CharacteristicSocialMediaHours   CharacteristicType = "social_media_hours"
	CharacteristicPurchasesLastMonth CharacteristicType = "purchases_last_month"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type AgeGroup struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type SocialMediaHours struct {
	Min float64 `json:"min"`
	Max float64 `json:"max,omitempty"`
}

//TODO: merge with AgeGroup.
type PurchasesLastMonth struct {
	Min int `json:"min"`
	Max int `json:"max,omitempty"`
}


