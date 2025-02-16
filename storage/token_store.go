package storage

import (
	"sync"
)

var tokenStore = struct {
	sync.Mutex
	tokens map[string]string
}{tokens: make(map[string]string)}

func StoreToken(userID, token string) {
	tokenStore.Lock()
	defer tokenStore.Unlock()
	tokenStore.tokens[userID] = token
}

func GetToken(userID string) (string, bool) {
	tokenStore.Lock()
	defer tokenStore.Unlock()
	token, exists := tokenStore.tokens[userID]
	return token, exists
}
