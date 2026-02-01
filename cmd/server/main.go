package main

import (
	"fmt"
	"net/http"

	"github.com/kongrillas/platform-go-challenge/internal/http/handlers"
	"github.com/kongrillas/platform-go-challenge/internal/http/router"
	"github.com/kongrillas/platform-go-challenge/internal/store"
)

func main() {
	s := store.NewMemoryStore()
	h := handlers.NewFavoritesHandler(s)

	app := router.New(h)

	fmt.Println("listening on :8080")
	_ = http.ListenAndServe(":8080", app)
}
