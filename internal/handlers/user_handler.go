package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"platform-go-challenge/internal/auth"
	"platform-go-challenge/internal/common"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/models"
	"platform-go-challenge/internal/storage"
)

type UserHandler struct {
	cfg   *config.Config
	users *storage.UserStore
	jwt   auth.JWTGenerator
}

func NewUserHandler(cfg *config.Config, users *storage.UserStore, jwt auth.JWTGenerator) *UserHandler {
	return &UserHandler{cfg: cfg, users: users, jwt: jwt}
}

func (uh *UserHandler) register(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		errorResponse := common.ErrorResponse{
			Error: "invalid JSON: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Create a new user model
	user, err := models.CreateUser(payload.Username, payload.Password, models.Role(payload.Role))
	if err != nil {
		errorResponse := common.ErrorResponse{
			Error: "failed to create user: " + err.Error(),
		}
		switch err.Error() {
		case common.ErrEmptyUsername, common.ErrEmptyPassword, common.ErrInvalidRole:
			w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		case common.ErrHashPassword:
			w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		default:
			w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Store in user storage
	if err := uh.users.AddUser(&user); err != nil {
		errorResponse := common.ErrorResponse{
			Error: "failed to add user: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Prepare response
	resp := struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		Role     string    `json:"role"`
	}{
		ID:       user.ID,
		Username: user.Username,
		Role:     string(user.UserRole),
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		errorResponse := common.ErrorResponse{
			Error: "failed to encode response: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
	}
}

func (uh *UserHandler) login(w http.ResponseWriter, r *http.Request) {
	// Parse the username and password from the request body
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		errorResponse := common.ErrorResponse{
			Error: "Invalid input: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest) // 400 Bad Request
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Check if the user exists
	user, exists := uh.users.GetUserByUsername(credentials.Username)
	if !exists {
		errorResponse := common.ErrorResponse{
			Error: "Invalid username or password",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized) // 401 Unauthorized
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Compare hashed password with provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		errorResponse := common.ErrorResponse{
			Error: "Invalid username or password",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized) // 401 Unauthorized
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Generate a JWT token for the user
	token, err := uh.jwt.Generate(uh.cfg, user.ID, user.UserRole)
	if err != nil {
		errorResponse := common.ErrorResponse{
			Error: "Error generating JWT token: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
		return
	}

	// Return the JWT token to the client
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{"token": token}); err != nil {
		// Return structured error response if encoding fails
		errorResponse := common.ErrorResponse{
			Error: "failed to encode response: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError) // 500 Internal Server Error
		if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
			http.Error(w, "Failed to encode error response", http.StatusInternalServerError)
		}
	}
}
