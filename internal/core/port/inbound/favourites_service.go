package inbound

import (
	"github.com/Kercyn/crud_template/internal/core/domain"
	port "github.com/Kercyn/crud_template/internal/core/port/outbound"
)

type FavouritesService interface {
	GetUserFavourites(
		userID port.DataSourceID,
	) ([]domain.Asset, error)

	MarkAssetAsUserFavourite(
		userID port.DataSourceID,
		assetID port.DataSourceID,
	) error

	UnmarkAssetAsUserFavourite(
		userID port.DataSourceID,
		assetID port.DataSourceID,
	) error

	UpdateAssetDescription(
		userID port.DataSourceID,
		assetID port.DataSourceID,
		newDescription string,
	) error
}
