package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	gwi "github.com/josedelrio85/platform-go-challenge/internal"
)

func main() {
	log.Println("GWI Go Coding Challenge starting...")

	repo, err := gwi.NewMemoryRepository()
	if err != nil {
		panic(err)
	}
	handler := gwi.Handler{
		Repo: repo,
	}

	router := mux.NewRouter()

	router.PathPrefix("/user/{userid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/fav/list").
		Handler(handler.GetFavsFromUser()).
		Methods(http.MethodGet)

	router.PathPrefix("/user/{userid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/fav/add").
		Handler(handler.AddNewFav()).
		Methods(http.MethodPost)

	router.PathPrefix("/user/{userid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/fav/edit").
		Handler(handler.UpdateFav()).
		Methods(http.MethodPut)

	router.PathPrefix("/user/{userid:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}}/fav/delete").
		Handler(handler.DeleteFav()).
		Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":4567", router))
}
