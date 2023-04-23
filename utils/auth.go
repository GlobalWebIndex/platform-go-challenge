package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"platform2.0-go-challenge/models"
)

var encryptionKey = "someKey"
var errorm models.Error

func GenerateJWT(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
	})

	tokenString, err := token.SignedString([]byte(encryptionKey + strconv.Itoa(int(user.ID))))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthorizationToken(endpoint http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["user_id"]

		if r.Header["Authorization"] != nil {

			token, err := jwt.Parse(r.Header["Authorization"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error on token")
				}

				return []byte(encryptionKey + id), nil
			})

			if err != nil {
				errorm.Message = "Unauthorized Access"
				SendError(w, http.StatusUnauthorized, errorm)
				return
			}

			if token.Valid {
				endpoint.ServeHTTP(w, r)
			} else {
				errorm.Message = "Unauthorized Access"
				SendError(w, http.StatusUnauthorized, errorm)
				return
			}
		} else {
			errorm.Message = "Unauthorized Access"
			SendError(w, http.StatusUnauthorized, errorm)
			return
		}
	})
}
