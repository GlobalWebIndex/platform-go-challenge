package mapper_test

import (
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/mapper"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

func TestChartEntityMapping(t *testing.T) {
	t.Run("happy path: domain -> entity -> domain", func(t *testing.T) {
		// Arrange
		now := time.Now().UTC()
		domainChart := &domain.Chart{
			AssetBase: domain.AssetBase{
				ID:          "chart1",
				Type:        domain.AssetTypeChart,
				Title:       "Sales Chart",
				Description: "Monthly sales",
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			AxesTitles: []string{"Month", "Revenue"},
			Data:       [][]float64{{1, 100}, {2, 200}},
		}

		// Act
		entity := mapper.ChartEntityFromDomain(domainChart)
		resultDomain := mapper.ChartEntityToDomain(entity)

		// Assert
		if resultDomain.ID != domainChart.ID {
			t.Errorf("expected ID %s, got %s", domainChart.ID, resultDomain.ID)
		}
		if len(resultDomain.AxesTitles) != len(domainChart.AxesTitles) {
			t.Errorf("expected AxesTitles %v, got %v", domainChart.AxesTitles, resultDomain.AxesTitles)
		}
		if len(resultDomain.Data) != len(domainChart.Data) {
			t.Errorf("expected Data %v, got %v", domainChart.Data, resultDomain.Data)
		}
	})

	t.Run("unhappy path: nil domain input", func(t *testing.T) {
		// Arrange & Act
		entity := mapper.ChartEntityFromDomain(nil)

		// Assert
		if entity != nil {
			t.Errorf("expected nil entity for nil input, got %+v", entity)
		}
	})

	t.Run("unhappy path: nil entity input", func(t *testing.T) {
		// Arrange & Act
		domainObj := mapper.ChartEntityToDomain(nil)

		// Assert
		if domainObj != nil {
			t.Errorf("expected nil domain for nil input, got %+v", domainObj)
		}
	})

	t.Run("safeUnmarshalStringArray returns empty on invalid JSON", func(t *testing.T) {
		// Arrange
		jsonStr := "invalid-json"

		// Act
		result := mapper.ChartEntityToDomain(&entities.ChartEntity{
			AssetBaseEntity: entities.AssetBaseEntity{ID: "chart2"},
			AxesTitles:      jsonStr,
		}).AxesTitles

		// Assert
		if len(result) != 0 {
			t.Errorf("expected empty slice on invalid JSON, got %+v", result)
		}
	})

	t.Run("safeUnmarshalFloatArray returns empty on invalid JSON", func(t *testing.T) {
		// Arrange
		jsonStr := "invalid-json"

		// Act
		result := mapper.ChartEntityToDomain(&entities.ChartEntity{
			AssetBaseEntity: entities.AssetBaseEntity{ID: "chart3"},
			Data:            jsonStr,
		}).Data

		// Assert
		if len(result) != 0 {
			t.Errorf("expected empty slice on invalid JSON, got %+v", result)
		}
	})

	t.Run("safeMarshalToString_returns_fallback_on_nil_value", func(t *testing.T) {
		fallback := "[]"
		chart := &domain.Chart{
			AssetBase:  domain.AssetBase{},
			AxesTitles: nil,
			Data:       nil,
		}

		entity := mapper.ChartEntityFromDomain(chart)
		if entity.AxesTitles != fallback {
			t.Errorf("expected fallback '%s' for AxesTitles, got %v", fallback, entity.AxesTitles)
		}
		if entity.Data != fallback {
			t.Errorf("expected fallback '%s' for Data, got %v", fallback, entity.Data)
		}
	})

	t.Run("nil_entity_input", func(t *testing.T) {
		var entity *entities.ChartEntity
		domainObj := mapper.ChartEntityToDomain(entity)
		if domainObj != nil {
			t.Errorf("expected nil, got %+v", domainObj)
		}
	})

}
