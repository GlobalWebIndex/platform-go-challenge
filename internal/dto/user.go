package dto

type BriefUser struct {
	ChainId     int    `json:"chain_id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	UserRole    string `json:"role"`
	PubKey      string `json:"pub_key"`
	WalletType  string `json:"wallet_type"`
}
