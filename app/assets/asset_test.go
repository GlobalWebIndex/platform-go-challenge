package assets

import (
	"reflect"
	"testing"
)

func TestGetAsset(t *testing.T) {
	tests := []struct {
		name      string
		assetType string
		wantType  reflect.Type
		wantOk    bool
	}{
		{
			name:      "Get Insight Asset",
			assetType: "INSIGHT",
			wantType:  reflect.TypeOf(&Insight{}),
			wantOk:    true,
		},
		{
			name:      "Get Chart Asset",
			assetType: "CHART",
			wantType:  reflect.TypeOf(&Chart{}),
			wantOk:    true,
		},
		{
			name:      "Get Audience Asset",
			assetType: "AUDIENCE",
			wantType:  reflect.TypeOf(&Audience{}),
			wantOk:    true,
		},
		{
			name:      "Get Invalid Asset",
			assetType: "INVALID",
			wantType:  nil,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := GetAsset(tt.assetType)

			if ok != tt.wantOk {
				t.Errorf("GetAsset() ok = %v, want %v", ok, tt.wantOk)
				return
			}

			if !tt.wantOk {
				if got != nil {
					t.Errorf("GetAsset() got = %v, want nil", got)
				}
				return
			}

			if reflect.TypeOf(got) != tt.wantType {
				t.Errorf("GetAsset() got type = %v, want %v", reflect.TypeOf(got), tt.wantType)
			}
		})
	}
}

func TestAssetInterface(t *testing.T) {
	assets := []Asset{
		&Insight{},
		&Chart{},
		&Audience{},
	}

	for _, asset := range assets {
		if _, ok := asset.(Asset); !ok {
			t.Errorf("%T does not implement Asset interface", asset)
		}
	}
}
