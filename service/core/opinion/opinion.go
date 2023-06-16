package opinion

import (
	"context"
	"fmt"

	"x-gwi/app/storage"
	store_pb "x-gwi/proto/core/_store/v1"
	opinon_pb "x-gwi/proto/core/opinion/v1"
	"x-gwi/service"
)

//nolint:unused
type CoreOpinion struct {
	opinon   *opinon_pb.OpinionAsset
	idx      *store_pb.StoreIDX
	storage  *storage.ServiceStorage
	coreName service.CoreName
}

func NewCore(storage *storage.ServiceStorage) (*CoreOpinion, error) {
	c := &CoreOpinion{ //nolint:exhaustruct
		coreName: service.NameOpinion,
		storage:  storage,
	}

	if c.storage.CoreName() != c.coreName {
		return nil, fmt.Errorf("wrong storage coreName") //nolint:goerr113
	}

	return c, nil
}

//nolint:nilnil,revive
func (c *CoreOpinion) Create(ctx context.Context, in *opinon_pb.OpinionAsset) (*store_pb.StoreIDX, error) {
	return nil, nil
}
