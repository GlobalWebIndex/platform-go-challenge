package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"platform-go-challenge/models"
	"platform-go-challenge/repository"

	"github.com/gorilla/mux"
)

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("GetAllUsers endpoint called")

	userList, err := repository.GetAllUsers()

	jsonDoc, err := json.Marshal(userList)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonDoc)

}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if newUser.Username == "" || newUser.Email == "" || newUser.Password == "" {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Missing user info, user cannot be created")
		return
	}

	fmt.Println("AddUser endpoint called for user with username : " + newUser.Username)

	// Add user to keycloak and if success store user in db
	if err := repository.AddUser(newUser); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["user_id"]

	fmt.Println("GetUserById endpoint called for user with id : " + userId)

	if userId == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User with that Id was not found")
		return
	}

	jsonDoc, err := json.Marshal(&models.User{
		Id:       userId,
		Username: "Bill",
		Email:    "pagalosb@gmail.com",
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonDoc)

}
