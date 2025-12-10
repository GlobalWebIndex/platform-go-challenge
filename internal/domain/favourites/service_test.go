package favourites

import (
	"context"
	"errors"
	"testing"
	"time"
)

// stubRepository is a test stub for the Repository interface.
type stubRepository struct {
	// ListFavourites
	listResult []FavouriteEntry
	listErr    error
	lastList   struct {
		userID string
		filter ListFilter
	}

	// CreateFavourite
	createResult FavouriteEntry
	createErr    error
	lastCreate   struct {
		userID string
		asset  Asset
	}

	// RemoveFavourite
	removeErr  error
	lastRemove struct {
		userID  string
		assetID string
	}

	// UpdateDescriptionForUser
	updateResult FavouriteEntry
	updateErr    error
	lastUpdate   struct {
		userID      string
		assetID     string
		description string
	}
}

func (s *stubRepository) CreateFavourite(
	_ context.Context,
	userID string,
	asset Asset,
) (FavouriteEntry, error) {
	s.lastCreate.userID = userID
	s.lastCreate.asset = asset
	return s.createResult, s.createErr
}

func (s *stubRepository) ListFavourites(
	_ context.Context,
	userID string,
	filter ListFilter,
) ([]FavouriteEntry, error) {
	s.lastList.userID = userID
	s.lastList.filter = filter
	return s.listResult, s.listErr
}

func (s *stubRepository) RemoveFavourite(
	_ context.Context,
	userID, assetID string,
) error {
	s.lastRemove.userID = userID
	s.lastRemove.assetID = assetID
	return s.removeErr
}

func (s *stubRepository) UpdateDescriptionForUser(
	_ context.Context,
	userID, assetID, description string,
) (FavouriteEntry, error) {
	s.lastUpdate.userID = userID
	s.lastUpdate.assetID = assetID
	s.lastUpdate.description = description
	return s.updateResult, s.updateErr
}

func TestNewServicePanicsOnNilRepo(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic when repo is nil, got none")
		}
	}()

	_ = NewService(nil) // should panic
}

func TestListUserFavourites_Validation(t *testing.T) {
	t.Parallel()

	svc := NewService(&stubRepository{})
	ctx := context.Background()

	_, err := svc.ListUserFavourites(ctx, "", ListFilter{})
	if !errors.Is(err, ErrUserIDRequired) {
		t.Fatalf("expected ErrUserIDRequired, got %v", err)
	}
}

func TestListUserFavourites_DelegatesToRepository(t *testing.T) {
	t.Parallel()

	stub := &stubRepository{
		listResult: []FavouriteEntry{
			{
				Favourite: Favourite{
					UserID:  "user-1",
					AssetID: "asset-1",
				},
				Asset: Asset{
					ID:   "asset-1",
					Type: AssetTypeInsight,
				},
			},
		},
	}
	svc := NewService(stub)
	ctx := context.Background()

	filter := ListFilter{Type: nil}
	result, err := svc.ListUserFavourites(ctx, "user-1", filter)
	if err != nil {
		t.Fatalf("ListUserFavourites returned error: %v", err)
	}

	if got, want := stub.lastList.userID, "user-1"; got != want {
		t.Fatalf("ListUserFavourites called with userID=%q, want %q", got, want)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 favourite entry, got %d", len(result))
	}
}

func TestListUserFavourites_PropagatesRepositoryError(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("list error")
	stub := &stubRepository{
		listErr: repoErr,
	}
	svc := NewService(stub)
	ctx := context.Background()

	_, err := svc.ListUserFavourites(ctx, "user-1", ListFilter{})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error %v, got %v", repoErr, err)
	}
}

