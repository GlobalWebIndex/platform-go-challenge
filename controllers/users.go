package controllers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"platform2.0-go-challenge/models"
	Repository "platform2.0-go-challenge/repository"
	"platform2.0-go-challenge/utils"
)

type UserController struct{}

var users []models.User
var loginresp []models.Login

func (u UserController) AddUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		var UserID int
		var user models.User
		encryptedPassword, err := encrypt(user.Password)

		if err != nil {
			error.Message = "Error on password encryption"

			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" || user.Password == "" {
			error.Message = "Enter missing fields."
			utils.SendError(w, http.StatusBadRequest, error) //400
			return
		}
		user.Password = encryptedPassword
		userRepo := Repository.UserRepository{}
		UserID, errcreated := userRepo.AddUser(db, user)

		if errcreated != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, UserID)
	}
}

func (u UserController) LoginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var error models.Error
		var user models.User

		json.NewDecoder(r.Body).Decode(&user)

		userRepo := Repository.UserRepository{}
		row, err := userRepo.LoginUser(db, user, user.Email)

		if err != nil {
			error.Message = "Error retrieving user"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		token, err := utils.GenerateJWT(row)

		if err != nil {
			error.Message = "Error retrieving token"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}
		result := models.Login{ID: row.ID, Token: token}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, result)
	}
}

func encrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}

	return string(hash), err
}
