package storage

import (
	"context"

	"github.com/VictoriaMetrics/fastcache"

	"x-gwi/app/instance"
	"x-gwi/service"
)

type AppStoreKVBC struct {
	inst   *instance.Instance
	stores map[service.CoreName]*ServiceStoreKVBC
}

type ServiceStoreKVBC struct {
	inst  *instance.Instance
	cache *CacheKVBC
	name  service.CoreName
}

type CacheKVBC = fastcache.Cache

func (st *AppStoreKVBC) initAppStoreKVBC(_ context.Context, apSt *AppStorage) error { //nolint:unparam
	st.inst = apSt.inst
	st.stores = make(map[service.CoreName]*ServiceStoreKVBC)

	for _, coreName := range service.CoreNames() {
		if st.stores[coreName] == nil {
			st.stores[coreName] = new(ServiceStoreKVBC)

			st.stores[coreName].inst = st.inst
			st.stores[coreName].name = coreName

			st.stores[coreName].cache = fastcache.New(1024) //nolint:gomnd
		}
	}

	return nil
}

/*
func (st *ServiceStoreKVBC) CreateDocument(ctx context.Context, doc any) (MetaDocAQL, error) {
	meta, err := st.col.CreateDocument(ctx, doc)
	if err != nil {
		return MetaDocAQL{}, fmt.Errorf("col.CreateDocument: %w", err)
	}

	return meta, nil
}

func (st *ServiceStoreKVBC) ReadDocument(ctx context.Context, key string, result any) (MetaDocAQL, error) {
	meta, err := st.col.ReadDocument(ctx, key, result)
	if err != nil {
		return MetaDocAQL{}, fmt.Errorf("col.ReadDocument: %w", err)
	}

	return meta, nil
}

func (st *ServiceStoreAQL) DocumentExists(ctx context.Context, key string) (bool, error) {
	exists, err := st.col.DocumentExists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("col.DocumentExists: %w", err)
	}

	return exists, nil
}
*/
