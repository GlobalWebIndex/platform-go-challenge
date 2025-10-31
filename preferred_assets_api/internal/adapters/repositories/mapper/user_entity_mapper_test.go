package mapper_test

import (
	"testing"
	"time"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/mapper"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

func TestUserEntityToDomain_HappyPath(t *testing.T) {
	// Arrange
	now := time.Now().UTC()
	usrEntity := entities.UserEntity{
		Id:        "1",
		Name:      "Alice",
		Email:     "alice@example.com",
		Password:  "hashed",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Act
	domainUser := mapper.UserEntityToDomain(&usrEntity)

	// Assert
	if domainUser.Id != usrEntity.Id {
		t.Errorf("expected Id %s, got %s", usrEntity.Id, domainUser.Id)
	}
	if domainUser.Name != usrEntity.Name {
		t.Errorf("expected Name %s, got %s", usrEntity.Name, domainUser.Name)
	}
	if domainUser.Email != usrEntity.Email {
		t.Errorf("expected Email %s, got %s", usrEntity.Email, domainUser.Email)
	}
	if domainUser.Password != usrEntity.Password {
		t.Errorf("expected Password %s, got %s", usrEntity.Password, domainUser.Password)
	}
	if !domainUser.CreatedAt.Equal(usrEntity.CreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", usrEntity.CreatedAt, domainUser.CreatedAt)
	}
	if !domainUser.UpdatedAt.Equal(usrEntity.UpdatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", usrEntity.UpdatedAt, domainUser.UpdatedAt)
	}
}

func TestUserEntityFromDomain_HappyPath(t *testing.T) {
	// Arrange
	now := time.Now().UTC()
	user := domain.User{
		Id:        "1",
		Name:      "Alice",
		Email:     "alice@example.com",
		Password:  "hashed",
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Act
	entity := mapper.UserEntityFromDomain(user)

	// Assert
	if entity.Id != user.Id {
		t.Errorf("expected Id %s, got %s", user.Id, entity.Id)
	}
	if entity.Name != user.Name {
		t.Errorf("expected Name %s, got %s", user.Name, entity.Name)
	}
	if entity.Email != user.Email {
		t.Errorf("expected Email %s, got %s", user.Email, entity.Email)
	}
	if entity.Password != user.Password {
		t.Errorf("expected Password %s, got %s", user.Password, entity.Password)
	}
	if !entity.CreatedAt.Equal(user.CreatedAt) {
		t.Errorf("expected CreatedAt %v, got %v", user.CreatedAt, entity.CreatedAt)
	}
	if !entity.UpdatedAt.Equal(user.UpdatedAt) {
		t.Errorf("expected UpdatedAt %v, got %v", user.UpdatedAt, entity.UpdatedAt)
	}
}

func TestUserEntityListConversions(t *testing.T) {
	// Arrange
	now := time.Now().UTC()
	domainUsers := []domain.User{
		{Id: "1", Name: "Alice", Email: "alice@example.com", CreatedAt: now, UpdatedAt: now},
		{Id: "2", Name: "Bob", Email: "bob@example.com", CreatedAt: now, UpdatedAt: now},
	}

	// Act
	entitiesList := mapper.UserEntintyFromDomainList(domainUsers)
	convertedDomain := mapper.UserEntintyToDomainList(entitiesList)

	// Assert
	if len(entitiesList) != len(domainUsers) {
		t.Errorf("expected %d entities, got %d", len(domainUsers), len(entitiesList))
	}
	if len(convertedDomain) != len(domainUsers) {
		t.Errorf("expected %d domain users, got %d", len(domainUsers), len(convertedDomain))
	}
	for i := range domainUsers {
		if convertedDomain[i].Id != domainUsers[i].Id {
			t.Errorf("expected Id %s, got %s", domainUsers[i].Id, convertedDomain[i].Id)
		}
		if convertedDomain[i].Name != domainUsers[i].Name {
			t.Errorf("expected Name %s, got %s", domainUsers[i].Name, convertedDomain[i].Name)
		}
		if convertedDomain[i].Email != domainUsers[i].Email {
			t.Errorf("expected Email %s, got %s", domainUsers[i].Email, convertedDomain[i].Email)
		}
	}
}

func TestEmptyListConversions(t *testing.T) {
	// Arrange
	var emptyDomain []domain.User
	var emptyEntities []entities.UserEntity

	// Act
	convertedEntities := mapper.UserEntintyFromDomainList(emptyDomain)
	convertedDomain := mapper.UserEntintyToDomainList(emptyEntities)

	// Assert
	if len(convertedEntities) != 0 {
		t.Errorf("expected 0 entities, got %d", len(convertedEntities))
	}
	if len(convertedDomain) != 0 {
		t.Errorf("expected 0 domain users, got %d", len(convertedDomain))
	}
}

func TestSafeUserEntityToDomain(t *testing.T) {
	now := time.Now().UTC()

	t.Run("happy path", func(t *testing.T) {
		// Arrange
		entity := &entities.UserEntity{
			Id:        "1",
			Name:      "Alice",
			Email:     "alice@example.com",
			Password:  "hashed",
			CreatedAt: now,
			UpdatedAt: now,
		}

		// Act
		user := mapper.UserEntityToDomain(entity)

		// Assert
		if user == nil {
			t.Fatal("expected non-nil user")
		}
		if user.Id != entity.Id {
			t.Errorf("expected Id %s, got %s", entity.Id, user.Id)
		}
	})

	t.Run("nil input", func(t *testing.T) {
		// Arrange
		var entity *entities.UserEntity = nil

		// Act
		user := mapper.UserEntityToDomain(entity)

		// Assert
		if user != nil {
			t.Errorf("expected nil user for nil input, got %+v", user)
		}
	})
}
