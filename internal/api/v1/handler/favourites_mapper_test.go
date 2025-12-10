package handler

import (
	"testing"
	"time"

	"github.com/GlobalWebIndex/platform-go-challenge/internal/api/v1/dto"
	"github.com/GlobalWebIndex/platform-go-challenge/internal/domain/favourites"
)

func TestParseAssetType_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected favourites.AssetType
	}{
		{"chart", string(dto.AssetTypeChart), favourites.AssetTypeChart},
		{"insight", string(dto.AssetTypeInsight), favourites.AssetTypeInsight},
		{"audience", string(dto.AssetTypeAudience), favourites.AssetTypeAudience},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := parseAssetType(tt.input)
			if err != nil {
				t.Fatalf("parseAssetType(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Fatalf("parseAssetType(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseAssetType_Invalid(t *testing.T) {
	t.Parallel()

	if _, err := parseAssetType("unknown"); err == nil {
		t.Fatalf("expected error for unknown asset type, got nil")
	}
}

func TestDTOToDomainAsset_SuccessVariants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		request  dto.CreateFavouriteRequest
		wantType favourites.AssetType
		check    func(t *testing.T, asset favourites.Asset)
	}{
		{
			name:     "chart",
			wantType: favourites.AssetTypeChart,
			request: dto.CreateFavouriteRequest{
				Asset: dto.CreateAssetRequest{
					Type:        dto.AssetTypeChart,
					Description: "chart-desc",
					Chart: &dto.ChartPayload{
						Title: "t",
						XAxis: "x",
						YAxis: "y",
						Data:  []float64{1, 2},
					},
				},
			},
			check: func(t *testing.T, asset favourites.Asset) {
				if asset.Chart == nil {
					t.Fatalf("expected Chart to be non-nil")
				}
				if asset.Chart.Title != "t" || asset.Chart.XAxis != "x" || asset.Chart.YAxis != "y" {
					t.Fatalf("unexpected Chart mapping: %+v", asset.Chart)
				}
			},
		},
		{
			name:     "insight",
			wantType: favourites.AssetTypeInsight,
			request: dto.CreateFavouriteRequest{
				Asset: dto.CreateAssetRequest{
					Type:        dto.AssetTypeInsight,
					Description: "insight-desc",
					Insight: &dto.InsightPayload{
						Text: "hello",
					},
				},
			},
			check: func(t *testing.T, asset favourites.Asset) {
				if asset.Insight == nil {
					t.Fatalf("expected Insight to be non-nil")
				}
				if asset.Insight.Text != "hello" {
					t.Fatalf("unexpected Insight.Text: %q", asset.Insight.Text)
				}
			},
		},
		{
			name:     "audience",
			wantType: favourites.AssetTypeAudience,
			request: dto.CreateFavouriteRequest{
				Asset: dto.CreateAssetRequest{
					Type:        dto.AssetTypeAudience,
					Description: "audience-desc",
					Audience: &dto.AudiencePayload{
						Gender:                   "Male",
						BirthCountry:             "GR",
						AgeGroup:                 "25-34",
						HoursOnSocialMediaPerDay: 3,
						PurchasesLastMonth:       2,
					},
				},
			},
			check: func(t *testing.T, asset favourites.Asset) {
				if asset.Audience == nil {
					t.Fatalf("expected Audience to be non-nil")
				}
				if asset.Audience.Gender != "Male" ||
					asset.Audience.BirthCountry != "GR" ||
					asset.Audience.AgeGroup != "25-34" {
					t.Fatalf("unexpected Audience mapping: %+v", asset.Audience)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := dtoToDomainAsset(tt.request)
			if err != nil {
				t.Fatalf("dtoToDomainAsset unexpected error: %v", err)
			}
			if got.Type != tt.wantType {
				t.Fatalf("expected Type %q, got %q", tt.wantType, got.Type)
			}
			if got.Description != tt.request.Asset.Description {
				t.Fatalf("expected Description %q, got %q", tt.request.Asset.Description, got.Description)
			}
			tt.check(t, got)
		})
	}
}

func TestDTOToDomainAsset_MissingPayload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		request dto.CreateFavouriteRequest
	}{
		{
			name: "chart_without_payload",
			request: dto.CreateFavouriteRequest{
				Asset: dto.CreateAssetRequest{
					Type:        dto.AssetTypeChart,
					Description: "chart-desc",
					Chart:       nil,
				},
			},
		},
		{
			name: "insight_without_payload",
			request: dto.CreateFavouriteRequest{
				Asset: dto.CreateAssetRequest{
					Type:        dto.AssetTypeInsight,
					Description: "insight-desc",
					Insight:     nil,
				},
			},
		},
		{
			name: "audience_without_payload",
			request: dto.CreateFavouriteRequest{
				Asset: dto.CreateAssetRequest{
					Type:        dto.AssetTypeAudience,
					Description: "audience-desc",
					Audience:    nil,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := dtoToDomainAsset(tt.request)
			if err == nil {
				t.Fatalf("expected error for missing payload, got nil")
			}
			if !errorsIsInvalidAssetPayload(err) {
				t.Fatalf("expected ErrInvalidAssetPayload, got %v", err)
			}
		})
	}
}

// errorsIsInvalidAssetPayload is a tiny helper to avoid importing `errors`
// just to call errors.Is in multiple places.
func errorsIsInvalidAssetPayload(err error) bool {
	return err != nil && err == favourites.ErrInvalidAssetPayload
}

func TestFavouriteEntryToDTO_MapsFields(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	entry := favourites.FavouriteEntry{
		Favourite: favourites.Favourite{
			UserID:    "user-1",
			AssetID:   "asset-1",
			CreatedAt: now,
		},
		Asset: favourites.Asset{
			ID:          "asset-1",
			Type:        favourites.AssetTypeInsight,
			Description: "desc",
			Insight: &favourites.Insight{
				Text: "hello",
			},
		},
	}

	dtoEntry := favouriteEntryToDTO(entry)

	if dtoEntry.Asset.ID != "asset-1" {
		t.Fatalf("expected Asset.ID %q, got %q", "asset-1", dtoEntry.Asset.ID)
	}
	if dtoEntry.Asset.Type != dto.AssetType(favourites.AssetTypeInsight) {
		t.Fatalf("expected Asset.Type %q, got %q", dto.AssetTypeInsight, dtoEntry.Asset.Type)
	}
	if dtoEntry.CreatedAt != now {
		t.Fatalf("expected CreatedAt %v, got %v", now, dtoEntry.CreatedAt)
	}
	if dtoEntry.Asset.Insight == nil || dtoEntry.Asset.Insight.Text != "hello" {
		t.Fatalf("unexpected Insight mapping: %+v", dtoEntry.Asset.Insight)
	}
}

func TestAssetToDTO_MapsAllPayloads(t *testing.T) {
	t.Parallel()

	asset := favourites.Asset{
		ID:          "asset-1",
		Type:        favourites.AssetTypeAudience,
		Description: "desc",
		Chart: &favourites.Chart{
			Title: "t",
			XAxis: "x",
			YAxis: "y",
			Data:  []float64{1, 2},
		},
		Insight: &favourites.Insight{
			Text: "hello",
		},
		Audience: &favourites.Audience{
			Gender:                   "Female",
			BirthCountry:             "US",
			AgeGroup:                 "18-24",
			HoursOnSocialMediaPerDay: 4,
			PurchasesLastMonth:       5,
		},
	}

	dtoAsset := assetToDTO(asset)

	if dtoAsset.ID != "asset-1" {
		t.Fatalf("expected ID %q, got %q", "asset-1", dtoAsset.ID)
	}
	if dtoAsset.Type != dto.AssetType(favourites.AssetTypeAudience) {
		t.Fatalf("expected Type %q, got %q", dto.AssetTypeAudience, dtoAsset.Type)
	}
	if dtoAsset.Description != "desc" {
		t.Fatalf("expected Description %q, got %q", "desc", dtoAsset.Description)
	}

	if dtoAsset.Chart == nil || dtoAsset.Chart.Title != "t" {
		t.Fatalf("unexpected Chart mapping: %+v", dtoAsset.Chart)
	}
	if dtoAsset.Insight == nil || dtoAsset.Insight.Text != "hello" {
		t.Fatalf("unexpected Insight mapping: %+v", dtoAsset.Insight)
	}
	if dtoAsset.Audience == nil ||
		dtoAsset.Audience.Gender != "Female" ||
		dtoAsset.Audience.BirthCountry != "US" {
		t.Fatalf("unexpected Audience mapping: %+v", dtoAsset.Audience)
	}
}
