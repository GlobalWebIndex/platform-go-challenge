package repository

import (
	"database/sql"
	"fmt"
	"platform-go-challenge/integrations"
	"platform-go-challenge/models"
)

func GetAllUsers() ([]models.User, error) {
	db := integrations.DB

	qry := "SELECT user_id, username, email FROM tUser"

	users, err := db.Query(qry)

	if err != nil {
		// handle this error
		panic(err)
	}

	usersFound := make([]models.User, 0)

	for users.Next() {
		var user models.User
		err = users.Scan(&user.Id, &user.Username, &user.Email)
		if err != nil {
			// handle this error
			panic(err)
		}
		fmt.Println(user.Id, user.Username, user.Email)
		usersFound = append(usersFound, user)
	}

	return usersFound, nil
}

func GetUserById(id string) (models.User, error) {
	db := integrations.DB

	qry := "SELECT user_id, username, email FROM tUser where user_id = $1"

	user := db.QueryRow(qry, id)

	userFound := models.User{}

	var user_id string
	var username string
	var email string

	switch err := user.Scan(&user_id, &username, &email); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(user_id, username, email)
		userFound = models.User{
			Id:       user_id,
			Username: username,
			Email:    email,
		}
	default:
		panic(err)
	}

	return userFound, nil
}

func AddUser(user models.User) error {
	db := integrations.DB

	qry := "INSERT INTO tUser (user_id, username, email) VALUES ($1, $2, $3)"

	_, err := db.Exec(qry, user.Id, user.Username, user.Email)

	return err
}
