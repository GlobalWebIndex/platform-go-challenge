package asset

import (
	"context"
	"fmt"

	"x-gwi/app/storage"
	store_pb "x-gwi/proto/core/_store/v1"
	asset_pb "x-gwi/proto/core/asset/v1"
	asset_srvpb "x-gwi/proto/serv/asset/v1"
	"x-gwi/service"
)

//nolint:unused
type CoreAsset struct {
	asset_srvpb.UnimplementedAssetServiceServer
	idx      *store_pb.StoreIDX
	storage  *storage.ServiceStorage
	asset    *asset_pb.AssetInstance
	coreName service.CoreName
}

func NewCore(storage *storage.ServiceStorage) (*CoreAsset, error) {
	c := &CoreAsset{ //nolint:exhaustruct
		coreName: service.NameAsset,
		storage:  storage,
	}

	if c.storage.CoreName() != c.coreName {
		return nil, fmt.Errorf("wrong storage coreName") //nolint:goerr113
	}

	return c, nil
}

//nolint:nilnil,revive
func (c *CoreAsset) Create(ctx context.Context, in *asset_pb.AssetInstance) (*store_pb.StoreIDX, error) {
	return nil, nil
}
