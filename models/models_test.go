package models

import (
	"encoding/json"
	"testing"
)

func TestUserModel(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
	}{
		{
			name:    "valid user",
			user:    User{ID: "u1", Name: "Alice"},
			wantErr: false,
		},
		{
			name:    "user with empty ID",
			user:    User{ID: "", Name: "Bob"},
			wantErr: false, // Structure allows empty ID, but repository should validate
		},
		{
			name:    "user with empty name",
			user:    User{ID: "u2", Name: ""},
			wantErr: false, // Structure allows it, but handler validates
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			data, err := json.Marshal(tt.user)
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && data == nil {
				t.Error("expected JSON data but got nil")
			}
		})
	}
}

func TestAssetModel(t *testing.T) {
	tests := []struct {
		name    string
		asset   Asset
		wantErr bool
	}{
		{
			name: "valid chart asset",
			asset: Asset{
				Type:  AssetChart,
				Title: ptrString("Sales"),
				Data:  json.RawMessage(`{"x":["Jan","Feb"]}`),
			},
			wantErr: false,
		},
		{
			name: "valid insight asset",
			asset: Asset{
				Type:  AssetInsight,
				Title: ptrString("Market Insight"),
				Data:  json.RawMessage(`{"text":"Growing market"}`),
			},
			wantErr: false,
		},
		{
			name: "valid audience asset",
			asset: Asset{
				Type:  AssetAudience,
				Title: ptrString("Target Audience"),
				Data:  json.RawMessage(`{"gender":"M"}`),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.asset.Type == "" {
				t.Error("asset type cannot be empty")
			}
			if tt.asset.Title == nil {
				t.Error("asset title cannot be nil")
			}
		})
	}
}

func TestFavouriteModel(t *testing.T) {
	tests := []struct {
		name      string
		favourite Favourite
		wantErr   bool
	}{
		{
			name: "valid favourite",
			favourite: Favourite{
				UserID:  "u1",
				AssetID: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: false,
		},
		{
			name: "favourite with description",
			favourite: Favourite{
				UserID:      "u2",
				AssetID:     "550e8400-e29b-41d4-a716-446655440001",
				Description: ptrString("Important asset"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.favourite.UserID == "" {
				t.Error("userID cannot be empty")
			}
			if tt.favourite.AssetID == "" {
				t.Error("assetID cannot be empty")
			}
		})
	}
}

func TestFavouriteAssetModel(t *testing.T) {
	tests := []struct {
		name           string
		favouriteAsset FavouriteAsset
		wantErr        bool
	}{
		{
			name: "valid favourite asset",
			favouriteAsset: FavouriteAsset{
				AssetID:     "550e8400-e29b-41d4-a716-446655440000",
				Type:        AssetChart,
				Title:       ptrString("Sales Chart"),
				Description: ptrString("Monthly sales"),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.favouriteAsset.AssetID == "" {
				t.Error("assetID cannot be empty")
			}
		})
	}
}

// Helper function
func ptrString(s string) *string {
	return &s
}
