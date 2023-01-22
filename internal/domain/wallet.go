package domain

type Account struct {
	Id         int    `json:"id"`
	SeedCipher string `json:"seed_cipher"`
}
type Wallet struct {
	UserRole Role      `json:"user_role"`
	UserId   string    `json:"user_id"`
	Account  []Account `json:"account"`
}
