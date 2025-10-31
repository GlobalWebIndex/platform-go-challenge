package mapper_test

import (
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/mapper"
)

func TestFavouriteEntityToDomain(t *testing.T) {
	now := time.Now().UTC()

	t.Run("happy path", func(t *testing.T) {
		// Arrange
		entity := &entities.FavouriteEntity{
			UserId:    "user1",
			AssetId:   "asset1",
			CreatedAt: now,
		}

		// Act
		fav := mapper.FavouriteEntityToDomain(entity)

		// Assert
		if fav == nil {
			t.Fatal("expected non-nil favourite")
		}
		if fav.UserID != entity.UserId {
			t.Errorf("expected UserID %s, got %s", entity.UserId, fav.UserID)
		}
	})

	t.Run("nil input", func(t *testing.T) {
		// Arrange
		var entity *entities.FavouriteEntity = nil

		// Act
		fav := mapper.FavouriteEntityToDomain(entity)

		// Assert
		if fav != nil {
			t.Errorf("expected nil favourite for nil input, got %+v", fav)
		}
	})
}

func TestFavouriteEntityToDomainList(t *testing.T) {
	now := time.Now().UTC()

	t.Run("happy path with multiple favourites", func(t *testing.T) {
		// Arrange
		entitiesList := []entities.FavouriteEntity{
			{UserId: "u1", AssetId: "a1", CreatedAt: now},
			{UserId: "u2", AssetId: "a2", CreatedAt: now.Add(time.Minute)},
		}

		// Act
		domainList := mapper.FavouriteEntityToDomainList(entitiesList)

		// Assert
		if len(domainList) != 2 {
			t.Fatalf("expected 2 favourites, got %d", len(domainList))
		}
		if domainList[0].UserID != "u1" || domainList[1].UserID != "u2" {
			t.Errorf("user IDs mismatch: %+v", domainList)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		// Arrange
		var entitiesList []entities.FavouriteEntity

		// Act
		domainList := mapper.FavouriteEntityToDomainList(entitiesList)

		// Assert
		if len(domainList) != 0 {
			t.Errorf("expected empty slice, got %+v", domainList)
		}
	})

	t.Run("nil slice input", func(t *testing.T) {
		// Act
		domainList := mapper.FavouriteEntityToDomainList(nil)

		// Assert
		if len(domainList) != 0 {
			t.Errorf("expected empty slice for nil input, got %+v", domainList)
		}
	})
}
