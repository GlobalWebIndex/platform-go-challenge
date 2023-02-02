package dto

import "strings"

type BriefUser struct {
	PubKey        string `json:"pub_key"`
	ChainId       int    `json:"chain_id"`
	UserId        string `json:"user_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	BirthDay      string `birthday`
	Gender        string `json:"gender"`
	Nationality   string `nationality`
	IdFingerprint string `id_fingerprint`
}

func (u *BriefUser) Valid() bool {
	return !(strings.TrimSpace(u.FirstName) == "" || strings.TrimSpace(u.LastName) == "" || strings.TrimSpace(u.UserId) == "" || strings.TrimSpace(u.Nationality) == "" || strings.TrimSpace(u.IdFingerprint) == "")
}
