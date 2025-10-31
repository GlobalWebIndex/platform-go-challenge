package mapper

import (
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

// FavouriteEntityToDomain converts entity to domain model
func FavouriteEntityToDomain(e *entities.FavouriteEntity) *domain.Favourite {
	if e == nil {
		return nil
	}
	return &domain.Favourite{
		UserID:    e.UserId,
		AssetID:   e.AssetId,
		CreatedAt: e.CreatedAt,
	}
}

// FavouriteEntityFromDomain converts domain model to entity
func FavouriteEntityFromDomain(favourite domain.Favourite) entities.FavouriteEntity {
	return entities.FavouriteEntity{
		UserId:    favourite.UserID,
		AssetId:   favourite.AssetID,
		CreatedAt: favourite.CreatedAt,
	}
}

// FavouriteEntityToDomainList converts a list of entity to a list of domain model, handling nil input
func FavouriteEntityToDomainList(favs []entities.FavouriteEntity) []domain.Favourite {
	if favs == nil {
		return []domain.Favourite{}
	}

	ent := make([]domain.Favourite, len(favs))
	for i, f := range favs {
		ent[i] = *FavouriteEntityToDomain(&f)
	}
	return ent
}
