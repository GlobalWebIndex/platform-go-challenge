package sqldb

import (
	"context"
	"platform-go-challenge/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCRUser(t *testing.T) {
	db, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	ctx := context.Background()
	du := domain.User{
		Username: "manos",
		Password: "hashed",
		IsAdmin:  false,
	}
	user, err := db.AddUser(ctx, du)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, uint(1), user.ID)

	nuser, err := db.FindUser(ctx, domain.LoginCredentials{Username: user.Username, Password: user.Password})
	assert.NoError(t, err)
	assert.NotNil(t, nuser)
	assert.Equal(t, user, nuser)
}
