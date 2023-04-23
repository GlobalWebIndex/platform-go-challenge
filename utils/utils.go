package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"platform2.0-go-challenge/models"
)

func SendError(w http.ResponseWriter, status int, err models.Error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

const UniqueConstrainViolationString = "Unique constraint violated"

var InvalidRequest error
var UniqueConstrainViolation error

func NewInvalidRequest(message string) error {
	InvalidRequest = errors.New(message)
	return InvalidRequest
}

func NewUniqueConstrainViolation(message string) error {
	UniqueConstrainViolation = errors.New(message)
	return UniqueConstrainViolation
}
