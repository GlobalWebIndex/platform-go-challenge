package sqldb

import (
	"context"
	"platform-go-challenge/domain"
)

func (d *DB) AddUser(ctx context.Context, user domain.User) (*domain.User, error) {
	u := &User{}
	u.FromDomain(&user)
	err := d.db.Create(u).Error
	if err != nil {
		return nil, err
	}
	nu := u.ToDomain()
	return nu, nil
}
func (d *DB) FindUser(ctx context.Context, cred domain.LoginCredentials) (*domain.User, error) {
	u := User{}
	err := d.db.Where("username = ? AND password = ?", cred.Username, cred.Password).First(&u).Error
	if err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}
