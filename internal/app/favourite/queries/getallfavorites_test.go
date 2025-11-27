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
	"github.com/stretchr/testify/mock"
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

func TestGetAllFavoritesRequestHandler_Handle(t *testing.T) {
	mockUserID := uuid.New()
	mockFavoriteID := uuid.New()
	mockTime := time.Now().UTC()

	tests := []struct {
		name          string
		mockReturn    []favourite.Favorite
		mockError     error
		expectedError string
		expectedCount int
	}{
		{
			name: "happy path - multiple favorites",
			mockReturn: []favourite.Favorite{
				{
					ID:          mockFavoriteID,
					Type:        "Chart",
					Description: "Favorite Chart",
					Data:        json.RawMessage(`{"x":1}`),
					CreatedAt:   mockTime,
					UpdatedAt:   mockTime,
				},
			},
			mockError:     nil,
			expectedError: "",
			expectedCount: 1,
		},
		{
			name:          "no favorites",
			mockReturn:    []favourite.Favorite{},
			mockError:     nil,
			expectedError: "",
			expectedCount: 0,
		},
		{
			name:          "repo returns error",
			mockReturn:    nil,
			mockError:     errors.New("repo failure"),
			expectedError: "failed to fetch favorites",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepositoryF{}
			mockRepo.On("GetAll", mockUserID).Return(tt.mockReturn, tt.mockError)

			handler := queries.NewGetAllFavoritesRequestHandler(mockRepo)

			result, err := handler.Handle(queries.GetAllFavoritesRequest{
				UserID: mockUserID,
			})

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
				for i, r := range result {
					assert.Equal(t, tt.mockReturn[i].ID, r.ID)
					assert.Equal(t, tt.mockReturn[i].Type, r.Type)
					assert.Equal(t, tt.mockReturn[i].Description, r.Description)
					assert.JSONEq(t, string(tt.mockReturn[i].Data), string(r.Data))
					assert.Equal(t, tt.mockReturn[i].CreatedAt, r.CreatedAt)
					assert.Equal(t, tt.mockReturn[i].UpdatedAt, r.UpdatedAt)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
