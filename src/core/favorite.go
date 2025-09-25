package core

import "time"

// Cursor edge
type FavoriteEdge struct {
	AssetID   AssetID   `json:"assetId"`
	CreatedAt time.Time `json:"createdAt"`
}

type FavoriteService interface {
	AddToFavorites(userID UserID, assetID AssetID, now time.Time) error
	RemoveFromFavorites(userID UserID, assetID AssetID) error
	ListFavorites(userID UserID, limit int, after *FavoriteEdge, typ *AssetType) (items []AnyAsset, next *FavoriteEdge, err error)
}

type favoriteService struct {
	assets    AssetRepository
	favorites FavoriteRepository
}

func NewFavoriteService(assets AssetRepository, favorites FavoriteRepository) FavoriteService {
	return &favoriteService{assets: assets, favorites: favorites}
}

func (s *favoriteService) AddToFavorites(userID UserID, assetID AssetID, now time.Time) error {
	exists, err := s.assets.Exists(assetID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrInvalidID
	}
	return s.favorites.Add(userID, assetID, now)
}

func (s *favoriteService) RemoveFromFavorites(userID UserID, assetID AssetID) error {
	return s.favorites.Remove(userID, assetID)
}

func (s *favoriteService) ListFavorites(userID UserID, limit int, after *FavoriteEdge, typ *AssetType) ([]AnyAsset, *FavoriteEdge, error) {
	edges, err := s.favorites.ListByCursor(userID, limit, after, typ)
	if err != nil {
		return nil, nil, err
	}
	ids := make([]AssetID, 0, len(edges))
	for _, e := range edges {
		ids = append(ids, e.AssetID)
	}
	items, err := s.assets.ListByIDsPreserveOrder(ids)
	if err != nil {
		return nil, nil, err
	}
	var next *FavoriteEdge
	if len(edges) == limit {
		last := edges[len(edges)-1]
		next = &FavoriteEdge{AssetID: last.AssetID, CreatedAt: last.CreatedAt}
	}
	return items, next, nil
}

type AssetService interface {
	EditDescription(assetID AssetID, newDesc string, now time.Time) error
}

type assetService struct{ assets AssetRepository }

func NewAssetService(assets AssetRepository) AssetService { return &assetService{assets: assets} }

func (s *assetService) EditDescription(assetID AssetID, newDesc string, now time.Time) error {
	return s.assets.UpdateDescription(assetID, newDesc, now)
}
