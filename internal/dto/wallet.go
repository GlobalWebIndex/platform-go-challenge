package dto

import (
	"net/mail"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

type Wallet struct {
	ChainId  int    `json:"chain_id"`
	UserRole string `json:"user_role"`
	Email    string `json:"email"`
	PubKey   string `json:"pub_addr"`
	UserId   string `json:"user_id"`
}

func (b *Wallet) Valid() bool {
	_, err := mail.ParseAddress(b.Email)
	_, addressErr := types.DecodeAddress(b.PubKey)
	return !(err != nil || strings.TrimSpace(b.PubKey) == "" || addressErr != nil || strings.TrimSpace(b.UserId) == "")
}
