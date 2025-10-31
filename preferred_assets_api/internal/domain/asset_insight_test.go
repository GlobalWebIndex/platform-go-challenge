package domain_test

import (
	"testing"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
	"github.com/stretchr/testify/require"
)

func newValidInsight() *domain.Insight {
	return &domain.Insight{
		AssetBase: domain.AssetBase{
			ID:    "1",
			Type:  domain.AssetType(domain.AssetTypeInsight),
			Title: "Example Insight",
		},
		Text: "Valid insight text",
	}
}

func TestInsightValidate(t *testing.T) {
	t.Run("valid insight", func(t *testing.T) {
		// Arrange
		i := newValidInsight()

		// Act
		err := i.Validate()

		// Assert
		require.NoError(t, err)
	})

	t.Run("empty text", func(t *testing.T) {
		// Arrange
		i := newValidInsight()
		i.Text = "" // empty text

		// Act
		err := i.Validate()

		// Assert
		require.Error(t, err)
		require.Contains(t, err.Error(), "insight text is required")
	})

	t.Run("text only whitespace", func(t *testing.T) {
		// Arrange
		i := newValidInsight()
		i.Text = "    " // whitespace only

		// Act
		err := i.Validate()

		// Assert
		require.Error(t, err)
		require.Contains(t, err.Error(), "insight text is required")
	})

	t.Run("AssetBase validation fails", func(t *testing.T) {
		// Arrange
		i := newValidInsight()
		i.ID = "" // AssetBase.Validate requires ID

		// Act
		err := i.Validate()

		// Assert
		require.Error(t, err)
		require.Contains(t, err.Error(), "id is required") // adjust if your AssetBase.Validate message differs
	})
}
