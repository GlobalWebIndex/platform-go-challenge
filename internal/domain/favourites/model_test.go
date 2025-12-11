package favourites

import "testing"

func TestAssetValidate_ValidVariants(t *testing.T) {
	chartAsset := Asset{
		Type: AssetTypeChart,
		Chart: &Chart{
			Title: "Example",
			XAxis: "Age",
			YAxis: "Score",
			Data:  []float64{1.0, 2.0},
		},
	}

	insightAsset := Asset{
		Type:    AssetTypeInsight,
		Insight: &Insight{Text: "40% of millennials..."},
	}

	audienceAsset := Asset{
		Type: AssetTypeAudience,
		Audience: &Audience{
			Gender:                   "Male",
			BirthCountry:             "UK",
			AgeGroup:                 "24-35",
			HoursOnSocialMediaPerDay: 3,
			PurchasesLastMonth:       2,
		},
	}

	tests := []struct {
		name  string
		asset Asset
	}{
		{"chart", chartAsset},
		{"insight", insightAsset},
		{"audience", audienceAsset},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.asset.Validate(); err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestAssetValidate_InvalidPayloadForType_Missing(t *testing.T) {
	tests := []struct {
		name  string
		asset Asset
	}{
		{
			name:  "chart without payload",
			asset: Asset{Type: AssetTypeChart},
		},
		{
			name:  "insight without payload",
			asset: Asset{Type: AssetTypeInsight},
		},
		{
			name:  "audience without payload",
			asset: Asset{Type: AssetTypeAudience},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.asset.Validate()
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if err != ErrInvalidAssetPayload {
				t.Fatalf("expected ErrInvalidAssetPayload, got %v", err)
			}
		})
	}
}

func TestAssetValidate_InvalidPayloadForType_ExtraFields(t *testing.T) {
	// Chart type but with extra non-nil Insight/Audience should fail.
	asset := Asset{
		Type: AssetTypeChart,
		Chart: &Chart{
			Title: "Chart",
		},
		Insight: &Insight{Text: "should not be set"},
	}

	if err := asset.Validate(); err == nil {
		t.Fatalf("expected error, got nil")
	} else if err != ErrInvalidAssetPayload {
		t.Fatalf("expected ErrInvalidAssetPayload, got %v", err)
	}
}

func TestAssetValidate_InvalidType(t *testing.T) {
	a := Asset{
		Type: "unknown",
	}

	err := a.Validate()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if err != ErrInvalidAssetType {
		t.Fatalf("expected ErrInvalidAssetType, got %v", err)
	}
}

func TestIsValidAssetType(t *testing.T) {
	tests := []struct {
		name string
		t    AssetType
		want bool
	}{
		{"chart", AssetTypeChart, true},
		{"insight", AssetTypeInsight, true},
		{"audience", AssetTypeAudience, true},
		{"unknown", AssetType("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidAssetType(tt.t); got != tt.want {
				t.Fatalf("IsValidAssetType(%q) = %v, want %v", tt.t, got, tt.want)
			}
		})
	}
}