func TestAddFavourite_Validation(t *testing.T) {
	t.Parallel()

	svc := NewService(&stubRepository{})
	ctx := context.Background()

	tests := []struct {
		name    string
		userID  string
		asset   Asset
		wantErr error
	}{
		{
			name:    "empty_userID",
			userID:  "",
			asset:   Asset{},
			wantErr: ErrUserIDRequired,
		},
		{
			name:   "asset_with_ID_is_rejected",
			userID: "user-1",
			asset: Asset{
				ID: "pre-set",
			},
			wantErr: ErrInvalidAssetPayload,
		},
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := svc.AddFavourite(ctx, tt.userID, tt.asset)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("AddFavourite error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestAddFavourite_PropagatesRepositoryError(t *testing.T) {
	t.Parallel()

	repoErr := errors.New("repo failure")
	stub := &stubRepository{
		createErr: repoErr,
	}
	svc := NewService(stub)
	ctx := context.Background()

	asset := Asset{
		Type: AssetTypeInsight,
	}

	_, err := svc.AddFavourite(ctx, "user-1", asset)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error %v, got %v", repoErr, err)
	}

	if stub.lastCreate.userID != "user-1" {
		t.Fatalf("CreateFavourite called with userID=%q, want %q", stub.lastCreate.userID, "user-1")
	}
}

func TestAddFavourite_SuccessDelegatesToRepository(t *testing.T) {
	t.Parallel()

	stub := &stubRepository{
		createResult: FavouriteEntry{
			Favourite: Favourite{
				UserID:  "user-1",
				AssetID: "asset-1",
			},
			Asset: Asset{
				ID:          "asset-1",
				Type:        AssetTypeInsight,
				Description: "desc",
			},
		},
	}
	svc := NewService(stub)
	ctx := context.Background()

	asset := Asset{ // ID must be empty; repo generates it
		Type: AssetTypeInsight,
	}

	entry, err := svc.AddFavourite(ctx, "user-1", asset)
	if err != nil {
		t.Fatalf("AddFavourite returned error: %v", err)
	}

	if got, want := stub.lastCreate.userID, "user-1"; got != want {
		t.Fatalf("CreateFavourite called with userID=%q, want %q", got, want)
	}
	if entry.Favourite.AssetID != "asset-1" {
		t.Fatalf("expected AssetID %q, got %q", "asset-1", entry.Favourite.AssetID)
	}
}

func TestRemoveFavourite_Validation(t *testing.T) {
	t.Parallel()

	svc := NewService(&stubRepository{})
	ctx := context.Background()

	tests := []struct {
		name    string
		userID  string
		assetID string
		wantErr error
	}{
		{
			name:    "empty_userID",
			userID:  "",
			assetID: "asset-1",
			wantErr: ErrUserIDRequired,
		},
		{
			name:    "empty_assetID",
			userID:  "user-1",
			assetID: "",
			wantErr: ErrAssetIDRequired,
		},
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := svc.RemoveFavourite(ctx, tt.userID, tt.assetID)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("RemoveFavourite error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestRemoveFavourite_DelegatesToRepository(t *testing.T) {
	t.Parallel()

	stub := &stubRepository{}
	svc := NewService(stub)
	ctx := context.Background()

	if err := svc.RemoveFavourite(ctx, "user-1", "asset-1"); err != nil {
		t.Fatalf("RemoveFavourite returned error: %v", err)
	}

	if got, want := stub.lastRemove.userID, "user-1"; got != want {
		t.Fatalf("RemoveFavourite called with userID=%q, want %q", got, want)
	}
	if got, want := stub.lastRemove.assetID, "asset-1"; got != want {
		t.Fatalf("RemoveFavourite called with assetID=%q, want %q", got, want)
	}
}

func TestEditFavouriteDescription_Validation(t *testing.T) {
	t.Parallel()

	svc := NewService(&stubRepository{})
	ctx := context.Background()

	tests := []struct {
		name        string
		userID      string
		assetID     string
		description string
		wantErr     error
	}{
		{
			name:        "empty_userID",
			userID:      "",
			assetID:     "asset-1",
			description: "desc",
			wantErr:     ErrUserIDRequired,
		},
		{
			name:        "empty_assetID",
			userID:      "user-1",
			assetID:     "",
			description: "desc",
			wantErr:     ErrAssetIDRequired,
		},
		{
			name:        "empty_description",
			userID:      "user-1",
			assetID:     "asset-1",
			description: "",
			wantErr:     ErrDescriptionRequired,
		},
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := svc.EditFavouriteDescription(ctx, tt.userID, tt.assetID, tt.description)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("EditFavouriteDescription error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestEditFavouriteDescription_DelegatesToRepositoryAndPropagatesErrors(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	wantAsset := Asset{
		ID:          "asset-1",
		Type:        AssetTypeChart,
		Description: "new-desc",
	}
	wantEntry := FavouriteEntry{
		Favourite: Favourite{
			UserID:    "user-1",
			AssetID:   "asset-1",
			CreatedAt: time.Now().UTC(),
		},
		Asset: wantAsset,
	}
	repoErr := errors.New("repo error")

	tests := []struct {
		name       string
		repo       *stubRepository
		wantErr    error
		expectCall bool
	}{
		{
			name: "success",
			repo: &stubRepository{
				updateResult: wantEntry,
				updateErr:    nil,
			},
			wantErr:    nil,
			expectCall: true,
		},
		{
			name: "propagates_repository_error",
			repo: &stubRepository{
				updateErr: repoErr,
			},
			wantErr:    repoErr,
			expectCall: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := NewService(tt.repo)

			got, err := svc.EditFavouriteDescription(ctx, "user-1", "asset-1", "new-desc")
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("EditFavouriteDescription error = %v, want %v", err, tt.wantErr)
			}

			if !tt.expectCall {
				return
			}

			if tt.wantErr == nil {
				if got.Asset.Description != wantAsset.Description {
					t.Fatalf("expected description %q, got %q", wantAsset.Description, got.Asset.Description)
				}
			}

			if tt.repo.lastUpdate.userID != "user-1" ||
				tt.repo.lastUpdate.assetID != "asset-1" ||
				tt.repo.lastUpdate.description != "new-desc" {
				t.Fatalf(
					"UpdateDescriptionForUser called with (%q, %q, %q), want (%q, %q, %q)",
					tt.repo.lastUpdate.userID,
					tt.repo.lastUpdate.assetID,
					tt.repo.lastUpdate.description,
					"user-1",
					"asset-1",
					"new-desc",
				)
			}
		})
	}
}
