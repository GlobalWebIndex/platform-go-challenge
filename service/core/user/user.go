package user

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"x-gwi/app/storage"
	"x-gwi/app/x/id"
	store_pb "x-gwi/proto/core/_store/v1"
	user_pb "x-gwi/proto/core/user/v1"
	"x-gwi/service"
)

type CoreUser struct {
	storage  *storage.ServiceStorage
	coreName service.CoreName
}

type WorkUser struct {
	In  *user_pb.UserInstance
	Out *user_pb.UserInstance
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

func (c *CoreUser) Create(_ context.Context, wrk *WorkUser) error {
	rev := time.Now().UnixNano()
	name := wrk.In.GetUsername()
	wrk.Out.Username = name
	// wrk.User2.Username = "aaaa"
	_ = wrk.In.GetUsername()
	// wrk.User2 = nil

	wrk.Out.Id = &store_pb.StoreIDX{
		Kind: string(c.storage.CoreName()),
		Uid:  id.XiD().String(),
		Rev:  strconv.FormatInt(rev, 36),
		Uuid: id.UUID().String(),
		RevN: &rev,
		Aql:  nil,
		// Aql: &store_pb.MetaDocAQL{
		// 	Key:    wrk.Out.Id.Uid,
		// 	Id:     fmt.Sprintf("%s/%s", c.coreName, wrk.Out.Id.Uid),
		// 	Rev:    wrk.Out.Id.Rev,
		// 	OldRev: "",
		// },
	}
	wrk.Out.Id.Aql = &store_pb.MetaDocAQL{
		Key:    wrk.Out.Id.Uid,
		Id:     fmt.Sprintf("%s/%s", c.coreName, wrk.Out.Id.Uid),
		Rev:    wrk.Out.Id.Rev,
		OldRev: "",
	}
	// wrk.Out.Id.Kind = string(c.storage.CoreName())
	// wrk.Out.Id.Uid = id.XiD().String()
	// wrk.Out.Id.Rev = strconv.FormatInt(rev, 36)
	// wrk.Out.Id.Uuid = id.UUID().String()
	// wrk.Out.Id.RevN = &rev
	// wrk.Out.Id.Aql = &store_pb.MetaDocAQL{
	// 	Key:    wrk.Out.Id.Uid,
	// 	Id:     fmt.Sprintf("%s/%s", c.coreName, wrk.Out.Id.Uid),
	// 	Rev:    wrk.Out.Id.Rev,
	// 	OldRev: "",
	// }

	return nil
}
