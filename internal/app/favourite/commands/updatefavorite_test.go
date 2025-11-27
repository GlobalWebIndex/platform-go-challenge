package commands_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/akazantzidis/gwi-ass/internal/app/favourite/commands"
	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateFavoriteRequestHandler_Handle(t *testing.T) {
	mockUserID := uuid.New()
	mockFavoriteID := uuid.New()

	existingFavorite := &favourite.Favorite{
		ID:          mockFavoriteID,
		Description: "Old Description",
		Type:        "OldType",
		Data:        json.RawMessage(`{"x":0}`),
	}

	updatedFavorite := favourite.Favorite{
		ID:          mockFavoriteID,
		Description: "New Description",
		Type:        "NewType",
		Data:        json.RawMessage(`{"x":1}`),
	}

	tests := []struct {
		name          string
		setupMock     func(m *MockRepositoryF)
		command       commands.UpdateFavoriteRequest
		expectedError string
	}{
		{
			name: "happy path - update succeeds",
			setupMock: func(m *MockRepositoryF) {
				m.On("GetByID", mockUserID, mockFavoriteID).Return(existingFavorite, nil)
				m.On("Update", mockUserID, updatedFavorite).Return(nil)
			},
			command: commands.UpdateFavoriteRequest{
				UserID:      mockUserID,
				ID:          mockFavoriteID,
				Type:        "NewType",
				Description: "New Description",
				Data:        json.RawMessage(`{"x":1}`),
			},
			expectedError: "",
		},
		{
			name: "favorite not found",
			setupMock: func(m *MockRepositoryF) {
				m.On("GetByID", mockUserID, mockFavoriteID).Return(nil, nil)
			},
			command: commands.UpdateFavoriteRequest{
				UserID:      mockUserID,
				ID:          mockFavoriteID,
				Type:        "NewType",
				Description: "New Description",
				Data:        json.RawMessage(`{"x":1}`),
			},
			expectedError: "favorite with ID",
		},
		{
			name: "GetByID returns error",
			setupMock: func(m *MockRepositoryF) {
				m.On("GetByID", mockUserID, mockFavoriteID).Return(nil, errors.New("repo error"))
			},
			command: commands.UpdateFavoriteRequest{
				UserID:      mockUserID,
				ID:          mockFavoriteID,
				Type:        "NewType",
				Description: "New Description",
				Data:        json.RawMessage(`{"x":1}`),
			},
			expectedError: "failed to fetch favorite",
		},
		{
			name: "Update returns error",
			setupMock: func(m *MockRepositoryF) {
				m.On("GetByID", mockUserID, mockFavoriteID).Return(existingFavorite, nil)
				m.On("Update", mockUserID, updatedFavorite).Return(errors.New("update failed"))
			},
			command: commands.UpdateFavoriteRequest{
				UserID:      mockUserID,
				ID:          mockFavoriteID,
				Type:        "NewType",
				Description: "New Description",
				Data:        json.RawMessage(`{"x":1}`),
			},
			expectedError: "failed to update favorite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepositoryF{}
			tt.setupMock(mockRepo)

			handler := commands.NewUpdateFavoriteRequestHandler(mockRepo)
			err := handler.Handle(tt.command)

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
