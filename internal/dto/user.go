package dto

import (
	"fmt"

	"github.com/algorand/go-algorand-sdk/v2/types"
	"github.com/go-playground/validator/v10"
)

type BriefUser struct {
	PubAddr       string `json:"pub_addr" validate:"required"`
	ChainId       int    `json:"chain_id"`
	UserId        string `json:"user_id" validate:"required"`
	FirstName     string `json:"first_name" validate:"required"`
	LastName      string `json:"last_name" validate:"required"`
	BirthDay      string `json:"birthday"`
	Gender        string `json:"gender"`
	Nationality   string `json:"nationality" validate:"required"`
	IdFingerprint string `json:"id_fingerprint" validate:"required,sha512"`
}

func (u *BriefUser) Valid() error {
	_, addressErr := types.DecodeAddress(u.PubAddr)
	if addressErr != nil {
		return addressErr
	}
	validate := validator.New()
	err := validate.Struct(u)
	return err
}

func (u *BriefUser) ValidMainInfo() error {
	_, addressErr := types.DecodeAddress(u.PubAddr)
	if addressErr != nil {
		return addressErr
	}

	if u.FirstName != "" && u.LastName != "" && u.BirthDay != "" && u.Gender != "" && u.Nationality != "" {
		return nil
	}
	return fmt.Errorf("[ERR] invalid user info: %s %s", u.FirstName, u.LastName)
}
