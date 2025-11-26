package unit

import (
	"platform-go-challenge/models"
	"testing"
)

func TestAssetType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		assetType models.AssetType
		want     bool
	}{
		{"Valid Audience", models.AssetTypeAudience, true},
		{"Valid Chart", models.AssetTypeChart, true},
		{"Valid Insight", models.AssetTypeInsight, true},
		{"Invalid empty", models.AssetType(""), false},
		{"Invalid random", models.AssetType("Random"), false},
		{"Invalid lowercase", models.AssetType("audience"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.assetType.IsValid(); got != tt.want {
				t.Errorf("AssetType.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssetType_String(t *testing.T) {
	tests := []struct {
		name     string
		assetType models.AssetType
		want     string
	}{
		{"Audience", models.AssetTypeAudience, "Audience"},
		{"Chart", models.AssetTypeChart, "Chart"},
		{"Insight", models.AssetTypeInsight, "Insight"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.assetType.String(); got != tt.want {
				t.Errorf("AssetType.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAssetType_Value(t *testing.T) {
	tests := []struct {
		name      string
		assetType models.AssetType
		wantValue string
		wantErr   bool
	}{
		{"Valid Audience", models.AssetTypeAudience, "Audience", false},
		{"Valid Chart", models.AssetTypeChart, "Chart", false},
		{"Valid Insight", models.AssetTypeInsight, "Insight", false},
		{"Invalid type", models.AssetType("Invalid"), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.assetType.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("AssetType.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantValue {
				t.Errorf("AssetType.Value() = %v, want %v", got, tt.wantValue)
			}
		})
	}
}

func TestAssetType_Scan(t *testing.T) {
	tests := []struct {
		name    string
		value   interface{}
		want    models.AssetType
		wantErr bool
	}{
		{"Valid string Audience", "Audience", models.AssetTypeAudience, false},
		{"Valid string Chart", "Chart", models.AssetTypeChart, false},
		{"Valid string Insight", "Insight", models.AssetTypeInsight, false},
		{"Valid bytes Audience", []byte("Audience"), models.AssetTypeAudience, false},
		{"Invalid string", "Invalid", "", true},
		{"Nil value", nil, "", true},
		{"Invalid type", 123, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var at models.AssetType
			err := at.Scan(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("AssetType.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && at != tt.want {
				t.Errorf("AssetType.Scan() = %v, want %v", at, tt.want)
			}
		})
	}
}
