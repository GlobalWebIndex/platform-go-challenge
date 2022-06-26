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

func (d *DB) FindUser(ctx context.Context, username string) (*domain.User, error) {
	u := User{}
	err := d.db.Where("username = ? ", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}

func (d *DB) UserExists(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := d.db.Model(&User{}).Select("count(*) > 0").Where("username = ? ", username).Find(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (d *DB) GetUser(ctx context.Context, userID uint) (*domain.User, error) {
	u := User{}
	err := d.db.First(&u, userID).Error
	if err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}
