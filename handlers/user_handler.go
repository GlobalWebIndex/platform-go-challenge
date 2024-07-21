package handlers

import (
	"encoding/json"
	"favourites/database"
	"favourites/middleware"
	"favourites/models"
	"favourites/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"io"
	"net/http"
	"strings"
	"time"
)

type UserHandler struct {
	service database.UserService
}

func NewUserHandler(service database.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(ctx *gin.Context) {
	users, err := h.service.GetAll(nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Found Users", "users": users})
}

// swagger:operation GET /users/{username} User GetByUsername
// # Get a user's details by username
// ---
// consumes:
//   - application/json
//
// produces:
//   - application/json
//
// parameters:
//   - in: path
//     name: username
//     schema:
//     type: string
//     required: true
//     description: The username of the user to get
//
// responses:
//
//		'200':
//		  description: The requested user's details
//		  schema:
//	   	     "$ref": "#/definitions/User"
//		'404':
//		  description: The requested user was not found
//		  schema:
//			type: string
func (h *UserHandler) GetByUsername(ctx *gin.Context) {
	username, _ := ctx.Params.Get("username")
	user, err := h.service.GetByUsername(ctx, username)
	if err != nil {
		if err.Error() == utils.ErrorNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	user.Password = ""
	ctx.JSON(http.StatusOK, gin.H{"success": "Found User", "user": user})
}

func (h *UserHandler) Add(ctx *gin.Context) {

	byteValue, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		fmt.Println(err)
	}
	var result *models.User
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	err = h.service.Create(ctx, result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (h *UserHandler) AddAll(ctx *gin.Context) {
	byteValue, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	var result []*models.User
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	err = h.service.CreateAll(ctx, result)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}
	ctx.Status(http.StatusCreated)

}

// swagger:operation POST /users/login User Login
//
// # User Login
//
// ---
// consumes:
//   - application/json
//
// produces:
//   - application/json
//
// responses:
//
//		'200':
//		  description: Successful login
//		  content:
//		    type:
//		      application/json:
//	             schema:
//	             type: string
//		'401':
//		  description: Unauthorized
//		  content:
//		    type:
//		      application/json:
//	             schema:
//	             type: string
func (h *UserHandler) Login(ctx *gin.Context) {

	var user *models.User

	byteValue, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		fmt.Println(err)
	}
	err = json.Unmarshal(byteValue, &user)

	existingUser, err := h.service.GetByUsername(ctx, user.Username)
	if err != nil {
		if err.Error() == utils.ErrorNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user does not exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	errHash := utils.CompareHashPassword(user.Password, existingUser.Password)

	if !errHash {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	expirationTime := time.Now().Add(10 * time.Minute)

	claims := &models.Claims{
		Role: existingUser.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   existingUser.Username,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(utils.JwtSecret))
	fmt.Println(tokenString)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("token", tokenString, 1000000, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"success": "user logged in"})
}

func (h *UserHandler) SignUp(ctx *gin.Context) {

	byteValue, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		fmt.Println(err)
	}
	var result *models.User
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred while creating user"})
		ctx.Abort()
		return
	}

	_, err = h.service.GetByUsername(ctx, result.Username)
	if err != nil && err.Error() != utils.ErrorNotFound {
		fmt.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else if err == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user already exists"})
		return
	}

	var errHash error
	result.Password, errHash = utils.GenerateHashPassword(result.Password)

	if errHash != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate password hash"})
		return
	}

	err = h.service.Create(ctx, result)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "user created", "user": result.Username})
}
func (h *UserHandler) LogOut(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"success": "user logged out"})
}

func (h *UserHandler) GetByRole(ctx *gin.Context) {
	role := ctx.GetString("role")
	if !strings.Contains(role, middleware.ROLE_SUFFIX) {
		fmt.Println(role)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized role"})
		return
	}
	username := strings.Split(role, middleware.ROLE_SUFFIX)[1]
	user, err := h.service.GetByUsername(ctx, username)
	if err != nil {
		if err.Error() == utils.ErrorNotFound {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": "Found User", "user": user})
}
