package mapping_test

import (
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/mapping"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

func TestAssetReqToDomain_Audience(t *testing.T) {
	// Arrange
	gender := "female"
	birthCountry := "Canada"
	ageGroup := "25-34"
	hours := 4.5
	purchases := 12
	req := dto.AssetRequest{
		ID:              "a1",
		Type:            "audience",
		Title:           "Audience Title",
		Description:     "Description",
		Gender:          &gender,
		BirthCountry:    &birthCountry,
		AgeGroup:        &ageGroup,
		HoursSocial:     &hours,
		PurchasesLastMo: &purchases,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Act
	asset, err := mapping.AssetReqToDomain(req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	aud, ok := asset.(*domain.Audience)
	if !ok {
		t.Fatalf("expected Audience, got %T", asset)
	}
	if aud.GetID() != "a1" || aud.Gender != "female" || aud.BirthCountry != "Canada" || aud.AgeGroup != "25-34" || aud.HoursSocial != 4.5 || aud.PurchasesLastMo != 12 {
		t.Errorf("audience fields not mapped correctly: %+v", aud)
	}
}

func TestAssetReqToDomain_Chart(t *testing.T) {
	// Arrange
	req := dto.AssetRequest{
		ID:          "c1",
		Type:        "chart",
		Title:       "Chart Title",
		Description: "Chart Description",
		AxesTitles:  []string{"X", "Y"},
		Data:        [][]float64{{1, 2}, {3, 4}},
	}

	// Act
	asset, err := mapping.AssetReqToDomain(req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	chart, ok := asset.(*domain.Chart)
	if !ok {
		t.Fatalf("expected Chart, got %T", asset)
	}
	if chart.GetID() != "c1" || len(chart.AxesTitles) != 2 || len(chart.Data) != 2 {
		t.Errorf("chart fields not mapped correctly: %+v", chart)
	}
}

func TestAssetReqToDomain_Insight(t *testing.T) {
	// Arrange
	req := dto.AssetRequest{
		ID:    "i1",
		Type:  "insight",
		Title: "Insight Title",
	}

	// Act
	asset, err := mapping.AssetReqToDomain(req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	insight, ok := asset.(*domain.Insight)
	if !ok {
		t.Fatalf("expected Insight, got %T", asset)
	}
	if insight.GetID() != "i1" || insight.GetTitle() != "Insight Title" {
		t.Errorf("insight fields not mapped correctly: %+v", insight)
	}
}

func TestAssetReqToDomain_UnsupportedType(t *testing.T) {
	// Arrange
	req := dto.AssetRequest{
		ID:    "x1",
		Type:  "unknown",
		Title: "Unknown Asset",
	}

	// Act
	asset, err := mapping.AssetReqToDomain(req)

	// Assert
	if err == nil {
		t.Fatal("expected error for unsupported type, got nil")
	}
	if asset != nil {
		t.Errorf("expected nil asset for unsupported type, got %+v", asset)
	}
}

func TestAssetReqToDomain_NilPointers(t *testing.T) {
	// Arrange: Audience with all optional fields nil
	req := dto.AssetRequest{
		ID:        "a2",
		Type:      "audience",
		Title:     "Nil Audience",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Act
	asset, err := mapping.AssetReqToDomain(req)

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	aud, ok := asset.(*domain.Audience)
	if !ok {
		t.Fatalf("expected Audience, got %T", asset)
	}
	if aud.Gender != "" || aud.BirthCountry != "" || aud.AgeGroup != "" || aud.HoursSocial != 0 || aud.PurchasesLastMo != 0 {
		t.Errorf("expected zero values for nil pointers, got %+v", aud)
	}
}

func TestAssetReqToDomain_MultipleAssets(t *testing.T) {
	// Arrange
	gender := "male"
	hours := 5.5
	reqs := []dto.AssetRequest{
		{
			ID:          "a1",
			Type:        "audience",
			Title:       "Audience 1",
			Gender:      &gender,
			HoursSocial: &hours,
		},
		{
			ID:         "c1",
			Type:       "chart",
			Title:      "Chart 1",
			AxesTitles: []string{"X", "Y"},
			Data:       [][]float64{{1, 2}},
		},
		{
			ID:    "i1",
			Type:  "insight",
			Title: "Insight 1",
		},
	}

	var assets []domain.Asset
	var err error

	// Act
	for _, req := range reqs {
		var a domain.Asset
		a, err = mapping.AssetReqToDomain(req)
		if err != nil {
			t.Fatalf("unexpected error mapping asset %s: %v", req.ID, err)
		}
		assets = append(assets, a)
	}

	// Assert
	if len(assets) != 3 {
		t.Fatalf("expected 3 assets, got %d", len(assets))
	}

	if aud, ok := assets[0].(*domain.Audience); !ok || aud.GetID() != "a1" || aud.Gender != "male" {
		t.Errorf("Audience asset not mapped correctly: %+v", assets[0])
	}
	if chart, ok := assets[1].(*domain.Chart); !ok || chart.GetID() != "c1" || len(chart.AxesTitles) != 2 {
		t.Errorf("Chart asset not mapped correctly: %+v", assets[1])
	}
	if insight, ok := assets[2].(*domain.Insight); !ok || insight.GetID() != "i1" {
		t.Errorf("Insight asset not mapped correctly: %+v", assets[2])
	}
}

func TestAssetReqToDomain_NilSlice(t *testing.T) {
	// Arrange
	var reqs []dto.AssetRequest // nil slice
	var assets []domain.Asset
	var err error

	// Act
	for _, req := range reqs {
		var a domain.Asset
		a, err = mapping.AssetReqToDomain(req)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		assets = append(assets, a)
	}

	// Assert
	if len(assets) != 0 {
		t.Errorf("expected empty asset slice, got %d elements", len(assets))
	}
}
