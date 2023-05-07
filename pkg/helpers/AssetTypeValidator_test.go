package helpers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssetTypeValidator_ValidateAssetType(t *testing.T) {
	type args struct {
		assetType string
	}
	tests := []struct {
		name     string
		args     args
		expected error
	}{
		{
			name: "chart",
			args: args{
				assetType: "chart",
			},
			expected: nil,
		},
		{
			name: "insight",
			args: args{
				assetType: "insight",
			},
			expected: nil,
		},
		{
			name: "audience",
			args: args{
				assetType: "audience",
			},
			expected: nil,
		},
		{
			name: "unknown",
			args: args{
				assetType: "fjslkdfjaslkdfjas",
			},
			expected: errors.New("unknown asset type: fjslkdfjaslkdfjas"),
		},
		{
			name:     "empty string",
			args:     args{},
			expected: errors.New("unknown asset type: "),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := AssetTypeValidator{}

			actual := validator.ValidateAssetType(tt.args.assetType)

			assert.Equal(t, tt.expected, actual)
		})
	}
}
