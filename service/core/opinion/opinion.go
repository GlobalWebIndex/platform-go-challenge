package opinion

import (
	"context"
	"fmt"

	"x-gwi/app/storage"
	"x-gwi/app/storage/storepb2"
	"x-gwi/app/x/id"
	sharepb "x-gwi/proto/core/_share/v1"
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

//nolint:nilnil,revive
func (c *CoreOpinion) Create1(ctx context.Context, in *opinonpb.OpinionCore) (*sharepb.ShareQID, error) {
	return nil, nil
}

func (c *CoreOpinion) Create(ctx context.Context, in *opinonpb.OpinionCore) error {
	var err error
	// check DocumentExists for from and to
	in.Qid.Kind = c.storage.CoreName().String()
	// in.Qid.Key = in.Qid.Key // use directly
	in.Qid.Uid = id.XiD().String()
	// in.Qid.Uuid = id.UUID().String() // lower size

	in.Qid.Key, err = c.newKeyFromTo(ctx, in)
	if err != nil {
		return fmt.Errorf("c.newKeyFromTo: %w", err)
	}

	// if c.storage.IsAQL()
	// "x-gwi/app/storage/storepb2"
	//nolint:exhaustruct
	dAQL := &storepb2.StoreAQL{
		Key:     in.Qid.Key,
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
	// if c.storage.IsAQL()
	// "x-gwi/app/storage/storepb2"
	//nolint:exhaustruct
	dAQL := &storepb2.StoreAQL{
		// Key:  in.Qid.Key,
		// Qid:  in.Qid,
		Opinion: in,
	}

	m, err := c.storage.AQL().ReadDocument(ctx, in.Qid.Key, dAQL)
	if err != nil {
		return fmt.Errorf("AQL().ReadDocument: %w", err)
	}

	// in.Qid.Key = m.Key
	in.Qid.Rev = m.Rev

	return nil
}

func (c *CoreOpinion) newKeyFromTo(ctx context.Context, in *opinonpb.OpinionCore) (string, error) {
	fromUser := in.GetQidFromUser().GetKey()
	toAsset := in.GetQidToAsset().GetKey()

	if fromUser == "" {
		return "", fmt.Errorf("missed QID fromUser") //nolint:goerr113
	} else if toAsset == "" {
		return "", fmt.Errorf("missed QID toAsset") //nolint:goerr113
	}

	favouriteKey := fmt.Sprintf("oua_%s_%s", fromUser, toAsset)

	exists, err := c.storage.AQL().DocumentExists(ctx, favouriteKey)
	if err != nil {
		return "", fmt.Errorf("AQL().DocumentExists: favourite %w", err)
	} else if exists {
		return "", fmt.Errorf("favourite key already exists") //nolint:goerr113
	}

	exists, err = c.storage.AQL().OtherCoreDocumentExists(ctx, fromUser, service.NameUser)
	if err != nil {
		return "", fmt.Errorf("AQL().OtherCoreDocumentExists: user %w", err)
	} else if !exists {
		return "", fmt.Errorf("unknown user") //nolint:goerr113
	}

	exists, err = c.storage.AQL().OtherCoreDocumentExists(ctx, toAsset, service.NameAsset)
	if err != nil {
		return "", fmt.Errorf("AQL().OtherCoreDocumentExists: asset %w", err)
	} else if !exists {
		return "", fmt.Errorf("unknown asset") //nolint:goerr113
	}

	return favouriteKey, nil
}
