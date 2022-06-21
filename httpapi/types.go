package httpapi

type InputUserCreation struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}
