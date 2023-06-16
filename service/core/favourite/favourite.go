package favourite

import (
	"context"
	"fmt"

	"x-gwi/app/storage"
	store_pb "x-gwi/proto/core/_store/v1"
	favourite_pb "x-gwi/proto/core/favourite/v1"
	"x-gwi/service"
)

//nolint:unused
type CoreFavourite struct {
	favourite *favourite_pb.FavouriteAsset
	idx       *store_pb.StoreIDX
	storage   *storage.ServiceStorage
	coreName  service.CoreName
}

func NewCore(storage *storage.ServiceStorage) (*CoreFavourite, error) {
	c := &CoreFavourite{ //nolint:exhaustruct
		coreName: service.NameFavourite,
		storage:  storage,
	}

	if c.storage.CoreName() != c.coreName {
		return nil, fmt.Errorf("wrong storage coreName") //nolint:goerr113
	}

	return c, nil
}

//nolint:nilnil,revive
func (c *CoreFavourite) Create(ctx context.Context, in *favourite_pb.FavouriteAsset) (*store_pb.StoreIDX, error) {
	return nil, nil
}
