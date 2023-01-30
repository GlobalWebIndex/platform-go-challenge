package dto

import (
	"net/mail"
	"strings"
)

type Wallet struct {
	ChainId  int    `json:"chain_id"`
	UserRole string `json:"user_role"`
	Email    string `json:"email"`
	PubKey   string `json:"pub_addr"`
}

func (b *Wallet) Valid() bool {
	_, err := mail.ParseAddress(b.Email)
	return !(err != nil || strings.TrimSpace(b.PubKey) == "")
}
