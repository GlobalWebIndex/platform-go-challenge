package user

import (
	"context"
	"fmt"

	"x-gwi/app/storage"
	"x-gwi/app/storage/storepb2"
	"x-gwi/app/x/id"
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

	// storepb "x-gwi/proto/core/_store/v1"
	// issue with json_name = "_key" for gen annotations json - proto/core/_store/v1/store.proto
	//nolint:exhaustruct
	// _ = &storepb.StoreAQL{
	// 	Key: in.Qid.Key,
	// 	Qid: in.Qid,
	// 	User: &userpb.UserCore{
	// 		BasicAccount: in.GetBasicAccount(),
	// 		Md:           in.GetMd(),
	// 	},
	// }

	// if c.storage.IsAQL()
	// "x-gwi/app/storage/storepb2"
	//nolint:exhaustruct
	dAQL := &storepb2.StoreAQL{
		Key:  in.Qid.Key,
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
	// func (st *ServiceStoreAQL) ReadDocument(ctx context.Context, key string, result any) (MetaDocAQL, error)

	// if c.storage.IsAQL()
	// "x-gwi/app/storage/storepb2"
	//nolint:exhaustruct
	dAQL := &storepb2.StoreAQL{
		// Key:  in.Qid.Key,
		// Qid:  in.Qid,
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

/*
func (c *CoreUser) Create2(ctx context.Context, in *userpb.UserCore) (*userpb.UserCore, error) {
	rev := id.Rev()
	// QID on create
	in.Md.Label = &rev

	in.Qid.Kind = c.storage.CoreName().String()
	in.Qid.Key = in.BasicAccount.GetUsername()
	in.Qid.Uid = id.XiD().String()
	in.Qid.Uuid = id.UUID().String()

	// storepb "x-gwi/proto/core/_store/v1"
	// issue with json_name = "_key" for gen annotations json - proto/core/_store/v1/store.proto
	//nolint:exhaustruct
	// _ = &storepb.StoreAQL{
	// 	Key: in.Qid.Key,
	// 	Qid: in.Qid,
	// 	User: &userpb.UserCore{
	// 		BasicAccount: in.GetBasicAccount(),
	// 		Md:           in.GetMd(),
	// 	},
	// }

	// "x-gwi/app/storage/storepb2"
	//nolint:exhaustruct
	dAQL := &storepb2.StoreAQL{
		Key:  in.Qid.Key,
		Qid:  in.Qid,
		User: in,
	}

	m, err := c.storage.AQL().CreateDocument(ctx, dAQL)
	if err != nil {
		return nil, fmt.Errorf("AQL().CreateDocument: %w", err)
	}

	out := &userpb.UserCore{
		Qid: &sharepb.ShareQID{
			Rev: m.Rev,
			Key: m.Key,
		},
		BasicAccount: in.GetBasicAccount(),
		Md:           in.GetMd(),
	}

	out.Qid = in.GetQid()

	out.Qid.Rev = m.Rev
	out.Qid.Key = m.Key

	out.BasicAccount = in.GetBasicAccount()
	out.Md = in.GetMd()

	return out, nil
}
*/
