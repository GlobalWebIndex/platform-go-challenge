package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/SokratisChaimanas/platform-go-challenge/internal/app"
	"github.com/SokratisChaimanas/platform-go-challenge/internal/domain"
	"github.com/google/uuid"
)

// UserHandler serves user-related endpoints.
type UserHandler struct {
	userService *app.UserService
}

func NewUserHandler(userService *app.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Get godoc
// @Summary      Get user
// @Description  Returns a user by ID.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user_id  path      string  true  "User ID (UUID)"
// @Success      200      {object}  handlers.UserResponse
// @Failure      400      {object}  handlers.ErrorResponse
// @Failure      404      {object}  handlers.ErrorResponse
// @Failure      500      {object}  handlers.ErrorResponse
// @Router       /users/{user_id} [get]
func (handler *UserHandler) Get(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "application/json")

	userId, ok := parseUUIDParam(writer, req, "user_id")
	if !ok {
		WriteJsonError(writer, "invalid user_id", http.StatusBadRequest)
		return
	}

	u, err := handler.userService.Get(req.Context(), userId)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			WriteJsonError(writer, "user not found", http.StatusNotFound)
			return
		}

		WriteJsonError(writer, "internal error", http.StatusInternalServerError)
		return
	}

	type resp struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt string    `json:"created_at"`
	}

	_ = json.NewEncoder(writer).Encode(resp{
		ID:        u.ID,
		CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
	})
}
