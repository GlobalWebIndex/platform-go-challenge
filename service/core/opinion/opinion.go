package opinion

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"x-gwi/app/storage"
	"x-gwi/app/x/id"
	sharepb "x-gwi/proto/core/_share/v1"
	storepb "x-gwi/proto/core/_store/v1"
	opinion "x-gwi/proto/core/opinion/v1"
	opinionpbapiv1 "x-gwi/proto/serv/opinion/v1"
	opinionpbapiv2 "x-gwi/proto/serv/opinion/v2"
	"x-gwi/service"
)

var (
	errKeyFromTo = errors.New("wrong key vs from to")
)

//nolint:unused
type CoreOpinion struct {
	opinon   *opinion.OpinionCore
	idx      *sharepb.ShareQID
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

func (c *CoreOpinion) Create(ctx context.Context, in *opinion.OpinionCore) error {
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
	// in.Qid.Key, from, to, err = c.edgeKeyFromTo(ctx, in)
	in.Qid.Key, from, to, err = c.edgeKeyFromTo(ctx, in)
	if err != nil {
		return fmt.Errorf("c.edgeKeyFromTo: %w", err)
	}

	dAQL := &storepb.StoreAQL{ //nolint:exhaustruct
		XFrom:   from,
		XTo:     to,
		XKey:    in.Qid.Key,
		Opinion: in,
	}

	m, err := c.storage.AQL().CreateDocument(ctx, dAQL, nil)
	if err != nil {
		return fmt.Errorf("AQL().CreateDocument: %w", err)
	}

	// if m.Key != in.Qid.Key {todo delete wronk key}
	in.Qid.Key = m.Key
	in.Qid.Rev = m.Rev

	return nil
}

func (c *CoreOpinion) Get(ctx context.Context, in *opinion.OpinionCore) error {
	// c.storage.IsAQL()
	//nolint:exhaustruct
	dAQL := &storepb.StoreAQL{
		Opinion: in,
	}

	m, err := c.storage.AQL().ReadDocument(ctx, in.Qid.Key, dAQL)
	if err != nil {
		return fmt.Errorf("AQL().ReadDocument: %w", err)
	}

	in.Qid.Rev = m.Rev

	return nil
}

func (c *CoreOpinion) Update(ctx context.Context, in *opinion.OpinionCore) error {
	key, _, _ := c.keyFromTo(in)
	if key != in.Qid.Key {
		return errKeyFromTo
	}
	// c.storage.IsAQL()
	//nolint:exhaustruct
	dAQL := &storepb.StoreAQL{
		Opinion: in,
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

// opinionpbapiv1 "x-gwi/proto/serv/opinion/v1"
func (c *CoreOpinion) ListV1(in *opinion.OpinionCore, stream opinionpbapiv1.OpinionService_ListServer) error { //nolint:lll
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

		o := dAQL.GetOpinion()
		o.Qid.Rev = m.Rev

		out := &opinionpbapiv1.ListResponse{
			Opinion: o,
		}

		err = stream.Send(out)
		if err != nil {
			return fmt.Errorf("stream.Send: %w", err)
		}
	}

	return nil
}

// opinionpbapiv2 "x-gwi/proto/serv/opinion/v2"
func (c *CoreOpinion) ListV2(in *opinion.OpinionCore, stream opinionpbapiv2.OpinionService_ListServer) error { //nolint:lll
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

		o := dAQL.GetOpinion()
		o.Qid.Rev = m.Rev

		err = stream.Send(o)
		if err != nil {
			return fmt.Errorf("stream.Send: %w", err)
		}
	}

	return nil
}

func (c *CoreOpinion) keyFromTo(in *opinion.OpinionCore) (string, string, string) {
	from := in.GetQidFromUser().GetKey()
	to := in.GetQidToAsset().GetKey()
	key := fmt.Sprintf("uaf:(%s,%s)", from, to)

	return key, from, to
}

func (c *CoreOpinion) edgeKeyFromTo(ctx context.Context, in *opinion.OpinionCore) (string, string, string, error) { //nolint:lll
	// from := in.GetQidFromUser().GetKey()
	// to := in.GetQidToAsset().GetKey()
	// key := fmt.Sprintf("uaf:(%s,%s)", from, to)
	key, from, to := c.keyFromTo(in)

	in.Qid.Key = key

	if from == "" {
		return "", "", "", fmt.Errorf("missed QID fromUser") //nolint:goerr113
	} else if to == "" {
		return "", "", "", fmt.Errorf("missed QID toAsset") //nolint:goerr113
	}

	exists, err := c.storage.AQL().DocumentExists(ctx, key)
	if err != nil {
		return "", "", "", fmt.Errorf("AQL().DocumentExists: opinion %w", err)
	} else if exists {
		return "", "", "", fmt.Errorf("opinion key already exists") //nolint:goerr113
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
