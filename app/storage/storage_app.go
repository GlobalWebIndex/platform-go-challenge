package storage

import (
	"context"
	"fmt"

	"x-gwi/app/auth"
	"x-gwi/app/instance"
	"x-gwi/app/logs"
	"x-gwi/service"
)

type AppStorage struct {
	config *ConfigAppStorage
	inst   *instance.Instance
	stAQL  *AppStoreAQL
	stKVBC *AppStoreKVBC
	stores map[service.CoreName]*ServiceStorage
}

func NewAppStorage(ctx context.Context, config *ConfigAppStorage, inst *instance.Instance, _ *auth.Auth) (*AppStorage, error) { //nolint:lll
	apSt := &AppStorage{
		config: config,
		inst:   inst,
		stAQL:  new(AppStoreAQL),
		stKVBC: new(AppStoreKVBC),
		stores: make(map[service.CoreName]*ServiceStorage),
	}

	if err := apSt.stAQL.initAppStoreAQL(ctx, apSt); err != nil {
		// allows cache only (in dev)
		if apSt.inst.Mode() != string(instance.ModeDev) {
			return nil, fmt.Errorf("stAQL.initAppStoreAQL: %w", err)
		}

		apSt.stAQL = nil

		logs.Warn().
			Err(err).
			Str("mode", apSt.inst.Mode()).
			Msg("env mode dev allows cache only")
	}

	if err := apSt.stKVBC.initAppStoreKVBC(ctx, apSt); err != nil {
		return nil, fmt.Errorf("stKVB.initAppStoreKVBC: %w", err)
	}

	for _, coreName := range service.CoreNames() {
		var ok bool

		apSt.stores[coreName] = new(ServiceStorage)
		apSt.stores[coreName].inst = apSt.inst
		apSt.stores[coreName].name = coreName

		// apSt.stores[coreName].cache = fastcache.New(1024) //nolint:gomnd

		// allows cache only in dev
		if apSt.stAQL != nil {
			apSt.stores[coreName].sstAQL, ok = apSt.stAQL.stores[coreName]
			if !ok {
				return nil, fmt.Errorf("apSt.stAQL.stores[coreName]") //nolint:goerr113
			}

			apSt.stores[coreName].isAQL = true
		}

		apSt.stores[coreName].sstKVBC, ok = apSt.stKVBC.stores[coreName]
		if !ok {
			return nil, fmt.Errorf("apSt.stKVBC.stores[coreName]") //nolint:goerr113
		}

		apSt.stores[coreName].isKVBC = true
	}

	return apSt, nil
}

func (s *AppStorage) Valid() bool {
	return s.config.Valid() &&
		s.inst.Valid()
}

func (s *AppStorage) ServiceStore(coreName service.CoreName) *ServiceStorage {
	srvStore, ok := s.stores[coreName]
	if !ok {
		logs.Fatal().Str("coreName", string(coreName)).Send()
	}

	return srvStore
}
