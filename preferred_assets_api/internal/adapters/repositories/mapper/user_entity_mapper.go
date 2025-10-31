package mapper

import (
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

// UserEntityToDomain converts entity to domain model
func UserEntityToDomain(e *entities.UserEntity) *domain.User {
	if e == nil {
		return nil
	}
	return &domain.User{
		Id:        e.Id,
		Name:      e.Name,
		Email:     e.Email,
		Password:  e.Password,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// UserEntityFromDomain converts domain model to entity
func UserEntityFromDomain(user domain.User) entities.UserEntity {
	return entities.UserEntity{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// UserEntintyFromDomainList converts slice domain models to slice entities
func UserEntintyFromDomainList(users []domain.User) []entities.UserEntity {
	ent := make([]entities.UserEntity, len(users))
	for i, user := range users {
		ent[i] = UserEntityFromDomain(user)
	}
	return ent
}

// UserEntintyToDomainList converts slice entities to slice domain models
func UserEntintyToDomainList(users []entities.UserEntity) []domain.User {
	ent := make([]domain.User, len(users))
	for i, usr := range users {
		ent[i] = *UserEntityToDomain(&usr)
	}
	return ent
}
