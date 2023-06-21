package opinion

import (
	"context"
	"fmt"

	"x-gwi/app/storage"
	"x-gwi/app/x/id"
	sharepb "x-gwi/proto/core/_share/v1"
	storepb "x-gwi/proto/core/_store/v1"
	opinonpb "x-gwi/proto/core/opinion/v1"
	"x-gwi/service"
)

//nolint:unused
type CoreOpinion struct {
	opinon   *opinonpb.OpinionCore
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

func (c *CoreOpinion) Create(ctx context.Context, in *opinonpb.OpinionCore) error {
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
		Qid:     in.Qid,
		Opinion: in,
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

func (c *CoreOpinion) Get(ctx context.Context, in *opinonpb.OpinionCore) error {
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

func (c *CoreOpinion) edgeKeyFromTo(ctx context.Context, in *opinonpb.OpinionCore) (string, string, string, error) { //nolint:lll
	from := in.GetQidFromUser().GetKey()
	to := in.GetQidToAsset().GetKey()

	in.Qid.Key = fmt.Sprintf("uao:(%s,%s)", from, to)

	if from == "" {
		return "", "", "", fmt.Errorf("missed QID fromUser") //nolint:goerr113
	} else if to == "" {
		return "", "", "", fmt.Errorf("missed QID toAsset") //nolint:goerr113
	}

	key := fmt.Sprintf("uao:(%s,%s)", from, to)

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
