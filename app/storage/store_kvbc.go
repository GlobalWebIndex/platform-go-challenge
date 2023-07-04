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

func (st *ServiceStoreKVBC) Create(key []byte, value []byte) error {
	if len(key) > maxValueCreate {
		return errTooBigValue
	}

	if st.cache.Has(key) {
		return errAlreadyExists
	}

	st.cache.Set(key, value)

	return nil
}

func (st *ServiceStoreKVBC) CreateBig(key []byte, value []byte) error {
	if st.cache.Has(key) {
		return errAlreadyExists
	}

	st.cache.SetBig(key, value)

	return nil
}

func (st *ServiceStoreKVBC) Exists(key []byte) bool {
	return st.cache.Has(key)
}

func (st *ServiceStoreKVBC) Get(key []byte) ([]byte, error) {
	if !st.cache.Has(key) {
		return []byte{}, errNotFound
	}

	return st.cache.Get([]byte{}, key), nil
}

func (st *ServiceStoreKVBC) GetBig(key []byte) ([]byte, error) {
	if !st.cache.Has(key) {
		return []byte{}, errNotFound
	}

	return st.cache.GetBig([]byte{}, key), nil
}

func (st *ServiceStoreKVBC) GetOK(key []byte) ([]byte, bool) {
	if !st.cache.Has(key) {
		return []byte{}, false
	}

	return st.cache.Get([]byte{}, key), true
}

func (st *ServiceStoreKVBC) GetBigOK(key []byte) ([]byte, bool) {
	if !st.cache.Has(key) {
		return []byte{}, false
	}

	return st.cache.GetBig([]byte{}, key), true
}
