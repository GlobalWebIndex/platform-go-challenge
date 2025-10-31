package mapper

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

// ChartEntityToDomain converts entity to domain model
func ChartEntityToDomain(e *entities.ChartEntity) *domain.Chart {
	if e == nil {
		return nil
	}

	return &domain.Chart{
		AssetBase:  AssetBaseEntityToDomain(e.AssetBaseEntity),
		AxesTitles: safeUnmarshalStringArray(e.AxesTitles, e.ID, "axes titles"),
		Data:       safeUnmarshalFloatArray(e.Data, e.ID, "chart data"),
	}
}

// ChartEntityFromDomain converts domain model to entity
func ChartEntityFromDomain(c *domain.Chart) *entities.ChartEntity {
	if c == nil {
		return nil
	}

	return &entities.ChartEntity{
		AssetBaseEntity: *AssetBaseEntityFromDomain(c.AssetBase),
		AxesTitles:      safeMarshalToString(c.AxesTitles, "[]", fmt.Sprintf("axes titles for chart %s", c.GetID())),
		Data:            safeMarshalToString(c.Data, "[]", fmt.Sprintf("chart data for chart %s", c.GetID())),
	}
}

func safeUnmarshalStringArray(jsonStr, assetID, fieldName string) []string {
	var result []string
	if jsonStr == "" {
		return []string{}
	}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		log.Printf("Warning: failed to unmarshal %s for asset %s: %v", fieldName, assetID, err)
		return []string{}
	}
	return result
}

func safeUnmarshalFloatArray(jsonStr, assetID, fieldName string) [][]float64 {
	var result [][]float64
	if jsonStr == "" {
		return [][]float64{}
	}
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		log.Printf("Warning: failed to unmarshal %s for asset %s: %v", fieldName, assetID, err)
		return [][]float64{}
	}
	return result
}

func safeMarshalToString(value interface{}, fallback string, logPrefix string) string {
	if value == nil {
		return fallback
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		log.Printf("WARNING: %s - failed to marshal: %v", logPrefix, err)
		return fallback
	}

	// Ensure empty slices are marshaled as "[]"
	if string(bytes) == "null" {
		return fallback
	}

	return string(bytes)
}
