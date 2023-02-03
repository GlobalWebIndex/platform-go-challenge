package dto

import (
	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/go-playground/validator/v10"
)

type BriefUser struct {
	PubKey        string `json:"pub_key" validate:"required"`
	ChainId       int    `json:"chain_id" validate:"required"`
	UserId        string `json:"user_id" validate:"required"`
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name" validate:"required"`
	BirthDay      string `birthday`
	Gender        string `json:"gender"`
	Nationality   string `nationality,validate:"required"`
	IdFingerprint string `id_fingerprint,validate:"required,md5"`
}

func (u *BriefUser) Valid() bool {
	_, addressErr := types.DecodeAddress(u.PubKey)
	validate := validator.New()
	err := validate.Struct(u)
	return !(addressErr != nil || err != nil)
}
