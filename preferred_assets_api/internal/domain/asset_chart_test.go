package domain_test

import (
	"testing"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"

	"github.com/stretchr/testify/require"
)

func TestChartValidate(t *testing.T) {
	t.Run("valid chart", func(t *testing.T) {
		// Arrange
		chart := &domain.Chart{
			AssetBase: domain.AssetBase{
				ID:    "1",
				Type:  domain.AssetTypeChart,
				Title: "Sales Data",
			},
			AxesTitles: []string{"Month", "Revenue"},
			Data: [][]float64{
				{1, 100},
				{2, 200},
				{3, 300},
			},
		}

		// Act
		err := chart.Validate()

		// Assert
		require.NoError(t, err)
	})

	t.Run("too many axes titles", func(t *testing.T) {
		// Arrange
		chart := &domain.Chart{
			AssetBase: domain.AssetBase{
				ID:    "2",
				Type:  domain.AssetTypeChart,
				Title: "Invalid Axes",
			},
			AxesTitles: []string{"X", "Y", "Z"},
			Data: [][]float64{
				{1, 2, 3},
			},
		}

		// Act
		err := chart.Validate()

		// Assert
		require.Error(t, err)
		require.Contains(t, err.Error(), "maximum 2 axes titles allowed")
	})

	t.Run("inconsistent row lengths", func(t *testing.T) {
		// Arrange
		chart := &domain.Chart{
			AssetBase: domain.AssetBase{
				ID:    "3",
				Type:  domain.AssetTypeChart,
				Title: "Bad Data",
			},
			AxesTitles: []string{"X", "Y"},
			Data: [][]float64{
				{1, 2},
				{3}, // inconsistent length
				{4, 5},
			},
		}

		// Act
		err := chart.Validate()

		// Assert
		require.Error(t, err)
		require.Contains(t, err.Error(), "all data rows must have the same length")
	})

	t.Run("empty data", func(t *testing.T) {
		// Arrange
		chart := &domain.Chart{
			AssetBase: domain.AssetBase{
				ID:    "4",
				Type:  domain.AssetTypeChart,
				Title: "Empty Data",
			},
			AxesTitles: []string{"X", "Y"},
			Data:       [][]float64{},
		}

		// Act
		err := chart.Validate()

		// Assert
		require.Error(t, err)
		require.Contains(t, err.Error(), "data cannot be empty")
	})
}
