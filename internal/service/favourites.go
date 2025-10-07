package service

import "context"

type AssetsRepository interface {
	GetByID(ctx context.Context, id int64) (Asset, error)
	Create(ctx context.Context, a Asset) (Asset, error)
}

type FavouritesRepository interface {
	ListByUser(ctx context.Context, userID int64, limit, offset int) ([]Favourite, int, error)
	Add(ctx context.Context, userID int64, asset Asset, customDesc *string) (Favourite, error)
	Remove(ctx context.Context, userID, assetID int64) error
	UpdateDescription(ctx context.Context, userID, assetID int64, desc string) error
	ClearAll(ctx context.Context, userID int64) error
}

type FavouritesService struct {
	Favs   FavouritesRepository
	Assets AssetsRepository
}

func (s *FavouritesService) List(ctx context.Context, userID int64, limit, offset int) ([]Favourite, int, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return s.Favs.ListByUser(ctx, userID, limit, offset)
}
func (s *FavouritesService) Add(ctx context.Context, userID int64, a Asset, custom *string) (Favourite, error) {
	return s.Favs.Add(ctx, userID, a, custom)
}
func (s *FavouritesService) Remove(ctx context.Context, userID, assetID int64) error {
	return s.Favs.Remove(ctx, userID, assetID)
}
func (s *FavouritesService) EditDescription(ctx context.Context, userID, assetID int64, desc string) error {
	return s.Favs.UpdateDescription(ctx, userID, assetID, desc)
}

func (s *FavouritesService) ClearAll(ctx context.Context, userID int64) error {
	return s.Favs.ClearAll(ctx, userID)
}
