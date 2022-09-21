package httpapi

import (
	"net/http"
	"platform-go-challenge/domain"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (s *Server) createUserHandler(c echo.Context) error {
	in := RequestUserCreation{}
	err := c.Bind(&in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}
	u := domain.User{}
	u.Username = in.Username
	u.Password = in.Password
	u.IsAdmin = in.IsAdmin
	_, err = s.domain.CreateUser(c.Request().Context(), u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, ResponseStatus{
		Status: SuccessStatus,
	})
}

func (s *Server) loginUserHandler(c echo.Context) error {
	in := RequestUserLogin{}
	err := c.Bind(&in)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}

	u, err := s.domain.LoginUser(c.Request().Context(), domain.LoginCredentials{
		Username: in.Username,
		Password: in.Password,
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}

	if in.ExpiresInMinutes <= 0 {
		in.ExpiresInMinutes = 60
	}

	expiresAt := time.Now().Add(time.Minute * time.Duration(in.ExpiresInMinutes)).Unix()
	claims := &JwtUserClaims{
		u.ID,
		in.Username,
		u.IsAdmin,
		jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, ResponseLogin{
		Status:    SuccessStatus,
		Token:     &t,
		ExpiresAt: &expiresAt,
		Username:  &in.Username,
	})
}

func (s *Server) meHandler(c echo.Context) error {
	user, err := getUserDomain(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
			"error":  err.Error(),
		})
	}
	uj := fromUserDomainToUserJson(*user)
	return c.JSON(http.StatusOK, uj)
}
