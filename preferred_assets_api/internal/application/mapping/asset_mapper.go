// internal/transport/http/mapping/asset_mapping.go
package mapping

import (
	"fmt"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

func AssetReqToDomain(req dto.AssetRequest) (domain.Asset, error) {
	assetType, err := domain.ParseAssetType(req.Type)
	if err != nil {
		return nil, err
	}

	base := domain.AssetBase{
		ID:          req.ID,
		Type:        assetType,
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   req.CreatedAt,
		UpdatedAt:   req.UpdatedAt,
	}

	switch assetType {
	case domain.AssetTypeAudience:
		return &domain.Audience{
			AssetBase:       base,
			Gender:          safeDerefString(req.Gender),
			BirthCountry:    safeDerefString(req.BirthCountry),
			AgeGroup:        safeDerefString(req.AgeGroup),
			HoursSocial:     safeDerefFloat64(req.HoursSocial),
			PurchasesLastMo: safeDerefInt(req.PurchasesLastMo),
		}, nil

	case domain.AssetTypeChart:
		return &domain.Chart{
			AssetBase:  base,
			AxesTitles: req.AxesTitles,
			Data:       req.Data,
		}, nil

	case domain.AssetTypeInsight:
		return &domain.Insight{
			AssetBase: base,
			// Add insight-specific fields as needed
		}, nil

	default:
		return nil, fmt.Errorf("unsupported asset type: %s", req.Type)
	}
}

// AssetDomainToCreationResponse maps a domain Asset to an AssetCreationResponse DTO
// This function converts the internal domain representation to the external API response format
func AssetDomainToCreationResponse(asset domain.Asset) dto.AssetCreationResponse {
	base := dto.AssetBaseResponse{
		ID:          asset.GetID(),
		Type:        asset.GetType().String(),
		Title:       asset.GetTitle(),
		Description: asset.GetDescription(),
		CreatedAt:   asset.GetCreatedAt(),
		UpdatedAt:   asset.GetUpdatedAt(),
	}

	switch a := asset.(type) {
	case *domain.Audience:
		return dto.AssetCreationResponse{
			AssetBaseResponse: base,
			Gender:            safeRefString(a.Gender),
			BirthCountry:      safeRefString(a.BirthCountry),
			AgeGroup:          safeRefString(a.AgeGroup),
			HoursSocial:       safeRefFloat64(a.HoursSocial),
			PurchasesLastMo:   safeRefInt(a.PurchasesLastMo),
		}

	case *domain.Chart:
		return dto.AssetCreationResponse{
			AssetBaseResponse: base,
			AxesTitles:        a.AxesTitles,
			Data:              a.Data,
		}

	case *domain.Insight:
		return dto.AssetCreationResponse{
			AssetBaseResponse: base,
			// Insight-specific fields can be added here as needed
		}

	default:
		// Fallback for unknown types - return base response only
		return dto.AssetCreationResponse{
			AssetBaseResponse: base,
		}
	}
}

// Helper functions for safe referencing
func safeRefString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func safeRefInt(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

// Helper functions for safe dereferencing
func safeDerefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func safeDerefInt(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}
func safeRefFloat64(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}
func safeDerefFloat64(i *float64) float64 {
	if i == nil {
		return 0
	}
	return *i
}
