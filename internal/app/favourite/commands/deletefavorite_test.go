package commands_test

import (
	"errors"
	"testing"

	"github.com/akazantzidis/gwi-ass/internal/app/favourite/commands"
	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteFavoriteRequestHandler_Handle(t *testing.T) {
	mockUserID := uuid.New()
	mockFavoriteID := uuid.New()
	existingFavorite := &favourite.Favorite{
		ID:          mockFavoriteID,
		Description: "Test Favorite",
	}

	tests := []struct {
		name          string
		setupMock     func(m *MockRepositoryF)
		expectedError string
	}{
		{
			name: "happy path - delete succeeds",
			setupMock: func(m *MockRepositoryF) {
				m.On("GetByID", mockUserID, mockFavoriteID).Return(existingFavorite, nil)
				m.On("Delete", mockUserID, mockFavoriteID).Return(nil)
			},
			expectedError: "",
		},
		{
			name: "favorite does not exist",
			setupMock: func(m *MockRepositoryF) {
				m.On("GetByID", mockUserID, mockFavoriteID).Return(nil, nil)
			},
			expectedError: "favorite with ID",
		},
		{
			name: "GetByID returns error",
			setupMock: func(m *MockRepositoryF) {
				m.On("GetByID", mockUserID, mockFavoriteID).Return(nil, errors.New("repo error"))
			},
			expectedError: "failed to check favorite existence",
		},
		{
			name: "Delete returns error",
			setupMock: func(m *MockRepositoryF) {
				m.On("GetByID", mockUserID, mockFavoriteID).Return(existingFavorite, nil)
				m.On("Delete", mockUserID, mockFavoriteID).Return(errors.New("delete failed"))
			},
			expectedError: "failed to delete favorite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepositoryF{}
			tt.setupMock(mockRepo)

			handler := commands.NewDeleteFavoriteRequestHandler(mockRepo)

			err := handler.Handle(commands.DeleteFavoriteRequest{
				UserID:     mockUserID,
				FavoriteID: mockFavoriteID,
			})

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
