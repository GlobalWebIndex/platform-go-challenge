package model

import (
	"time"

	"github.com/google/uuid"
)

type Chart struct {
	Id         uuid.UUID
	Title      string
	XAxisTitle string
	YAxisTitle string
	DataPoints []DataPoint
}

type DataPoint struct {
	X float64
	Y float64
}

type Insight struct {
	Id          uuid.UUID
	Description string
}

type Audience struct {
	Id             uuid.UUID
	Gender         rune
	BirthCountry   string
	AgeFrom        int
	AgeTo          int
	HoursSpentFrom int
	HoursSpentTo   int
	PurchasesFrom  int
	PurchasesTo    int
}

type User struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	Gender    rune
	BirthDate time.Time
}

// This is sent when starring an item
type UserFavouriteShort struct {
	Description  string
	ResourceType string
	Id           uuid.UUID
}

// This is returned when listing a user's favourites
type UserFavourite struct {
	Description string
	Chart       *Chart    `json:",omitempty"`
	Insight     *Insight  `json:",omitempty"`
	Audience    *Audience `json:",omitempty"`
}

type UserFavourites struct {
	Favourites []UserFavourite
}
