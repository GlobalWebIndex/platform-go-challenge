package main

import "github.com/google/uuid"

// AssetType represents the type of asset
const (
	ChartType    = "chart"
	InsightType  = "insight"
	AudienceType = "audience"
)

type Asset interface {
	GetID() uuid.UUID
	GetType() string
	GetDescription() string
	SetDescription(desc string)
	IsFavorite() bool
	SetFavorite(isFav bool)
}

type Chart struct {
	ID          uuid.UUID
	Title       string
	XAxisTitle  string
	YAxisTitle  string
	Data        []float64
	Description string
	Favorite    bool
}

func (c *Chart) GetID() uuid.UUID           { return c.ID }
func (c *Chart) GetType() string            { return ChartType }
func (c *Chart) GetDescription() string     { return c.Description }
func (c *Chart) SetDescription(desc string) { c.Description = desc }
func (c *Chart) IsFavorite() bool           { return c.Favorite }
func (c *Chart) SetFavorite(isFav bool)     { c.Favorite = isFav }

type Insight struct {
	ID          uuid.UUID
	Text        string
	Description string
	Favorite    bool
}

func (i *Insight) GetID() uuid.UUID           { return i.ID }
func (i *Insight) GetType() string            { return InsightType }
func (i *Insight) GetDescription() string     { return i.Description }
func (i *Insight) SetDescription(desc string) { i.Description = desc }
func (i *Insight) IsFavorite() bool           { return i.Favorite }
func (i *Insight) SetFavorite(isFav bool)     { i.Favorite = isFav }

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

type Audience struct {
	ID           uuid.UUID
	Gender       Gender
	BirthCountry string
	AgeGroup     string
	SocialHours  int
	Purchases    int
	Description  string
	Favorite     bool
}

func (a *Audience) GetID() uuid.UUID           { return a.ID }
func (a *Audience) GetType() string            { return AudienceType }
func (a *Audience) GetDescription() string     { return a.Description }
func (a *Audience) SetDescription(desc string) { a.Description = desc }
func (a *Audience) IsFavorite() bool           { return a.Favorite }
func (a *Audience) SetFavorite(isFav bool)     { a.Favorite = isFav }

type User struct {
	ID         uuid.UUID
	Favourites []Asset
}
