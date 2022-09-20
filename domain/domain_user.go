package domain

import (
	"context"
	"errors"
	"fmt"
)

func (d *Domain) CreateUser(ctx context.Context, user User) (*User, error) {
	err := d.validate.Struct(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWrongUserInput, err)
	}

	exists, err := d.repo.UserExists(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInternalDBFailure, err)
	}
	if exists {
		return nil, fmt.Errorf("%w: %v", ErrWrongUserInput, errors.New("user exists"))
	}

	pass, err := hashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWrongUserInput, err)
	}
	user.Password = pass
	newUser, err := d.repo.AddUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWrongUserInput, err)
	}
	newUser.Password = ""
	return newUser, nil
}

func (d *Domain) LoginUser(ctx context.Context, cred LoginCredentials) (*User, error) {
	err := d.validate.Struct(cred)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrWrongLoginInput, err)
	}
	user, err := d.repo.FindUser(ctx, cred.Username)
	if err != nil {
		return nil, ErrUserNotFound
	}
	ok := checkPasswordHash(cred.Password, user.Password)
	if !ok {
		return nil, ErrUnauthorized
	}

	user.Password = ""
	return user, nil
}
