package storage

import (
	"context"

	"github.com/VictoriaMetrics/fastcache"

	"x-gwi/app/auth"
	"x-gwi/app/instance"
)

type Storage struct {
	config    *ConfigStorage
	inst      *instance.Instance
	FastCache *fastcache.Cache
	FastCaIDX *fastcache.Cache
}

func NewStorage(_ context.Context, config *ConfigStorage, inst *instance.Instance, _ *auth.Auth) (*Storage, error) {
	s := &Storage{ //nolint:exhaustruct
		config: config,
		inst:   inst,
	}

	s.FastCache = fastcache.New(1024) //nolint:gomnd
	s.FastCaIDX = fastcache.New(1024) //nolint:gomnd

	return s, nil
}

func (s *Storage) Valid() bool {
	return s.config.Valid() &&
		s.inst.Valid() &&
		s.FastCache != nil &&
		s.FastCaIDX != nil
}
