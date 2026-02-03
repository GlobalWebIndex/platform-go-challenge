package favorites

import (
	"testing"
)

func TestBaseAsset_ValidateCommon(t *testing.T) {
	tests := []struct {
		name         string
		base         BaseAsset
		expectedType AssetType
		wantErr      bool
		errMsg       string
	}{
		{
			name: "valid base asset",
			base: BaseAsset{
				ID:   "123",
				Name: "Test Asset",
				Type: AssetTypeChart,
			},
			expectedType: AssetTypeChart,
			wantErr:      false,
		},
		{
			name: "missing id",
			base: BaseAsset{
				Name: "Test Asset",
				Type: AssetTypeChart,
			},
			expectedType: AssetTypeChart,
			wantErr:      true,
			errMsg:       "validation failed: id is required",
		},
		{
			name: "missing name",
			base: BaseAsset{
				ID:   "123",
				Type: AssetTypeChart,
			},
			expectedType: AssetTypeChart,
			wantErr:      true,
			errMsg:       "validation failed: name is required",
		},
		{
			name: "type mismatch",
			base: BaseAsset{
				ID:   "123",
				Name: "Test Asset",
				Type: AssetTypeChart,
			},
			expectedType: AssetTypeInsight,
			wantErr:      true,
			errMsg:       "validation failed: invalid asset type: expected insight, got chart",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.base.ValidateCommon(tt.expectedType)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: BaseAsset.ValidateCommon() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("%s: BaseAsset.ValidateCommon() error = %v, wantErrMsg %v", tt.name, err, tt.errMsg)
			}
		})
	}
}

func TestBaseAsset_Getters(t *testing.T) {
	base := BaseAsset{
		ID:   "123",
		Name: "Test",
		Type: AssetTypeChart,
	}

	if base.GetID() != "123" {
		t.Errorf("GetID() = %v, want %v", base.GetID(), "123")
	}

	if base.GetType() != AssetTypeChart {
		t.Errorf("GetType() = %v, want %v", base.GetType(), AssetTypeChart)
	}
}
