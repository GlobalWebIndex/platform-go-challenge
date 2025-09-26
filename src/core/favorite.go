package core

import "time"

// Cursor edge + per-user override
type FavoriteEdge struct {
	AssetID   AssetID            `json:"assetId"`
	CreatedAt time.Time          `json:"createdAt"`
	Override  map[string]any     `json:"-"` // internal; merged into asset before returning
}

type FavoriteService interface {
	AddToFavorites(userID UserID, assetID AssetID, now time.Time) error
	RemoveFromFavorites(userID UserID, assetID AssetID) error
	ListFavorites(userID UserID, limit int, after *FavoriteEdge, typ *AssetType) (items []AnyAsset, next *FavoriteEdge, err error)

	// NEW: per-user override (e.g., description, title)
	SetFavoriteOverride(userID UserID, assetID AssetID, override map[string]any, now time.Time) error
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
	if err != nil { return err }
	if !exists { return ErrInvalidID }
	return s.favorites.Add(userID, assetID, now)
}

func (s *favoriteService) RemoveFromFavorites(userID UserID, assetID AssetID) error {
	return s.favorites.Remove(userID, assetID)
}

func (s *favoriteService) SetFavoriteOverride(userID UserID, assetID AssetID, override map[string]any, now time.Time) error {
	// Optionally: validate fields (only allow known keys)
	// Keep it permissive now; HTTP does whitelisting
	return s.favorites.SetOverride(userID, assetID, override, now)
}

func (s *favoriteService) ListFavorites(userID UserID, limit int, after *FavoriteEdge, typ *AssetType) ([]AnyAsset, *FavoriteEdge, error) {
	edges, err := s.favorites.ListByCursor(userID, limit, after, typ)
	if err != nil { return nil, nil, err }

	ids := make([]AssetID, 0, len(edges))
	for _, e := range edges { ids = append(ids, e.AssetID) }

	items, err := s.assets.ListByIDsPreserveOrder(ids)
	if err != nil { return nil, nil, err }

	// Merge per-user overrides onto returned assets (only shared/base fields)
	for i := range items {
		base := items[i].GetAssetBase()
		ov := edges[i].Override
		if ov == nil { continue }
		if v, ok := ov["title"].(string); ok && v != "" {
			base.Title = v
		}
		if v, ok := ov["description"].(string); ok && v != "" {
			base.Description = v
		}
		// You can add more shared fields here if needed later
	}

	var next *FavoriteEdge
	if len(edges) == limit {
		last := edges[len(edges)-1]
		next = &FavoriteEdge{AssetID: last.AssetID, CreatedAt: last.CreatedAt}
	}
	return items, next, nil
}
