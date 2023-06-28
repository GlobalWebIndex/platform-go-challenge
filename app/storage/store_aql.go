package storage

import (
	"context"
	"fmt"
	"net"
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
	inst   *instance.Instance
	apSt   *AppStoreAQL
	col    CollectionAQL
	name   service.CoreName
	isEdge bool
}

type ClientAQL = arango.Client
type DatabaseAQL = arango.Database
type CollectionAQL = arango.Collection
type GraphAQL = arango.Graph
type SearchViewAQL = arango.ArangoSearchView
type MetaDocAQL = arango.DocumentMeta
type CursorAQL = arango.Cursor
type QueryAQL struct {
	Bind  map[string]any
	Query string
}

func (st *AppStoreAQL) initAppStoreAQL(ctx context.Context, apSt *AppStorage) error {
	st.inst = apSt.inst
	st.stores = make(map[service.CoreName]*ServiceStoreAQL)

	var err error

	t := time.Now()
	l := logs.LogC3.With().
		Str("mode", st.inst.Mode()).
		Strs("host_ip", apSt.config.HostIP).
		Strs("endpoints", apSt.config.AQL.Endpoints).
		Dur("duration_ms", time.Since(t)).
		Logger()

	defer logs.DebugOnDefer(&l, t, err)

	st.verifyEndpoints(ctx, apSt)
	l = l.With().Strs("endpoints_ver", apSt.config.AQL.Endpoints).Logger()

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
	l = l.With().Interface("dbInfo", dbInfo).Logger()

	return nil
}

// ads host ip in containerized env in dev and test modes
func (st *AppStoreAQL) verifyEndpoints(_ context.Context, apSt *AppStorage) {
	if (st.inst.Mode() == instance.ModeDev.String() || st.inst.Mode() == instance.ModeTest.String()) &&
		len(apSt.config.HostIP) >= 1 &&
		len(apSt.config.AQL.Endpoints) == 1 &&
		apSt.config.AQL.Endpoints[0] == defAdrAQL {
		hostIP := []string{}
		emptyV := true

		for _, v := range apSt.config.HostIP {
			if v != "" {
				ip := net.ParseIP(v)
				if ip != nil {
					v = fmt.Sprintf("http://%s", net.JoinHostPort(ip.String(), defPortAQL))
					hostIP = append(hostIP, v)
					emptyV = false
				}
			}
		}

		if !emptyV {
			apSt.config.AQL.Endpoints = append(apSt.config.AQL.Endpoints, hostIP...)
		}
	}
}

func (st *AppStoreAQL) initDB(ctx context.Context) error {
	dbName := fmt.Sprintf("db_%s_%s", st.inst.Name(), st.inst.Mode())
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
			st.stores[coreName].apSt = st
			st.stores[coreName].name = coreName
			st.stores[coreName].isEdge = coreName.IsEdge()
		}

		if err = st.initCollectionCore(ctxDB, coreName); err != nil {
			return fmt.Errorf("st.initCollection: %w", err)
		}
	}

	return nil
}

func (st *AppStoreAQL) initCollectionCore(ctx context.Context, coreName service.CoreName) error {
	// name := fmt.Sprintf("srv_core_%s", coreName)
	name := string(coreName)

	exists, err := st.db.CollectionExists(ctx, name)
	if err != nil {
		return fmt.Errorf("db.CollectionExists: %w", err)
	}

	if exists { //nolint:nestif
		if st.stores[coreName].col, err = st.db.Collection(ctx, name); err != nil {
			return fmt.Errorf("db.Collection: %w", err)
		}
	} else {
		opts := &arango.CreateCollectionOptions{ //nolint:exhaustruct
			CacheEnabled: func(b bool) *bool { return &b }(true),
			WaitForSync:  true,
		}

		if st.stores[coreName].isEdge {
			opts.Type = arango.CollectionTypeEdge
			if st.stores[coreName].col, err = st.db.CreateCollection(ctx, name, opts); err != nil {
				return fmt.Errorf("db.CreateCollection: %w", err)
			}
		} else {
			if st.stores[coreName].col, err = st.db.CreateCollection(ctx, name, opts); err != nil {
				return fmt.Errorf("db.CreateCollection: %w", err)
			}
		}
	}

	return nil
}

func (st *ServiceStoreAQL) CollectionName() string {
	if st == nil || st.col == nil {
		return ""
	}

	return st.col.Name()
}

func (st *ServiceStoreAQL) OtherCollectionName(coreName service.CoreName) string {
	srvStoreAQL, ok := st.apSt.stores[coreName]
	if ok {
		return srvStoreAQL.CollectionName()
	}

	return ""
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

// DocumentExists checks if a document with given key exists in the core collection.
func (st *ServiceStoreAQL) DocumentExists(ctx context.Context, key string) (bool, error) {
	exists, err := st.col.DocumentExists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("col.DocumentExists: %w", err)
	}

	return exists, nil
}

