package models

type Login struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}
