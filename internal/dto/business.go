package dto

import (
	"net/mail"
	"strings"
)

type BriefBusiness struct {
	UserId       string `json:"user_id"`
	Email        string `json:"email" validate:"required,email"`
	Pin          string `json:"pin"`
	Recover      string `json:"recover"`
	Business     string `json:"business"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Location     string `json:"location"`
	PhoneNumber  string `json:"phone_number"`
	IsApproved   bool   `json:"is_approved"`
	IsSubscribed bool   `json:"is_subscribed"`
}

func (b *BriefBusiness) Valid() bool {
	_, err := mail.ParseAddress(b.Email)
	return !(err != nil || strings.TrimSpace(b.UserId) == "" || strings.TrimSpace(b.Email) == "" || strings.TrimSpace(b.Business) == "")
}
