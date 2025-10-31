package services_test

import (
	"errors"
	"testing"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/repositories/entities"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/services"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
)

// Mocks
type mockUserRepo struct {
	users map[string]entities.UserEntity
}

func (m *mockUserRepo) GetByID(id string) (entities.UserEntity, error) {
	u, ok := m.users[id]
	if !ok {
		return entities.UserEntity{}, errors.New("not found")
	}
	return u, nil
}

func (m *mockUserRepo) GetAll() ([]entities.UserEntity, error) {
	list := []entities.UserEntity{}
	for _, u := range m.users {
		list = append(list, u)
	}
	return list, nil
}

func (m *mockUserRepo) Save(user entities.UserEntity) error {
	if m.users == nil {
		m.users = make(map[string]entities.UserEntity)
	}
	m.users[user.Id] = user
	return nil
}
func (m *mockUserRepo) Update(user entities.UserEntity) error {
	if m.users == nil {
		m.users = make(map[string]entities.UserEntity)
	}
	m.users[user.Id] = user
	return nil
}

func (m *mockUserRepo) Delete(id string) error { return nil }
func (m *mockUserRepo) GetFavouritesByID(id string) ([]entities.FavouriteEntity, error) {
	return []entities.FavouriteEntity{
		{UserId: "1", AssetId: "a1"},
		{UserId: "1", AssetId: "a2"},
	}, nil
}

type mockAssetRepository struct {
	assets map[string]entities.AssetEntity
}

func newMockAssetRepo() *mockAssetRepository {
	return &mockAssetRepository{
		assets: make(map[string]entities.AssetEntity),
	}
}

func (m *mockAssetRepository) GetByIDs(ids []string) ([]entities.AssetEntity, error) {
	assets := make([]entities.AssetEntity, 0, len(ids))
	for _, id := range ids {
		assets = append(assets, &entities.AssetBaseEntity{
			ID:    id,
			Type:  entities.AssetTypeChart,
			Title: "Mock Asset " + id,
		})
	}
	return assets, nil
}

func (m *mockAssetRepository) Save(asset entities.AssetEntity) (entities.AssetEntity, error) {
	return nil, nil
}
func (m *mockAssetRepository) GetByID(id string) (entities.AssetEntity, error) { return nil, nil }
func (m *mockAssetRepository) GetAll() ([]entities.AssetEntity, error)         { return nil, nil }
func (m *mockAssetRepository) GetByType(assetType entities.AssetType) ([]entities.AssetEntity, error) {
	return nil, nil
}
func (m *mockAssetRepository) Update(asset entities.AssetEntity) error { return nil }
func (m *mockAssetRepository) Delete(id string) error                  { return nil }
func (m *mockAssetRepository) Exists(id string) (bool, error)          { return true, nil }

// Happy PathTests

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		users     map[string]entities.UserEntity
		wantError bool
	}{
		{"user exists", "1", map[string]entities.UserEntity{"1": {Id: "1", Name: "Alice"}}, false},
		{"user missing", "42", map[string]entities.UserEntity{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			repo := &mockUserRepo{users: tt.users}
			assetRepo := newMockAssetRepo()
			service := services.NewUserService(repo, assetRepo)

			// Act
			user, err := service.GetUserByID(tt.userID)

			// Assert
			if tt.wantError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if user != nil {
					t.Errorf("expected nil user, got %v", user)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if user.Id != tt.userID {
					t.Errorf("expected ID '%s', got '%s'", tt.userID, user.Id)
				}
			}
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	// Arrange
	users := map[string]entities.UserEntity{
		"1": {Id: "1", Name: "Alice"},
		"2": {Id: "2", Name: "Bob"},
	}
	service := services.NewUserService(&mockUserRepo{users: users}, &mockAssetRepository{})

	// Act
	list, err := service.GetAllUsers()

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(list) != len(users) {
		t.Errorf("expected %d users, got %d", len(users), len(list))
	}
}

func TestCreateUser_HashesPassword(t *testing.T) {
	// Arrange
	service := services.NewUserService(&mockUserRepo{}, &mockAssetRepository{})
	user := domain.User{
		Id:       "1",
		Name:     "Alice",
		Password: "plain",
	}

	// Act
	err := service.CreateUser(user)
	userEntity, _ := service.GetUserByID("1")

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userEntity.Password == "plain" {
		t.Error("expected password to be hashed")
	}
	if userEntity.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}
}

func TestDeleteUser(t *testing.T) {
	// Arrange
	service := services.NewUserService(&mockUserRepo{}, &mockAssetRepository{})

	// Act
	err := service.DeleteUser("1")

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	// Arrange
	service := services.NewUserService(&mockUserRepo{}, &mockAssetRepository{})
	user := domain.User{Id: "1", Name: "Alice"}

	// Act
	err := service.UpdateUser(user)
	userEntity, _ := service.GetUserByID("1")

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if userEntity.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestGetFavouritesByUser(t *testing.T) {
	// Arrange
	service := services.NewUserService(&mockUserRepo{}, &mockAssetRepository{})

	// Act
	favs, err := service.GetFavouritesByUser("1")

	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(favs) != 0 {
		t.Errorf("expected 0 favourites, got %d", len(favs))
	}
	for _, fav := range favs {
		if fav.AssetID == "" {
			t.Errorf("favourite %s has nil Asset", fav.UserID)
		}
	}
}
