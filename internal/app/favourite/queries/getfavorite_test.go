package queries_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/akazantzidis/gwi-ass/internal/app/favourite/queries"
	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetFavoriteRequestHandler_Handle(t *testing.T) {
	mockUserID := uuid.New()
	mockFavoriteID := uuid.New()
	mockTime := time.Now().UTC()

	tests := []struct {
		name          string
		mockReturn    *favourite.Favorite
		mockError     error
		expectedError string
	}{
		{
			name: "happy path - favorite exists",
			mockReturn: &favourite.Favorite{
				ID:          mockFavoriteID,
				Type:        "Chart",
				Description: "Favorite Chart",
				Data:        json.RawMessage(`{"x":1}`),
				CreatedAt:   mockTime,
			},
			mockError:     nil,
			expectedError: "",
		},
		{
			name:          "favorite not found",
			mockReturn:    nil,
			mockError:     nil,
			expectedError: "not found",
		},
		{
			name:          "repo returns error",
			mockReturn:    nil,
			mockError:     errors.New("repo failure"),
			expectedError: "failed to fetch favorite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepositoryF{}
			mockRepo.On("GetByID", mockUserID, mockFavoriteID).Return(tt.mockReturn, tt.mockError)

			handler := queries.NewGetFavoriteRequestHandler(mockRepo)

			result, err := handler.Handle(queries.GetFavoriteRequest{
				UserID:     mockUserID,
				FavoriteID: mockFavoriteID,
			})

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockReturn.ID, result.ID)
				assert.Equal(t, tt.mockReturn.Type, result.Type)
				assert.Equal(t, tt.mockReturn.Description, result.Description)
				assert.JSONEq(t, string(tt.mockReturn.Data), string(result.Data))
				assert.Equal(t, tt.mockReturn.CreatedAt, result.CreatedAt)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
