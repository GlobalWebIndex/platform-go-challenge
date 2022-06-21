package domain

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	_, err := dom.CreateUser(ctx, User{
		Username: "manos",
		Password: "",
	})
	assert.ErrorIs(t, err, ErrWrongUserInput)
	_, err = dom.CreateUser(ctx, User{
		Username: "",
		Password: "secret",
	})
	assert.ErrorIs(t, err, ErrWrongUserInput)
}

func TestLoginUser(t *testing.T) {
	dom := NewDomain(&MockDB{})
	ctx := context.Background()
	_, err := dom.LoginUser(ctx, LoginCredentials{
		Username: "manos",
		Password: "",
	})
	assert.ErrorIs(t, err, ErrWrongLoginInput)
	_, err = dom.LoginUser(ctx, LoginCredentials{
		Username: "",
		Password: "secret",
	})
	assert.ErrorIs(t, err, ErrWrongLoginInput)
}
