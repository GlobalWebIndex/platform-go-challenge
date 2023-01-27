package dto

import "strings"

type BriefBusiness struct {
	UserId      string `json:"user_id"`
	Email       string `json:"email"`
	Pin         string `json:"pin"`
	Recover     string `json:"recover"`
	Business    string `json:"business"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Location    string `json:"location"`
	PhoneNumber string `json:"phone_number"`
	IsApproved  bool   `json:"is_approved"`
}

func (b *BriefBusiness) Valid() bool {
	return !(strings.TrimSpace(b.UserId) == "" || strings.TrimSpace(b.Email) == "" || strings.TrimSpace(b.Pin) == "" || strings.TrimSpace(b.Business) == "")
}
