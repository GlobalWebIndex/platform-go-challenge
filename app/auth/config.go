package auth

import (
	"x-gwi/app/x/env"
)

type ConfigAuth struct {
	AppRoot *BasicAuth
}

type BasicAuth struct {
	UserName string
	PassWord string
	PassHash string
}

func NewConfigAuth() *ConfigAuth {
	return &ConfigAuth{
		AppRoot: &BasicAuth{
			UserName: env.Env("AUTH_APPROOT_USERNAME", "app-root"),
			PassWord: env.Env("AUTH_APPROOT_PASSWORD", ""),
			PassHash: env.Env("AUTH_APPROOT_PASSHASH", "$2a$08$rvzL6XK2K/UPnK1fkHS3TeDWyc9nbYMAnZBSkiiWCamNilcYxA4jK"),
		},
	}
}

func (a *ConfigAuth) Valid() bool {
	return true
}
