package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	errMalformedUser   = errors.New("malformed user")
)

var (
	//go:embed users.txt
	usersText string
	// users hold the users database
	users     []User

	//go:embed secret.txt
	configuredSecret []byte

	//go:embed assets.json
	assetsData []byte
	// assets are the static inmemory assets for the server
	assets     Assets
)

// decodeAssets decodes assetsData to assets.
// Should be called at init.
func decodeAssets() error {
	return json.NewDecoder(bytes.NewBuffer(assetsData)).Decode(&assets)
}

// newUser parses s, expected to be a pair "email password", and returns a User.
func newUser(s string) (User, error) {
	toks := strings.Fields(s)
	if len(toks) != 2 || !strings.Contains(toks[0], "@") {
		return User{}, fmt.Errorf("failed to parse for user %q: %w", s, errMalformedUser)
	}
	return User{toks[0], toks[1]}, nil
}

// parseUsers converts usersText to users.
// usersText contains lines of email, password pairs.
// Should be called at init.
func parseUsers() error {
	for _, s := range strings.Split(usersText, "\n") {
		if len(s) > 0 {
			user, err := newUser(s)
			if err != nil {
				return err
			}
			users = append(users, user)
		}
	}
	return nil
}

// findUser return the user for the credentials, if exists.
func findUser(email, passwd string) (User, bool) {
	ind := slices.Index(users, User{email, passwd})
	if ind == -1 {
		return User{}, false
	}
	return users[ind], true
}
