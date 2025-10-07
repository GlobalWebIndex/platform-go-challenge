package service

import "context"

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type UsersService struct{ Repo UsersRepository }

type UsersRepository interface {
	Create(ctx context.Context, email string) (int64, error)
	List(ctx context.Context, limit, offset int) ([]User, error)
	Get(ctx context.Context, id int64) (User, error)
}

func (s *UsersService) Create(ctx context.Context, email string) (int64, error) {
	return s.Repo.Create(ctx, email)
}

func (s *UsersService) List(ctx context.Context, limit, offset int) ([]User, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.Repo.List(ctx, limit, offset)
}

func (s *UsersService) Get(ctx context.Context, id int64) (User, error) {
	return s.Repo.Get(ctx, id)
}
