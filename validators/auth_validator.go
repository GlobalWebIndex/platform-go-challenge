package validators

import (
	"github.com/dintzeler/platform-go-challenge/customerrors"
	"net/http"
	"encoding/json"
)

type LoginRequest struct {
	Email *string `json:"email"`
	Password *string `json:"password"`
}

func ValidateLoginRequest(r *http.Request) (*LoginRequest, error) {
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		return nil, &customerrors.ValidationError{
			Message:   "Invalid request body",
			ErrorCode: "INVALID_REQUEST_BODY",
		}
	}

	err = validateEmail(loginRequest.Email)
	if err != nil {
		return nil, err
	}

	err = validatePassword(loginRequest.Password)
	if err != nil {
		return nil, err
	}

	return &loginRequest, nil
}

func validateEmail(email *string) error {
	if email == nil || *email == "" {
		return &customerrors.ValidationError{
			Message:   "Email is required",
			ErrorCode: "EMAIL_REQUIRED",
		}
	}
	return nil
}

func validatePassword(password *string) error {
	if password == nil || *password == "" {
		return &customerrors.ValidationError{
			Message:   "Password is required",
			ErrorCode: "PASSWORD_REQUIRED",
		}
	}
	return nil
}