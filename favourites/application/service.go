package application

import (
	"github.com/gwi/platform-go-challenge/favourites/domain"
)

type FavouritesService struct {
	repo         FavouritesRepository
	assetCatalog AssetCatalogClient // optional: nil means no validation of source_asset_id
}

// NewFavouritesService creates the use case. assetCatalog may be nil to skip catalog lookup and add-by-reference.
func NewFavouritesService(repo FavouritesRepository, assetCatalog AssetCatalogClient) *FavouritesService {
	return &FavouritesService{repo: repo, assetCatalog: assetCatalog}
}

var _ FavouritesUseCase = (*FavouritesService)(nil)

func (s *FavouritesService) ListFavouritesPaginated(userID string, limit, offset int) ([]domain.Asset, int) {
	return s.repo.ListPaginated(userID, limit, offset)
}

func (s *FavouritesService) GetFavourite(userID string, assetID string) (*domain.Asset, bool) {
	return s.repo.Get(userID, assetID)
}

func (s *FavouritesService) AddFavourite(userID string, asset domain.Asset) (domain.Asset, error) {
	requestDescription := asset.Description
	if s.assetCatalog != nil && asset.SourceAssetID != "" && asset.Type != "" {
		catalogAsset, err := s.assetCatalog.GetAsset(userID, asset.SourceAssetID, asset.Type)
		if err != nil {
			return domain.Asset{}, err
		}
		if catalogAsset == nil {
			return domain.Asset{}, ErrAssetNotFound
		}
		// Add-by-reference: use catalog data as base; allow request to override description.
		asset = *catalogAsset
		asset.ID = ""
		if requestDescription != "" {
			asset.Description = requestDescription
		}
	}
	return s.repo.Add(userID, asset), nil
}

func (s *FavouritesService) RemoveFavourite(userID string, assetID string) bool {
	return s.repo.Remove(userID, assetID)
}

func (s *FavouritesService) UpdateFavouriteDescription(userID string, assetID string, description string) (*domain.Asset, bool) {
	return s.repo.UpdateDescription(userID, assetID, description)
}

func (s *FavouritesService) OnAssetDeleted(sourceAssetID string) int {
	return s.repo.RemoveFavouritesBySourceAssetID(sourceAssetID)
}
