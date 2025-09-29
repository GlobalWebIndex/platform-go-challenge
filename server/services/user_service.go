package services

import (
	"github.com/google/uuid"
	"favourite_assets/server/models"
	"favourite_assets/server/repositories"
	"favourite_assets/server/errors"
)

type UserService struct {
	repo *repositories.UserRepository
}


func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(name, email string) (*models.User, error) {

	for _, user := range s.repo.List() {
		if user.Email == email {
			return user, errors.ErrUserExists 
		}
	}

	user := &models.User{
		ID:    uuid.New(),
		Name:  name,
		Email: email,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(id uuid.UUID) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) UpdateUser(id uuid.UUID, name, email string) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.Name = name
	user.Email = email

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *UserService) ListUsers() []*models.User {
	return s.repo.List()
}
