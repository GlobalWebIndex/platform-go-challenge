package commands_test

import (
	"encoding/json"
	"errors"
	"testing"
	_ "time"

	"github.com/akazantzidis/gwi-ass/internal/app/notification"
	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/akazantzidis/gwi-ass/internal/app/favourite/commands"
)

// Mock repository for favorites
type MockRepositoryF struct {
	mock.Mock
}

func (m *MockRepositoryF) Add(userID uuid.UUID, fav favourite.Favorite) error {
	args := m.Called(userID, fav)
	return args.Error(0)
}

func (m *MockRepositoryF) GetByID(userID uuid.UUID, favoriteID uuid.UUID) (*favourite.Favorite, error) {
	args := m.Called(userID, favoriteID)
	fav := args.Get(0)
	if fav == nil {
		return nil, args.Error(1)
	}
	return fav.(*favourite.Favorite), args.Error(1)
}

func (m *MockRepositoryF) GetAll(userID uuid.UUID) ([]favourite.Favorite, error) {
	args := m.Called(userID)
	favs := args.Get(0)
	if favs == nil {
		return nil, args.Error(1)
	}
	return favs.([]favourite.Favorite), args.Error(1)
}

func (m *MockRepositoryF) Update(userID uuid.UUID, favorite favourite.Favorite) error {
	args := m.Called(userID, favorite)
	return args.Error(0)
}

// Delete mock implementation
func (m *MockRepositoryF) Delete(userID uuid.UUID, favoriteID uuid.UUID) error {
	args := m.Called(userID, favoriteID)
	return args.Error(0)
}

// Mock notification service
type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) Notify(n notification.Notification) error {
	args := m.Called(n)
	return args.Error(0)
}

func TestAddFavoriteRequestHandler_Handle(t *testing.T) {
	mockUserID := uuid.New()
	tests := []struct {
		name              string
		repoError         error
		notificationError error
		expectedError     string
	}{
		{
			name:              "happy path",
			repoError:         nil,
			notificationError: nil,
			expectedError:     "",
		},
		{
			name:              "repo returns error",
			repoError:         errors.New("repo failed"),
			notificationError: nil,
			expectedError:     "failed to add favorite: repo failed",
		},
		{
			name:              "notification returns error",
			repoError:         nil,
			notificationError: errors.New("notify failed"),
			expectedError:     "favorite added but failed to send notification: notify failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepositoryF{}
			mockNotification := &MockNotificationService{}

			// Repo Add always expected
			mockRepo.On("Add", mock.Anything, mock.Anything).Return(tt.repoError)

			// Notify may be called, depending on repo result
			mockNotification.On("Notify", mock.Anything).Maybe().Return(tt.notificationError)

			handler := commands.NewAddFavoriteRequestHandler(mockRepo, mockNotification)

			req := commands.AddFavoriteRequest{
				UserID:      mockUserID,
				Type:        "Chart",
				Description: "My Favorite Chart",
				Data:        json.RawMessage(`{"x":1,"y":2}`),
			}

			err := handler.Handle(req)
			if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockNotification.AssertExpectations(t)
		})
	}
}
