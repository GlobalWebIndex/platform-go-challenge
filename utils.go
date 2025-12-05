package main

import (
	"net/http"
	"strings"
	"fmt"
	"log"
)

func parseUserPath(r *http.Request) (string, string, error) {
	// manually parsing path params
	// net/http does not provide APIs for easily parsing path params. Hence, we need to do it manually
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	parts := strings.Split(path, "/")

	if len(parts) < 1 {
		log.Printf("Invalid path requested: %s", r.URL.Path)
		return "", "", fmt.Errorf("invalid path")
	}

	userID := parts[0]

	assetInfo := ""
	if len(parts) == 3 {
		assetInfo = parts[2]
	}

	log.Printf("Parsed userID=%s, assetInfo (type or id)=%s", userID, assetInfo)
	
	return userID, assetInfo, nil
}