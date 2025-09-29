package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"favourite_assets/server/errors"
	"favourite_assets/server/authentication"
	"favourite_assets/server/services"

	"github.com/google/uuid"
)

type UserController struct {
	UserService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

// (admin-only)
func (c *UserController) CreateUserHandler(w http.ResponseWriter, r *http.Request) {

	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	if err := authentication.RequireRole(r.Context(), "admin"); err != nil {
		errors.WriteError(w, errors.ErrForbidden)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	user, err := c.UserService.CreateUser(req.Name, req.Email)
	if err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	errors.WriteJSON(w, http.StatusCreated, user)
}

// (all-roles)
func (c *UserController) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	idStr := r.URL.Query().Get("userId")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	user, err := c.UserService.GetUser(userID)
	if err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	errors.WriteJSON(w, http.StatusOK, user)
}

// (all-roles)
func (c *UserController) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	idStr := r.URL.Query().Get("userId")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	user, err := c.UserService.UpdateUser(userID, req.Name, req.Email)
	if err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	errors.WriteJSON(w, http.StatusOK, user)
}

// (admin-only)
func (c *UserController) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	if err := authentication.RequireRole(r.Context(), "admin"); err != nil {
		errors.WriteError(w, errors.ErrForbidden)
		return
	}

	idStr := r.URL.Query().Get("userId")
	userID, err := uuid.Parse(idStr)
	if err != nil {
		errors.WriteError(w, errors.ErrBadRequest)
		return
	}

	if err := c.UserService.DeleteUser(userID); err != nil {
		errors.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// (admin- only)
func (c *UserController) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	if authentication.GetUserInfo(r.Context()) == nil {
		errors.WriteError(w, errors.ErrUnauthorized)
		return
	}

	if err := authentication.RequireRole(r.Context(), "admin"); err != nil {
		errors.WriteError(w, errors.ErrForbidden)
		return
	}

	users := c.UserService.ListUsers()
	errors.WriteJSON(w, http.StatusOK, users)
}
