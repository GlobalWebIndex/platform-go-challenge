package asset

import (
	"context"
	"fmt"

	"x-gwi/app/storage"
	"x-gwi/app/x/id"
	storepb "x-gwi/proto/core/_store/v1"
	assetpb "x-gwi/proto/core/asset/v1"
	"x-gwi/service"
)

type CoreAsset struct {
	storage  *storage.ServiceStorage
	coreName service.CoreName
}

func NewCore(storage *storage.ServiceStorage) (*CoreAsset, error) {
	c := &CoreAsset{
		coreName: service.NameAsset,
		storage:  storage,
	}

	if c.storage.CoreName() != c.coreName {
		return nil, fmt.Errorf("wrong storage coreName") //nolint:goerr113
	}

	return c, nil
}

func (c *CoreAsset) Create(ctx context.Context, in *assetpb.AssetCore) error {
	in.Qid.Kind = c.storage.CoreName().String()
	// in.Qid.Key = in.Qid.Key // use directly
	in.Qid.Uid = id.XiD().String()
	in.Qid.Uuid = id.UUID().String()

	// c.storage.IsAQL()
	//nolint:exhaustruct
	dAQL := &storepb.StoreAQL{
		XKey:  in.Qid.Key,
		Asset: in,
	}

	m, err := c.storage.AQL().CreateDocument(ctx, dAQL, nil)
	if err != nil {
		return fmt.Errorf("AQL().CreateDocument: %w", err)
	}

	in.Qid.Rev = m.Rev

	return nil
}

func (c *CoreAsset) Get(ctx context.Context, in *assetpb.AssetCore) error {
	//nolint:exhaustruct
	dAQL := &storepb.StoreAQL{
		Asset: in,
	}

	m, err := c.storage.AQL().ReadDocument(ctx, in.Qid.Key, dAQL)
	if err != nil {
		return fmt.Errorf("AQL().ReadDocument: %w", err)
	}

	in.Qid.Rev = m.Rev

	return nil
}

func (c *CoreAsset) Update(ctx context.Context, in *assetpb.AssetCore) error {
	// c.storage.IsAQL()
	//nolint:exhaustruct
	dAQL := &storepb.StoreAQL{
		Asset: in,
	}

	m, err := c.storage.AQL().UpdateDocument(ctx, in.Qid.Key, in.Qid.Rev, dAQL, dAQL, nil)
	if err != nil {
		return fmt.Errorf("AQL().CreateDocument: %w", err)
	}

	// if m.Key != in.Qid.Key {todo delete wronk key}
	in.Qid.Key = m.Key
	in.Qid.Rev = m.Rev

	return nil
}
