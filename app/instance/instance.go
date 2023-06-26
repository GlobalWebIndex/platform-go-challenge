package instance

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"x-gwi/app/x/env"
	"x-gwi/app/x/id"
)

type Instance struct {
	name string
	mode InstMode
	id   id.XID
}

type InstMode string

const (
	DefAppName          = "app-gwi"
	ModeDev    InstMode = InstMode("dev")
	ModeProd   InstMode = InstMode("prod")
	ModeTest   InstMode = InstMode("test")
	ModeStage  InstMode = InstMode("stage")
)

func InstModes() []InstMode {
	return []InstMode{
		ModeDev,
		ModeProd,
		ModeTest,
		ModeStage,
	}
}

func (i InstMode) String() string {
	return string(i)
}

func (i InstMode) Valid() bool {
	for _, v := range InstModes() {
		if i == v {
			return true
		}
	}

	return false
}

func NewInstance() *Instance {
	return &Instance{
		id:   id.XiD(),
		name: env.Env("APP_NAME", DefAppName),
		mode: InstMode(env.Env("APP_MODE", ModeDev.String())),
		// TENANT
	}
}

func (i *Instance) Valid() bool {
	return i.mode.Valid() &&
		i.name != "" &&
		!i.id.IsNil()
}

func (i *Instance) Name() string {
	return i.name
}

func (i *Instance) Mode() string {
	return i.mode.String()
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
