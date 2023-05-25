package service

import (
	"github.com/Kercyn/crud_template/internal/core/domain"
	port "github.com/Kercyn/crud_template/internal/core/port/outbound"
)

type FavouritesServiceImpl struct {
	userRepository port.UserRepository
	requestBuilder port.UserPatchRequestBuilder
}

func NewUserFavouritesServiceImpl(
	userRepository port.UserRepository,
	requestBuilder port.UserPatchRequestBuilder,
) *FavouritesServiceImpl {
	return &FavouritesServiceImpl{
		userRepository: userRepository,
		requestBuilder: requestBuilder,
	}
}

func (f FavouritesServiceImpl) GetUserFavourites(
	userID port.DataSourceID,
) ([]domain.Asset, error) {
	user, err := f.userRepository.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	return f.extractFavouritesFromAssets(user.Assets), nil
}

func (f FavouritesServiceImpl) extractFavouritesFromAssets(
	assets []domain.Asset,
) []domain.Asset {
	var favourites []domain.Asset
	for _, a := range assets {
		if a.GetIsFavourite() {
			favourites = append(favourites, a)
		}
	}

	return favourites
}

func (f FavouritesServiceImpl) MarkAssetAsUserFavourite(
	userID port.DataSourceID,
	assetID port.DataSourceID,
) error {
	f.setAssetIsFavourite(userID, assetID, true)
	request := f.requestBuilder.Build()
	return f.userRepository.Patch(request)
}

func (f FavouritesServiceImpl) UnmarkAssetAsUserFavourite(
	userID port.DataSourceID,
	assetID port.DataSourceID,
) error {
	f.setAssetIsFavourite(userID, assetID, false)
	request := f.requestBuilder.Build()
	return f.userRepository.Patch(request)
}

func (f FavouritesServiceImpl) setAssetIsFavourite(
	userID port.DataSourceID,
	assetID port.DataSourceID,
	isFavourite bool,
) {
	f.requestBuilder.Reset()
	f.requestBuilder.WithUserID(userID)
	f.requestBuilder.SetAssetIsFavourite(assetID, isFavourite)
}

func (f FavouritesServiceImpl) UpdateAssetDescription(
	userID port.DataSourceID,
	assetID port.DataSourceID,
	description string,
) error {
	f.requestBuilder.Reset()
	f.requestBuilder.WithUserID(userID)
	f.requestBuilder.SetAssetDescription(assetID, description)
	request := f.requestBuilder.Build()
	return f.userRepository.Patch(request)
}