// DocumentExists checks if a document with given key exists in other core collection.
func (st *ServiceStoreAQL) OtherCoreDocumentExists(ctx context.Context, key string, coreName service.CoreName) (
	bool, error) {
	srvStoreAQL, ok := st.apSt.stores[coreName]
	if !ok {
		return false, fmt.Errorf("collectionCore doesn't exists") //nolint:goerr113
	}

	exists, err := srvStoreAQL.col.DocumentExists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("col.DocumentExists: %w", err)
	}

	return exists, nil
}

func (st *ServiceStoreAQL) ListFrom(ctx context.Context, keyFrom string, coreFrom service.CoreName) (CursorAQL, int64, error) { //nolint:ireturn,nolintlint,lll
	// FOR doc IN @@collection
	// FILTER doc._from == @value
	// RETURN doc
	fromID := fmt.Sprintf("%s/%s", st.OtherCollectionName(coreFrom), keyFrom)
	qu := fmt.Sprintf(`FOR doc IN %s FILTER doc._from == "%s" RETURN doc`, st.CollectionName(), fromID)
	// bi := make(map[string]any)
	// bi["fromUser"] = fromID

	ctx = arango.WithQueryCount(ctx, true)

	cursor, err := st.col.Database().Query(ctx, qu, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("ListFrom_Query: %w", err)
	}

	return cursor, cursor.Count(), nil
}

func (st *ServiceStoreAQL) ListTo(ctx context.Context, keyTo string, coreTo service.CoreName) (CursorAQL, int64, error) { //nolint:ireturn,nolintlint,lll
	// FOR doc IN @@collection
	// FILTER doc._to == @value
	// RETURN doc
	fromID := fmt.Sprintf("%s/%s", st.OtherCollectionName(coreTo), keyTo)
	qu := fmt.Sprintf(`FOR doc IN %s FILTER doc._to == "%s" RETURN doc`, st.CollectionName(), fromID)

	ctx = arango.WithQueryCount(ctx, true)

	cursor, err := st.col.Database().Query(ctx, qu, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("ListTo_Query: %w", err)
	}

	return cursor, cursor.Count(), nil
}

func (st *ServiceStoreAQL) ListAll(ctx context.Context) (CursorAQL, int64, error) { //nolint:ireturn
	// FOR doc IN @@collection
	// RETURN doc
	qu := fmt.Sprintf(`FOR doc IN %s RETURN doc`, st.CollectionName())

	ctx = arango.WithQueryCount(ctx, true)

	cursor, err := st.col.Database().Query(ctx, qu, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("ListAll_Query: %w", err)
	}

	return cursor, cursor.Count(), nil
}

// Query performs an AQL query, returning a cursor used to iterate over the returned documents.
// Note that the returned Cursor must always be closed to avoid holding on to resources in the server
// while they are no longer needed - cursor.Close().
func (st *ServiceStoreAQL) Query(ctx context.Context, q *QueryAQL) (CursorAQL, int64, error) { //nolint:ireturn
	ctx = arango.WithQueryCount(ctx, true)

	cursor, err := st.col.Database().Query(ctx, q.Query, q.Bind)
	if err != nil {
		return nil, 0, fmt.Errorf("col.Database().Query: %w", err)
	}

	return cursor, cursor.Count(), nil
}

// .WithRevision(

/* insert some values for the @@collection and @value bind parameters
uaf_favourite
u_user/fk-u:390
FOR doc IN @@collection
FILTER doc._from == @value
RETURN {
_key: doc._key,
_from: doc._from,
_to: doc._to,
_rev: doc._rev,
md: doc.favourite.md,
favourite: doc
}

FOR doc IN @@collection
  FILTER doc._from == @value
  LIMIT 2
  RETURN {
  _key: doc._key,
  _from: doc._from,
  _to: doc._to,
  _rev: doc._rev,
  md: doc.favourite.md,
  favourite: doc
  }
*/
/*
//uaf_favourite
//u_user/fk-u:390
//FOR doc IN @@collection
//FILTER doc._from == @value
FOR doc IN uaf_favourite
FILTER doc._from == 'u_user/fk-u:159'
//FILTER doc._to == 'a_asset/fk-a:1050'
RETURN doc
/*
RETURN {
_key: doc._key,
_from: doc._from,
_to: doc._to,
_rev: doc._rev,
md: doc.favourite.md,
favourite: doc
}
*/
/*
  	Query(ctx context.Context, query string, bindVars map[string]interface{}) (arango.Cursor, error)
	Query performs an AQL query, returning a cursor used to iterate over the returned documents.
	Note that the returned Cursor must always be closed to avoid holding on to resources in the server
	while they are no longer needed.
	query := ""
	bindVars := make(BindVarsAGL)
	cursor, err := st.col.Database().Query(ctx, query, bindVars)
	_ = cursor
	_ = err
	_ = cursor
	defer cursor.Close()

*/
