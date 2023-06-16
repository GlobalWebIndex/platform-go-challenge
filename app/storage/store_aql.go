package storage

import (
	"context"
	"fmt"
	"time"

	arango "github.com/arangodb/go-driver"
	arhttp "github.com/arangodb/go-driver/http"

	"x-gwi/app/instance"
	"x-gwi/app/logs"
	"x-gwi/service"
)

type AppStoreAQL struct {
	inst   *instance.Instance
	cli    ClientAQL
	db     DatabaseAQL
	stores map[service.CoreName]*ServiceStoreAQL
}

type ServiceStoreAQL struct {
	inst *instance.Instance
	col  CollectionAQL
	name service.CoreName
}

type ClientAQL = arango.Client
type DatabaseAQL = arango.Database
type CollectionAQL = arango.Collection
type GraphAQL = arango.Graph
type SearchViewAQL = arango.ArangoSearchView
type MetaDocAQL = arango.DocumentMeta

func (st *AppStoreAQL) initAppStoreAQL(ctx context.Context, apSt *AppStorage) error {
	t := time.Now()
	st.inst = apSt.inst
	st.stores = make(map[service.CoreName]*ServiceStoreAQL)

	// connection to arango db cluster
	conn, err := arhttp.NewConnection(arhttp.ConnectionConfig{ //nolint:exhaustruct
		// Endpoints: []string{"http://localhost:8529"},
		Endpoints: apSt.config.AQL.Endpoints,
	})
	if err != nil {
		return fmt.Errorf("arhttp.NewConnection: %w", err)
	}

	st.cli, err = arango.NewClient(arango.ClientConfig{ //nolint:exhaustruct
		Connection: conn,
		Authentication: arango.BasicAuthentication(
			apSt.config.AQL.Credentials.UserName,
			apSt.config.AQL.Credentials.PassWord),
	})
	if err != nil {
		return fmt.Errorf("arango.NewClient: %w", err)
	}

	if err = st.initDB(ctx); err != nil {
		return fmt.Errorf("st.initDB: %w", err)
	}

	dbInfo, _ := st.db.Info(ctx)

	logs.Debug().
		Interface("dbInfo", dbInfo).
		Dur("duration_ms", time.Since(t)).
		Send()

	return nil
}

func (st *AppStoreAQL) initDB(ctx context.Context) error {
	dbName := fmt.Sprintf("db_gwi_%s_%s", st.inst.Name(), st.inst.Mode())
	_ = dbName

	// Using WithArangoQueueTimeout we get Timout error if not response after time specified in WithArangoQueueTime
	ctxDB := arango.WithArangoQueueTimeout(ctx, true)
	ctxDB = arango.WithArangoQueueTime(ctxDB, 10*time.Second) //nolint:gomnd

	// Ensure app database mode
	dbExists, err := st.cli.DatabaseExists(ctxDB, dbName)
	if err != nil {
		return fmt.Errorf("cli.DatabaseExists: %w", err)
	}

	if dbExists {
		if st.db, err = st.cli.Database(ctxDB, dbName); err != nil {
			return fmt.Errorf("cli.CreateDatabase: %w", err)
		}
	} else {
		if st.db, err = st.cli.CreateDatabase(ctxDB, dbName, nil); err != nil {
			return fmt.Errorf("cli.CreateDatabase: %w", err)
		}
	}

	// collection per service
	for _, coreName := range service.CoreNames() {
		if st.stores[coreName] == nil {
			st.stores[coreName] = new(ServiceStoreAQL)

			st.stores[coreName].inst = st.inst
			st.stores[coreName].name = coreName
		}

		if err = st.initCollectionCore(ctxDB, coreName); err != nil {
			return fmt.Errorf("st.initCollection: %w", err)
		}
	}

	return nil
}

func (st *AppStoreAQL) initCollectionCore(ctx context.Context, coreName service.CoreName) error {
	name := fmt.Sprintf("srv_core_%s", coreName)
	// name := string(coreName)

	exists, err := st.db.CollectionExists(ctx, name)
	if err != nil {
		return fmt.Errorf("db.CollectionExists: %w", err)
	}

	if exists {
		if st.stores[coreName].col, err = st.db.Collection(ctx, name); err != nil {
			return fmt.Errorf("db.Collection: %w", err)
		}
	} else {
		if st.stores[coreName].col, err = st.db.CreateCollection(ctx, name, nil); err != nil {
			return fmt.Errorf("db.CreateCollection: %w", err)
		}
	}

	return nil
}

// CreateDocument creates a single document in the collection.
// The document data is loaded from the given document, the document meta data is returned.
// If the document data already contains a `_key` field, this will be used as key of the new document,
// otherwise a unique key is created. A ConflictError is returned when a `_key` field contains a duplicate key,
// other any other field violates an index constraint.
// To return the NEW document, prepare a context with WithReturnNew. To wait until document has been synced to disk,
// prepare a context with `WithWaitForSync`.
func (st *ServiceStoreAQL) CreateDocument(ctx context.Context, doc any) (MetaDocAQL, error) {
	meta, err := st.col.CreateDocument(ctx, doc)
	if err != nil {
		return MetaDocAQL{}, fmt.Errorf("col.CreateDocument: %w", err)
	}

	return meta, nil
}

// ReadDocument reads a single document with given key from the collection.
// The document data is stored into result, the document meta data is returned.
// If no document exists with given key, a NotFoundError is returned.
func (st *ServiceStoreAQL) ReadDocument(ctx context.Context, key string, result any) (MetaDocAQL, error) {
	meta, err := st.col.ReadDocument(ctx, key, result)
	if err != nil {
		return MetaDocAQL{}, fmt.Errorf("col.ReadDocument: %w", err)
	}

	return meta, nil
}

// DocumentExists checks if a document with given key exists in the collection.
func (st *ServiceStoreAQL) DocumentExists(ctx context.Context, key string) (bool, error) {
	exists, err := st.col.DocumentExists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("col.DocumentExists: %w", err)
	}

	return exists, nil
}

// .WithRevision(
