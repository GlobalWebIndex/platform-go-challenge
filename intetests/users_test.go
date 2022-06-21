package intetests

import (
	"context"
	"platform-go-challenge/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	dom, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	user, err := dom.CreateUser(ctx, domain.User{
		Username: "manos",
		Password: "password",
		IsAdmin:  true,
	})
	assert.NoError(t, err)
	assert.NotNil(t, user)

	user2, err := dom.LoginUser(ctx, domain.LoginCredentials{
		Username: "manos",
		Password: "password",
	})
	assert.NoError(t, err)
	assert.NotNil(t, user2)
	assert.EqualValues(t, user, user2)

	_, err = dom.CreateUser(ctx, domain.User{
		Username: "manos",
		Password: "password",
		IsAdmin:  true,
	})
	assert.Error(t, err)
}
