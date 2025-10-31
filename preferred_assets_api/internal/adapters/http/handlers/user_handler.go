package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/adapters/http/middleware"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/mapping"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/ports"
	"github.com/go-chi/chi/v5"
)

var _ ports.UserHandler = (*UserHandler)(nil)

type UserHandler struct {
	service ports.UserService
}

func NewUserHandler(s ports.UserService) *UserHandler {
	return &UserHandler{service: s}
}

// Create creates a new user
// @Summary Create a new user
// @Description Creates a new user account with the provided information
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User creation request"
// @Success 201 {object} dto.UserResponse "User created successfully"
// @Failure 400 {string} string "Invalid input data"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, ok := middleware.GetValidatedBody[dto.CreateUserRequest](r)
	if !ok {
		http.Error(w, "missing validated body", http.StatusBadRequest)
		return
	}

	usr, err := mapping.UserReqToDomain(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.CreateUser(usr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Get retrieves a user by ID
// @Summary Get a user by ID
// @Description Retrieves user details for the specified user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse "User found successfully"
// @Failure 400 {string} string "Invalid user ID"
// @Failure 404 {string} string "User not found"
// @Failure 405 {string} string "Method not allowed"
// @Security BearerAuth
// @Router /users/{id} [get]
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	u, err := h.service.GetUserByID(id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	usr := mapping.DomainToUserRes(*u)
	if err := json.NewEncoder(w).Encode(usr); err != nil {
		log.Fatal(err)
	}
}

// List retrieves all users
// @Summary Get all users
// @Description Retrieves a list of all users in the system
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {array} dto.UserResponse "List of users retrieved successfully"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, "error fetching users", http.StatusInternalServerError)
		return
	}

	usersResponse := mapping.UserReqToResponseList(users)
	if err := json.NewEncoder(w).Encode(usersResponse); err != nil {
		log.Fatal(err)
	}
}

// Delete removes a user by ID
// @Summary Delete a user
// @Description Permanently removes a user from the system
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 "User deleted successfully"
// @Failure 400 {string} string "Invalid user ID"
// @Failure 404 {string} string "User not found"
// @Failure 405 {string} string "Method not allowed"
// @Security BearerAuth
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteUser(id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
}

// Update modifies a user by ID
// @Summary Update a user
// @Description Updates user information for the specified user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body dto.UpdateUserRequest true "User update request"
// @Success 200 {string} string "User updated successfully"
// @Failure 400 {string} string "Invalid input data"
// @Failure 404 {string} string "User not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPatch {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	req, ok := middleware.GetValidatedBody[dto.UpdateUserRequest](r)
	if !ok {
		http.Error(w, "missing validated body", http.StatusBadRequest)
		return
	}

	existingUser, err := h.service.GetUserByID(id)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	updatedUser := mapping.UpdateReqToDomain(existingUser, req)
	if err := h.service.UpdateUser(*updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "updated"}); err != nil {
		http.Error(w, "Failed JSON serialization", http.StatusInternalServerError)
	}
}

// GetFavourites retrieves user favourites
// @Summary Get user favourites
// @Description Retrieves all favourite assets for the specified user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {array} dto.FavouriteResponse "User favourites retrieved successfully"
// @Failure 400 {string} string "Invalid user ID"
// @Failure 404 {string} string "Favourites not found"
// @Failure 405 {string} string "Method not allowed"
// @Failure 500 {string} string "Internal server error"
// @Security BearerAuth
// @Router /users/{id}/favourites [get]
func (h *UserHandler) GetFavourites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "missing user id", http.StatusBadRequest)
		return
	}

	favourites, err := h.service.GetFavouritesByUser(id)
	if err != nil {
		http.Error(w, "favourites not found", http.StatusNotFound)
		return
	}

	response := mapping.FavouritesToResponse(favourites)

	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("JSON marshaling error: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Favourites JSON response: %s", string(jsonBytes))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
