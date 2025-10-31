package mapping

import (
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

func FavouriteReqToDomain(req dto.FavouriteRequest) domain.Favourite {
	return domain.Favourite{
		UserID:  req.UserId,
		AssetID: req.AssetId,
	}
}

// Single favourite to DTO
func FavouriteToResponse(fav *domain.Favourite) dto.FavouriteResponse {
	if fav == nil {
		// Return an empty DTO if the domain object is nil
		return dto.FavouriteResponse{}
	}

	return dto.FavouriteResponse{
		UserID:    fav.UserID,
		AssetID:   fav.AssetID,
		Asset:     mapAssetToDTO(fav.GetAsset()),
		CreatedAt: fav.CreatedAt,
	}
}

// Multiple favourites to DTOs
func FavouritesToResponse(favourites []domain.Favourite) []dto.FavouriteResponse {
	responses := make([]dto.FavouriteResponse, len(favourites))
	for i, f := range favourites {
		responses[i] = FavouriteToResponse(&f)
	}
	return responses
}

// Map the asset to the appropriate DTO
func mapAssetToDTO(asset domain.Asset) interface{} {
	if asset == nil {
		return nil
	}

	// Extra check: in case asset is a typed nil pointer
	switch a := asset.(type) {
	case *domain.Audience:
		if a == nil {
			return nil
		}
		return mapAudienceToDTO(a)
	case *domain.Chart:
		if a == nil {
			return nil
		}
		return mapChartToDTO(a)
	case *domain.Insight:
		if a == nil {
			return nil
		}
		return mapInsightToDTO(a)
	default:
		return nil
	}
}

func mapAudienceToDTO(audience *domain.Audience) dto.AssetRequest {
	return dto.AssetRequest{
		ID:              audience.GetID(),
		Type:            assetTypeToString(audience.Type),
		Title:           audience.GetTitle(),
		Description:     audience.GetDescription(),
		Gender:          &audience.Gender,
		BirthCountry:    &audience.BirthCountry,
		AgeGroup:        &audience.AgeGroup,
		HoursSocial:     safeRefFloat64(audience.HoursSocial),
		PurchasesLastMo: &audience.PurchasesLastMo,
		CreatedAt:       audience.GetCreatedAt(),
		UpdatedAt:       audience.GetUpdatedAt(),
	}
}

func mapChartToDTO(chart *domain.Chart) dto.AssetRequest {
	return dto.AssetRequest{
		ID:          chart.GetID(),
		Type:        assetTypeToString(chart.Type),
		Title:       chart.GetTitle(),
		Description: chart.GetDescription(),
		AxesTitles:  chart.AxesTitles,
		Data:        chart.Data,
		CreatedAt:   chart.GetCreatedAt(),
		UpdatedAt:   chart.GetUpdatedAt(),
	}
}

func mapInsightToDTO(insight *domain.Insight) dto.AssetRequest {
	return dto.AssetRequest{
		ID:          insight.GetID(),
		Type:        assetTypeToString(insight.Type),
		Title:       insight.GetTitle(),
		Description: insight.GetDescription(),
		Text:        &insight.Text,
		CreatedAt:   insight.GetCreatedAt(),
		UpdatedAt:   insight.GetUpdatedAt(),
	}
}

// Helper to convert AssetType to string
func assetTypeToString(assetType domain.AssetType) string {
	switch assetType {
	case domain.AssetTypeAudience:
		return "audience"
	case domain.AssetTypeChart:
		return "chart"
	case domain.AssetTypeInsight:
		return "insight"
	default:
		return "unknown"
	}
}
