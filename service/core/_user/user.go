package user

import (
	"context"
	"fmt"

	"x-gwi/app/storage"
	"x-gwi/app/x/id"
	storepb "x-gwi/proto/core/_store/v1"
	userpb "x-gwi/proto/core/_user/v1"
	"x-gwi/service"
)

type CoreUser struct {
	storage  *storage.ServiceStorage
	coreName service.CoreName
}

func NewCore(storage *storage.ServiceStorage) (*CoreUser, error) {
	c := &CoreUser{
		coreName: service.NameUser,
		storage:  storage,
	}

	if c.storage.CoreName() != c.coreName {
		return nil, fmt.Errorf("wrong storage coreName") //nolint:goerr113
	}

	return c, nil
}

func (c *CoreUser) Create(ctx context.Context, in *userpb.UserCore) error {
	// rev := id.Rev()
	// QID on create
	// in.Md.Label = &rev
	in.Qid.Kind = c.storage.CoreName().String()
	in.Qid.Key = in.BasicAccount.GetUsername()
	in.Qid.Uid = id.XiD().String()
	in.Qid.Uuid = id.UUID().String()

	// c.storage.IsAQL()
	//nolint:exhaustruct
	dAQL := &storepb.StoreAQL{
		XKey: in.Qid.Key,
		Qid:  in.Qid,
		User: in,
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

func (c *CoreUser) Get(ctx context.Context, in *userpb.UserCore) error {
	// c.storage.IsAQL()
	dAQL := &storepb.StoreAQL{
		User: in,
	}

	m, err := c.storage.AQL().ReadDocument(ctx, in.Qid.Key, dAQL)
	if err != nil {
		return fmt.Errorf("AQL().ReadDocument: %w", err)
	}

	// in.Qid.Key = m.Key
	in.Qid.Rev = m.Rev

	return nil
}
