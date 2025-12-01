package controllers

import (
	"encoding/json"
	"net/http"
	"github.com/dintzeler/platform-go-challenge/services"
	"github.com/dintzeler/platform-go-challenge/validators"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	handleLogin(w, r)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	loginRequest, err := validators.ValidateLoginRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	token, status, err := services.AuthenticateUser(*loginRequest.Email, *loginRequest.Password)
	if err != nil {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}