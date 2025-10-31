package mapping_test

import (
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/mapping"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
	"github.com/stretchr/testify/assert"
)

// --- FavouriteReqToDomain ---
func TestFavouriteReqToDomain(t *testing.T) {
	// Arrange
	req := dto.FavouriteRequest{UserId: "user1", AssetId: "asset1"}

	// Act
	fav := mapping.FavouriteReqToDomain(req)

	// Assert
	assert.Equal(t, "user1", fav.UserID)
	assert.Equal(t, "asset1", fav.AssetID)
}

// --- FavouriteToResponse: nil asset ---
func TestFavouriteToResponse_NilAsset(t *testing.T) {
	// Arrange
	favDomain := domain.Favourite{UserID: "user1", AssetID: "asset1", CreatedAt: time.Now()}

	// Act
	res := mapping.FavouriteToResponse(&favDomain)

	// Assert
	assert.Nil(t, res.Asset)
	assert.Equal(t, "user1", res.UserID)
	assert.Equal(t, "asset1", res.AssetID)
	assert.WithinDuration(t, favDomain.CreatedAt, res.CreatedAt, time.Millisecond)
}

// --- FavouriteToResponse: Audience ---
func TestFavouriteToResponse_WithAudience(t *testing.T) {
	// Arrange
	audience := &domain.Audience{
		AssetBase: domain.AssetBase{
			ID:          "a1",
			Type:        domain.AssetTypeAudience,
			Title:       "Audience1",
			Description: "Test Audience",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Gender:          "M",
		BirthCountry:    "US",
		AgeGroup:        "25-34",
		HoursSocial:     4.5,
		PurchasesLastMo: 12,
	}

	favDomain := domain.Favourite{UserID: "u1", AssetID: "a1", CreatedAt: time.Now()}
	err := favDomain.SetAsset(audience)
	assert.NoError(t, err)

	// Act
	res := mapping.FavouriteToResponse(&favDomain)

	// Assert
	assert.Equal(t, "u1", res.UserID)
	assert.Equal(t, "a1", res.AssetID)

	dtoAsset, ok := res.Asset.(dto.AssetRequest)
	assert.True(t, ok, "Expected AssetRequest, got %T", res.Asset)
	assert.Equal(t, "a1", dtoAsset.ID)
	assert.Equal(t, "audience", dtoAsset.Type)
	assert.Equal(t, "Audience1", dtoAsset.Title)
	assert.Equal(t, "Test Audience", dtoAsset.Description)
	assert.Equal(t, "M", *dtoAsset.Gender)
	assert.Equal(t, "US", *dtoAsset.BirthCountry)
	assert.Equal(t, "25-34", *dtoAsset.AgeGroup)
	assert.Equal(t, 4.5, *dtoAsset.HoursSocial)
	assert.Equal(t, 12, *dtoAsset.PurchasesLastMo)
}

// --- FavouriteToResponse: Chart ---
func TestFavouriteToResponse_WithChart(t *testing.T) {
	// Arrange
	chart := &domain.Chart{
		AssetBase: domain.AssetBase{
			ID:          "c1",
			Type:        domain.AssetTypeChart,
			Title:       "Chart1",
			Description: "Test Chart",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		AxesTitles: []string{"X", "Y"},
		Data:       [][]float64{{1, 2}, {3, 4}},
	}

	favDomain := domain.Favourite{UserID: "u1", AssetID: "c1", CreatedAt: time.Now()}
	err := favDomain.SetAsset(chart)
	assert.NoError(t, err)

	// Act
	res := mapping.FavouriteToResponse(&favDomain)

	// Assert
	assert.Equal(t, "u1", res.UserID)
	assert.Equal(t, "c1", res.AssetID)

	dtoAsset, ok := res.Asset.(dto.AssetRequest)
	assert.True(t, ok, "Expected AssetRequest, got %T", res.Asset)
	assert.Equal(t, "c1", dtoAsset.ID)
	assert.Equal(t, "chart", dtoAsset.Type)
	assert.Equal(t, "Chart1", dtoAsset.Title)
	assert.Equal(t, "Test Chart", dtoAsset.Description)
	assert.Equal(t, []string{"X", "Y"}, dtoAsset.AxesTitles)
	assert.Equal(t, [][]float64{{1, 2}, {3, 4}}, dtoAsset.Data)
}

// --- FavouriteToResponse: Insight ---
func TestFavouriteToResponse_WithInsight(t *testing.T) {
	// Arrange
	insight := &domain.Insight{
		AssetBase: domain.AssetBase{
			ID:          "i1",
			Type:        domain.AssetTypeInsight,
			Title:       "Insight1",
			Description: "Test Insight",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		Text: "Some insight text",
	}

	favDomain := domain.Favourite{UserID: "u1", AssetID: "i1", CreatedAt: time.Now()}
	err := favDomain.SetAsset(insight)
	assert.NoError(t, err)

	// Act
	res := mapping.FavouriteToResponse(&favDomain)

	// Assert
	assert.Equal(t, "u1", res.UserID)
	assert.Equal(t, "i1", res.AssetID)

	dtoAsset, ok := res.Asset.(dto.AssetRequest)
	assert.True(t, ok, "Expected AssetRequest, got %T", res.Asset)
	assert.Equal(t, "i1", dtoAsset.ID)
	assert.Equal(t, "insight", dtoAsset.Type)
	assert.Equal(t, "Insight1", dtoAsset.Title)
	assert.Equal(t, "Test Insight", dtoAsset.Description)
	assert.Equal(t, "Some insight text", *dtoAsset.Text)
}

// --- FavouritesToResponse: mixed assets and nil ---
func TestFavouritesToResponse_MixedAssets(t *testing.T) {
	// Arrange
	aud := &domain.Audience{
		AssetBase: domain.AssetBase{
			ID:        "a1",
			Type:      domain.AssetTypeAudience,
			Title:     "Audience1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	chart := &domain.Chart{
		AssetBase: domain.AssetBase{
			ID:        "c1",
			Type:      domain.AssetTypeChart,
			Title:     "Chart1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	fav1 := domain.Favourite{UserID: "u1", AssetID: "a1", CreatedAt: time.Now()}
	err1 := fav1.SetAsset(aud)
	assert.NoError(t, err1)

	fav2 := domain.Favourite{UserID: "u2", AssetID: "c1", CreatedAt: time.Now()}
	err2 := fav2.SetAsset(chart)
	assert.NoError(t, err2)

	fav3 := domain.Favourite{UserID: "u3", AssetID: "i1", CreatedAt: time.Now()} // nil asset

	favs := []domain.Favourite{fav1, fav2, fav3}

	// Act
	responses := mapping.FavouritesToResponse(favs)

	// Assert
	assert.Len(t, responses, 3)
	assert.NotNil(t, responses[0].Asset)
	assert.NotNil(t, responses[1].Asset)
	assert.Nil(t, responses[2].Asset)

	// Verify the types are correct
	dtoAsset1, ok1 := responses[0].Asset.(dto.AssetRequest)
	assert.True(t, ok1, "First asset should be AssetRequest")
	assert.Equal(t, "audience", dtoAsset1.Type)

	dtoAsset2, ok2 := responses[1].Asset.(dto.AssetRequest)
	assert.True(t, ok2, "Second asset should be AssetRequest")
	assert.Equal(t, "chart", dtoAsset2.Type)
}

// --- Unhappy path: nil Favourite pointer (domain object) ---
func TestFavouriteToResponse_NilDomainObject(t *testing.T) {
	// Arrange
	var fav *domain.Favourite // nil pointer

	// Act
	res := mapping.FavouriteToResponse(fav)

	// Assert
	assert.Equal(t, "", res.UserID)
	assert.Equal(t, "", res.AssetID)
	assert.Nil(t, res.Asset)
	assert.True(t, res.CreatedAt.IsZero())
}
