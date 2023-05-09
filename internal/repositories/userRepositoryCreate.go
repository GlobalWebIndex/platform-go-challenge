package repositories

import (
	"context"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (repo *UserRepository) CreateUser(
	ctx context.Context,
	user *domain.User,
) (uint, error) {
	var err error

	// A new user will not have any favourite asset
	modelUser := User{
		Model:    gorm.Model{},
		Username: user.Username,
		Password: user.Password,
	}

	err = repo.db.WithContext(ctx).Create(&modelUser).Error
	if err != nil {
		return modelUser.ID, err
	}

	return modelUser.ID, nil
}
