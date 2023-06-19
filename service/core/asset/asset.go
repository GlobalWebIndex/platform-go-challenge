package asset

import (
	"context"
	"fmt"

	"x-gwi/app/storage"
	"x-gwi/app/storage/storepb2"
	"x-gwi/app/x/id"
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

	// if c.storage.IsAQL()
	// "x-gwi/app/storage/storepb2"
	//nolint:exhaustruct
	dAQL := &storepb2.StoreAQL{
		Key:   in.Qid.Key,
		Qid:   in.Qid,
		Asset: in,
	}

	m, err := c.storage.AQL().CreateDocument(ctx, dAQL)
	if err != nil {
		return fmt.Errorf("AQL().CreateDocument: %w", err)
	}

	// if m.Key != in.Qid.Key {todo delete wronk key}
	in.Qid.Key = m.Key
	in.Qid.Rev = m.Rev

	return nil
}

func (c *CoreAsset) Get(ctx context.Context, in *assetpb.AssetCore) error {
	// if c.storage.IsAQL()
	// "x-gwi/app/storage/storepb2"
	//nolint:exhaustruct
	dAQL := &storepb2.StoreAQL{
		// Key:  in.Qid.Key,
		// Qid:  in.Qid,
		Asset: in,
	}

	m, err := c.storage.AQL().ReadDocument(ctx, in.Qid.Key, dAQL)
	if err != nil {
		return fmt.Errorf("AQL().ReadDocument: %w", err)
	}

	// in.Qid.Key = m.Key
	in.Qid.Rev = m.Rev

	return nil
}
