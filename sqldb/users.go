package sqldb

import (
	"context"
	"platform-go-challenge/domain"
)

func (d *DB) CreateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	return nil, nil
}
func (d *DB) FindUser(ctx context.Context, cred domain.LoginCredentials) (*domain.User, error) {
	return nil, nil
}
