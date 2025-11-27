package commands_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/akazantzidis/gwi-ass/internal/app/favourite/commands"
	"github.com/akazantzidis/gwi-ass/internal/domain/favourite"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdatePartialFavoriteRequestHandler_HandlePartial(t *testing.T) {
	mockUserID := uuid.New()
	mockFavoriteID := uuid.New()

	oldType := favourite.AssetType("OldType")
	oldDescription := "OldDescription"
	oldData := json.RawMessage(`{"x":0}`)

	tests := []struct {
		name          string
		req           commands.PatchFavoriteRequest
		setupMock     func(m *MockRepositoryF, fav *favourite.Favorite)
		expectedFav   *favourite.Favorite
		expectedError string
	}{
		{
			name: "happy path - all fields updated",
			req: commands.PatchFavoriteRequest{
				Type:        ptrAssetType("NewType"),
				Description: ptrString("NewDescription"),
				Data:        ptrRawMessage(`{"x":1}`),
			},
			setupMock: func(m *MockRepositoryF, fav *favourite.Favorite) {
				m.On("GetByID", mock.Anything, mock.Anything).Return(fav, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedFav: &favourite.Favorite{
				ID:          mockFavoriteID,
				Type:        "NewType",
				Description: "NewDescription",
				Data:        json.RawMessage(`{"x":1}`),
			},
			expectedError: "",
		},
		{
			name: "partial update - only description",
			req: commands.PatchFavoriteRequest{
				Description: ptrString("UpdatedDescription"),
			},
			setupMock: func(m *MockRepositoryF, fav *favourite.Favorite) {
				m.On("GetByID", mock.Anything, mock.Anything).Return(fav, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedFav: &favourite.Favorite{
				ID:          mockFavoriteID,
				Type:        oldType,
				Description: "UpdatedDescription",
				Data:        oldData,
			},
			expectedError: "",
		},
		{
			name: "favorite not found",
			req:  commands.PatchFavoriteRequest{Description: ptrString("Anything")},
			setupMock: func(m *MockRepositoryF, fav *favourite.Favorite) {
				m.On("GetByID", mock.Anything, mock.Anything).Return(nil, nil)
			},
			expectedFav:   nil,
			expectedError: "favorite with ID",
		},
		{
			name: "GetByID returns error",
			req:  commands.PatchFavoriteRequest{Description: ptrString("Anything")},
			setupMock: func(m *MockRepositoryF, fav *favourite.Favorite) {
				m.On("GetByID", mock.Anything, mock.Anything).Return(nil, errors.New("repo error"))
			},
			expectedFav:   nil,
			expectedError: "failed to fetch favorite",
		},
		{
			name: "Update returns error",
			req: commands.PatchFavoriteRequest{
				Description: ptrString("NewDesc"),
			},
			setupMock: func(m *MockRepositoryF, fav *favourite.Favorite) {
				m.On("GetByID", mock.Anything, mock.Anything).Return(fav, nil)
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("update failed"))
			},
			expectedFav:   nil,
			expectedError: "failed to update favorite",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepositoryF{}

			// Create fresh favorite instance for this subtest
			fav := &favourite.Favorite{
				ID:          mockFavoriteID,
				Type:        oldType,
				Description: oldDescription,
				Data:        oldData,
			}

			tt.setupMock(mockRepo, fav)

			handler := commands.NewUpdatePartialFavoriteRequestHandler(mockRepo)
			result, err := handler.HandlePartial(mockUserID, mockFavoriteID, tt.req)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedFav.ID, result.ID)
				assert.Equal(t, tt.expectedFav.Type, result.Type)
				assert.Equal(t, tt.expectedFav.Description, result.Description)
				assert.JSONEq(t, string(tt.expectedFav.Data), string(result.Data))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// helper functions to create pointers
func ptrAssetType(s string) *favourite.AssetType { t := favourite.AssetType(s); return &t }
func ptrString(s string) *string                 { return &s }
func ptrRawMessage(s string) *json.RawMessage {
	r := json.RawMessage(s)
	return &r
}
