package domain_test

import (
	"testing"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"

	"github.com/stretchr/testify/require"
)

// Helper for a valid Audience
func newValidAudience() *domain.Audience {
	return &domain.Audience{
		AssetBase: domain.AssetBase{
			ID:    "aud-1",
			Type:  domain.AssetTypeAudience,
			Title: "Audience Example",
		},
		Gender:          "female",
		BirthCountry:    "Canada",
		AgeGroup:        "25-34",
		HoursSocial:     3,
		PurchasesLastMo: 5,
	}
}

func TestAudience_Validate(t *testing.T) {
	t.Run("valid audience", func(t *testing.T) {
		// Arrange
		a := newValidAudience()

		// Act
		err := a.Validate()

		// Assert
		require.NoError(t, err)
	})

	t.Run("empty ID fails validation", func(t *testing.T) {
		// Arrange
		a := newValidAudience()
		a.ID = "" // AssetBase requires ID

		// Act
		err := a.Validate()

		// Assert
		require.Error(t, err)
	})

	t.Run("empty Title fails validation", func(t *testing.T) {
		// Arrange
		a := newValidAudience()
		a.Title = "" // AssetBase requires Title

		// Act
		err := a.Validate()

		// Assert
		require.Error(t, err)
	})
}
