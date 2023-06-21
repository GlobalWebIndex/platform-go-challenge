package fake

import (
	"fmt"

	userpb "x-gwi/proto/core/_user/v1"
	"x-gwi/service"
)

// userpb

func keyUser(seed int64) string {
	if seed == 0 {
		return ""
	}

	return fmt.Sprintf("fk-u:%d", seed)
}

// UserCore
func FakeUserCore(seed int64) *userpb.UserCore {
	if seed == 0 {
		return &userpb.UserCore{}
	}
	return &userpb.UserCore{
		Qid:          FakeShareQID(seed, service.NameUser, keyUser(seed)),
		Md:           FakeMetaDescription(seed),
		Mmd:          FakeMetaMultiDescription(seed),
		BasicAccount: FakeBasicAccount(seed),
	}
}

// BasicAccount
func FakeBasicAccount(seed int64) *userpb.BasicAccount {
	if seed == 0 {
		return &userpb.BasicAccount{}
	}
	return &userpb.BasicAccount{
		Username: keyUser(seed),
	}
}
