package instance

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"x-gwi/internal/env"
	"x-gwi/internal/id"
)

type Instance struct {
	name string
	mode string
	id   id.XID
}

func NewInstance() *Instance {
	return &Instance{
		id:   id.XiD(),
		name: env.Env("APP_NAME", "app"),
		mode: env.Env("APP_MODE", modes()[0]),
		// TENANT
	}
}

func modes() []string {
	return []string{"dev", "prod", "stage", "test"}
}

func (i *Instance) Valid() bool {
	if i.mode == "" {
		return false
	}

	var validMode bool

	for _, v := range modes() {
		if v == i.mode {
			validMode = true

			break
		}
	}

	return validMode &&
		i.name != "" &&
		!i.id.IsNil()
}

func (i *Instance) Name() string {
	return i.name
}

func (i *Instance) Mode() string {
	return i.mode
}

func (i *Instance) pass() []byte {
	prfx := fmt.Sprintf("%s::%s::", i.name, i.mode)

	return append([]byte(prfx), i.id.Bytes()...)
}

func (i *Instance) PassHash() []byte {
	passHash, _ := bcrypt.GenerateFromPassword(i.pass(), bcrypt.MinCost)

	return passHash
}

func (i *Instance) PassVerifyHash(passHash []byte) bool {
	// it is reverse verification of the instance without exposition of internal id
	return bcrypt.CompareHashAndPassword(passHash, i.pass()) == nil // no err
}
