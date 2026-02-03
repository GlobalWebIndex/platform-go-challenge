package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid user",
			user: User{
				ID:    "123",
				Email: "test@example.com",
			},
			wantErr: false,
		},
		{
			name: "missing id",
			user: User{
				Email: "test@example.com",
			},
			wantErr: true,
			errMsg:  "id is required",
		},
		{
			name: "missing email",
			user: User{
				ID: "123",
			},
			wantErr: true,
			errMsg:  "email is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
