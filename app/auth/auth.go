package auth

import (
	"context"
	"fmt"

	"x-gwi/app/instance"
)

type Auth struct {
	config *ConfigAuth
	inst   *instance.Instance
}

func NewAuth(_ context.Context, config *ConfigAuth, inst *instance.Instance) (*Auth, error) {
	a := &Auth{
		config: config,
		inst:   inst,
	}

	if !a.Valid() {
		return nil, fmt.Errorf("auth failure") //nolint:goerr113
	}

	return a, nil
}

func (a *Auth) Valid() bool {
	return a.config.Valid() &&
		a.inst.Valid()
}
