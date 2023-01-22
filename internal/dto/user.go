package dto

type BriefUser struct {
	ChainId    int    `json:"chain_id"`
	Wallet     string `json:"wallet"`
	WalletType string `json:"wallet_type"`
}

