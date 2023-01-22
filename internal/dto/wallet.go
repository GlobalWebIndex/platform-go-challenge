package dto

type Wallet struct {
	ChainId  int    `json:"chain_id"`
	UserRole string `json:"user_role"`
	Email    string `json:"email"`
	PubKey   string `json:"pub_addr"`
}
