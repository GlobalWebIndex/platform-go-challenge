package utils

import (
	"net/mail"
	"regexp"
	"strings"

	"github.com/algorand/go-algorand-sdk/v2/types"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

func IsPubKey(pubKey string) error {
	_, err := types.DecodeAddress(pubKey)
	return err
}

func IsEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
