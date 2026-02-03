package service

import (
	"context"
	"errors"
	"testing"

	"go-favorites-app/internal/core/domain/auth"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Save(ctx context.Context, user auth.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (auth.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(auth.User), args.Error(1)
}

func TestAuthService_SignUp(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := NewAuthService(mockRepo, "secret")

	t.Run("success", func(t *testing.T) {
		email := "test@example.com"
		password := "password123"

		mockRepo.On("Save", mock.Anything, mock.MatchedBy(func(u auth.User) bool {
			return u.Email == email && u.ID != "" && u.PasswordHash != ""
		})).Return(nil)

		err := svc.SignUp(context.Background(), email, password)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		svc := NewAuthService(mockRepo, "secret")
		mockRepo.On("Save", mock.Anything, mock.Anything).Return(errors.New("db error"))

		err := svc.SignUp(context.Background(), "test@example.com", "pass")
		assert.Error(t, err)
	})
}

func TestAuthService_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	svc := NewAuthService(mockRepo, "mysecret")

	password := "password123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := auth.User{ID: "user1", Email: "test@example.com", PasswordHash: string(hashed)}

	t.Run("success", func(t *testing.T) {
		mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)

		token, err := svc.Login(context.Background(), "test@example.com", password)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Verify token
		parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return []byte("mysecret"), nil
		})
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		assert.True(t, ok)
		assert.Equal(t, "user1", claims["sub"])
	})

	t.Run("invalid credentials - wrong password", func(t *testing.T) {
		// Expect FindByEmail but validation fails after
		mockRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(user, nil)

		token, err := svc.Login(context.Background(), "test@example.com", "wrongpass")
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		assert.Empty(t, token)
	})

	t.Run("invalid credentials - user not found", func(t *testing.T) {
		mockRepo.On("FindByEmail", mock.Anything, "unknown@example.com").Return(auth.User{}, errors.New("not found"))

		token, err := svc.Login(context.Background(), "unknown@example.com", "pass")
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
		assert.Empty(t, token)
	})
}
