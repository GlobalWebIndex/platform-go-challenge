package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"platform-go-challenge/endpoints"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/favAssets", endpoints.FavoriteAssets)
	mux.HandleFunc("/editAssets", endpoints.EditAssets)
	mux.HandleFunc("/deleteAssets", endpoints.DeleteAsset)
	mux.HandleFunc("/addAssets", endpoints.AddAsset)

	err := http.ListenAndServe(":3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
