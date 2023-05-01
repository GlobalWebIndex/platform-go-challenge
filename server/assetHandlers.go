package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"platform-go-challenge/models"
	"platform-go-challenge/repository"

	"github.com/gorilla/mux"
)

func GetFavoriteForUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["user_id"]

	fmt.Println("GetFavoriteForUserById endpoint called for user with id : " + userId)

	if userId == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User with that Id was not found")
		return
	}

	assets, sqlErr := repository.GetFavoriteByUserId(userId)

	if sqlErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonDoc, err := json.Marshal(assets)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonDoc)
}

func AddFavoriteForUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["user_id"]

	var newFavorite models.Asset[any]
	err := json.NewDecoder(r.Body).Decode(&newFavorite)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body")
		return
	}

	fmt.Printf("AddFavoriteForUserById endpoint called for user with id %s and %+v\n", userId, newFavorite)

	// if favoriteType := newFavorite.Type; favoriteType == "Chart" {
	// 	var chart models.Chart
	// 	err := json.Unmarshal(newFavorite.Data, &chart)
	// 	fmt.Printf("This is a chart with X title : %s, Y title : %s and data : %s", chart.XAxisTitle)
	// }

	if userId == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User with that Id was not found")
		return
	}

	w.WriteHeader(http.StatusOK)
}
