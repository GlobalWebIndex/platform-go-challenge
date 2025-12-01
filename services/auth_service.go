package services

import (
	"github.com/dintzeler/platform-go-challenge/storage"
	"github.com/dintzeler/platform-go-challenge/customerrors"
	"net/http"
	"os"
	"github.com/dintzeler/platform-go-challenge/utils"

)

func AuthenticateUser(email string, password string) (string, int, error) {
	data, err := storage.LoadData(os.Getenv("DATA_FILE"))
	if err != nil {
		return "", http.StatusInternalServerError, &customerrors.ValidationError{
			Message:   "Error loading data",
			ErrorCode: "DATA_LOAD_ERROR",
		}
	}

	user := storage.GetUserByEmail(email, data.Users)
	if user == nil {
		return "", 404, &customerrors.ValidationError{
			Message: "User not found",
			ErrorCode: "USER_NOT_FOUND",
		}
	}

	if user.Password != password {
		return "", 401, &customerrors.ValidationError{
			Message: "Wrong password",
			ErrorCode: "WRONG_PASSWORD",
		}
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return "", 500, &customerrors.ValidationError{
			Message: "Failed to generate token",
			ErrorCode: "TOKEN_GENERATION_FAILED",
		}
	}

	return token, 200, nil
}