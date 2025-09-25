package core

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type (
	UserID  string
	AssetID string
)

type AssetType string

const (
	AssetTypeChart    AssetType = "chart"
	AssetTypeInsight  AssetType = "insight"
	AssetTypeAudience AssetType = "audience"
)

var (
	ErrInvalidID        = errors.New("invalid id")
	ErrEmptyTitle       = errors.New("title cannot be empty")
	ErrShortTitle       = errors.New("title must be at least 3 characters")
	ErrEmptyDescription = errors.New("description cannot be empty")
	ErrInvalidTime      = errors.New("timestamps must be non-zero")
	ErrInvalidAxisTitle = errors.New("axis title cannot be empty")
	ErrNoSeries         = errors.New("chart series cannot be empty")
	ErrNoPoints         = errors.New("chart series points cannot be empty")
	ErrInvalidGender    = errors.New("gender must be Male or Female")
	ErrInvalidAgeRange  = errors.New("age range invalid")
)

func nonZero(t time.Time) bool { return !t.IsZero() }

// Base fields shared by all assets (exported so JSON works without DTOs).
type AssetBase struct {
	ID          AssetID   `json:"id"`
	Type        AssetType `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (b *AssetBase) SetDescription(desc string, now time.Time) error {
	desc = strings.TrimSpace(desc)
	if desc == "" {
		return ErrEmptyDescription
	}
	if !nonZero(now) {
		return ErrInvalidTime
	}
	b.Description = desc
	b.UpdatedAt = now.UTC()
	return nil
}

func makeBase(id AssetID, typ AssetType, title, description string, now time.Time) (AssetBase, error) {
	title = strings.TrimSpace(title)
	description = strings.TrimSpace(description)
	if id == "" {
		return AssetBase{}, ErrInvalidID
	}
	if title == "" {
		return AssetBase{}, ErrEmptyTitle
	}
	if len([]rune(title)) < 3 {
		return AssetBase{}, ErrShortTitle
	}
	if description == "" {
		return AssetBase{}, ErrEmptyDescription
	}
	if !nonZero(now) {
		return AssetBase{}, ErrInvalidTime
	}
	now = now.UTC()
	return AssetBase{
		ID:          id,
		Type:        typ,
		Title:       title,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// ----- Chart -----

type ChartPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ChartSeries struct {
	Name   string       `json:"name"`
	Points []ChartPoint `json:"points"`
}

type Chart struct {
	AssetBase
	XAxisTitle string        `json:"xTitle"`
	YAxisTitle string        `json:"yTitle"`
	Series     []ChartSeries `json:"series"`
}

func MakeChart(id AssetID, title, description, xTitle, yTitle string, series []ChartSeries, now time.Time) (*Chart, error) {
	base, err := makeBase(id, AssetTypeChart, title, description, now)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(xTitle) == "" || strings.TrimSpace(yTitle) == "" {
		return nil, ErrInvalidAxisTitle
	}
	if len(series) == 0 {
		return nil, ErrNoSeries
	}
	for i, s := range series {
		if strings.TrimSpace(s.Name) == "" {
			series[i].Name = fmt.Sprintf("Series %d", i+1)
		}
		if len(s.Points) == 0 {
			return nil, ErrNoPoints
		}
	}
	return &Chart{AssetBase: base, XAxisTitle: xTitle, YAxisTitle: yTitle, Series: series}, nil
}

// ----- Insight -----

type Insight struct {
	AssetBase
	Text string `json:"text"`
}

func MakeInsight(id AssetID, title, description, text string, now time.Time) (*Insight, error) {
	base, err := makeBase(id, AssetTypeInsight, title, description, now)
	if err != nil {
		return nil, err
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.New("insight text cannot be empty")
	}
	return &Insight{AssetBase: base, Text: text}, nil
}

// ----- Audience -----

type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
)

type Audience struct {
	AssetBase
	Gender                *Gender `json:"gender,omitempty"`
	BirthCountry          *string `json:"birthCountry,omitempty"`
	AgeMin                *int    `json:"ageMin,omitempty"`
	AgeMax                *int    `json:"ageMax,omitempty"`
	MinHoursSocialPerDay  *int    `json:"minHoursSocialPerDay,omitempty"`
	MinPurchasesLastMonth *int    `json:"minPurchasesLastMonth,omitempty"`
}

func MakeAudience(
	id AssetID,
	title, description string,
	gender *Gender,
	birthCountry *string,
	ageMin, ageMax *int,
	minHoursPerDay, minPurchases *int,
	now time.Time,
) (*Audience, error) {
	base, err := makeBase(id, AssetTypeAudience, title, description, now)
	if err != nil {
		return nil, err
	}
	if gender != nil && *gender != GenderMale && *gender != GenderFemale {
		return nil, ErrInvalidGender
	}
	if ageMin != nil && ageMax != nil && *ageMin > *ageMax {
		return nil, ErrInvalidAgeRange
	}
	return &Audience{
		AssetBase:             base,
		Gender:                gender,
		BirthCountry:          birthCountry,
		AgeMin:                ageMin,
		AgeMax:                ageMax,
		MinHoursSocialPerDay:  minHoursPerDay,
		MinPurchasesLastMonth: minPurchases,
	}, nil
}

// Union for repos/services (concrete assets are returned directly).
type AnyAsset interface{ GetAssetBase() *AssetBase }

// Implement GetAssetBase for union
func (c *Chart) GetAssetBase() *AssetBase    { return &c.AssetBase }
func (i *Insight) GetAssetBase() *AssetBase  { return &i.AssetBase }
func (a *Audience) GetAssetBase() *AssetBase { return &a.AssetBase }