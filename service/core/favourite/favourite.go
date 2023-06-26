package favourite

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"x-gwi/app/storage"
	"x-gwi/app/x/id"
	sharepb "x-gwi/proto/core/_share/v1"
	storepb "x-gwi/proto/core/_store/v1"
	favouritepb "x-gwi/proto/core/favourite/v1"
	favouritepbapiv1 "x-gwi/proto/serv/favourite/v1"
	favouritepbapiv2 "x-gwi/proto/serv/favourite/v2"
	"x-gwi/service"
)

//nolint:unused
type CoreFavourite struct {
	favourite *favouritepb.FavouriteCore
	idx       *sharepb.ShareQID
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

func (c *CoreFavourite) Create(ctx context.Context, in *favouritepb.FavouriteCore) error {
	var (
		err  error
		from string
		to   string
	)
	// check DocumentExists for from and to
	in.Qid.Kind = c.storage.CoreName().String()
	// in.Qid.Key = in.Qid.Key // use directly
	in.Qid.Uid = id.XiD().String()
	// in.Qid.Uuid = id.UUID().String() // lower size

	// c.storage.IsAQL()
	in.Qid.Key, from, to, err = c.edgeKeyFromTo(ctx, in)
	if err != nil {
		return fmt.Errorf("c.edgeKeyFromTo: %w", err)
	}

	dAQL := &storepb.StoreAQL{ //nolint:exhaustruct
		XFrom:     from,
		XTo:       to,
		XKey:      in.Qid.Key,
		Qid:       in.Qid,
		Favourite: in,
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

func (c *CoreFavourite) Get(ctx context.Context, in *favouritepb.FavouriteCore) error {
	// c.storage.IsAQL()
	//nolint:exhaustruct
	dAQL := &storepb.StoreAQL{
		Favourite: in,
	}

	m, err := c.storage.AQL().ReadDocument(ctx, in.Qid.Key, dAQL)
	if err != nil {
		return fmt.Errorf("AQL().ReadDocument: %w", err)
	}

	in.Qid.Rev = m.Rev

	return nil
}

// favouritepbapiv1 "x-gwi/proto/serv/favourite/v1"
func (c *CoreFavourite) ListV1(in *favouritepb.FavouriteCore, stream favouritepbapiv1.FavouriteService_ListServer) error { //nolint:lll
	keyFrom := in.GetQidFromUser().GetKey()
	if keyFrom == "" {
		return status.Errorf(codes.InvalidArgument, "missed qid_from_user.key")
	}

	cursor, count, err := c.storage.AQL().ListFrom(stream.Context(), keyFrom, service.NameUser)
	if err != nil {
		return fmt.Errorf("AQL().ListFrom: %w", err)
	}
	defer cursor.Close()

	if count == 0 {
		return nil
	}

	for cursor.HasMore() {
		if stream.Context().Err() != nil {
			return fmt.Errorf("stream.Context: %w", err)
		}

		dAQL := new(storepb.StoreAQL)

		m, err := cursor.ReadDocument(stream.Context(), dAQL)
		if err != nil {
			return fmt.Errorf("cursor.ReadDocument: %w", err)
		}

		if m.Rev == "" {
			return nil
		}

		f := dAQL.GetFavourite()
		f.Qid.Rev = m.Rev

		out := &favouritepbapiv1.ListResponse{
			Favourite: f,
		}

		err = stream.Send(out)
		if err != nil {
			return fmt.Errorf("stream.Send: %w", err)
		}
	}

	return nil
}

// favouritepbapiv2 "x-gwi/proto/serv/favourite/v2"
func (c *CoreFavourite) ListV2(in *favouritepb.FavouriteCore, stream favouritepbapiv2.FavouriteService_ListServer) error { //nolint:lll
	keyFrom := in.GetQidFromUser().GetKey()
	if keyFrom == "" {
		return status.Errorf(codes.InvalidArgument, "missed qid_from_user.key")
	}

	cursor, count, err := c.storage.AQL().ListFrom(stream.Context(), keyFrom, service.NameUser)
	if err != nil {
		return fmt.Errorf("AQL().ListFrom: %w", err)
	}
	defer cursor.Close()

	if count == 0 {
		return nil
	}

	for cursor.HasMore() {
		if stream.Context().Err() != nil {
			return fmt.Errorf("stream.Context: %w", err)
		}

		dAQL := new(storepb.StoreAQL)

		m, err := cursor.ReadDocument(stream.Context(), dAQL)
		if err != nil {
			return fmt.Errorf("cursor.ReadDocument: %w", err)
		}

		if m.Rev == "" {
			return nil
		}

		f := dAQL.GetFavourite()
		f.Qid.Rev = m.Rev

		err = stream.Send(f)
		if err != nil {
			return fmt.Errorf("stream.Send: %w", err)
		}
	}

	return nil
}

func (c *CoreFavourite) edgeKeyFromTo(ctx context.Context, in *favouritepb.FavouriteCore) (string, string, string, error) { //nolint:lll
	from := in.GetQidFromUser().GetKey()
	to := in.GetQidToAsset().GetKey()

	in.Qid.Key = fmt.Sprintf("uaf:(%s,%s)", from, to)

	if from == "" {
		return "", "", "", fmt.Errorf("missed QID fromUser") //nolint:goerr113
	} else if to == "" {
		return "", "", "", fmt.Errorf("missed QID toAsset") //nolint:goerr113
	}

	key := fmt.Sprintf("uaf:(%s,%s)", from, to)

	exists, err := c.storage.AQL().DocumentExists(ctx, key)
	if err != nil {
		return "", "", "", fmt.Errorf("AQL().DocumentExists: favourite %w", err)
	} else if exists {
		return "", "", "", fmt.Errorf("favourite key already exists") //nolint:goerr113
	}

	exists, err = c.storage.AQL().OtherCoreDocumentExists(ctx, from, service.NameUser)
	if err != nil {
		return "", "", "", fmt.Errorf("AQL().OtherCoreDocumentExists: user %w", err)
	} else if !exists {
		return "", "", "", fmt.Errorf("unknown user") //nolint:goerr113
	}

	exists, err = c.storage.AQL().OtherCoreDocumentExists(ctx, to, service.NameAsset)
	if err != nil {
		return "", "", "", fmt.Errorf("AQL().OtherCoreDocumentExists: asset %w", err)
	} else if !exists {
		return "", "", "", fmt.Errorf("unknown asset") //nolint:goerr113
	}

	from = fmt.Sprintf("%s/%s", service.NameUser, from)
	to = fmt.Sprintf("%s/%s", service.NameAsset, to)

	return key, from, to, nil
}
