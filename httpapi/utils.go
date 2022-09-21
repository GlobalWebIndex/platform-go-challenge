package httpapi

import (
	"errors"
	"log"
	"platform-go-challenge/domain"
	"platform-go-challenge/sqldb"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func setupSuite(tb testing.TB) (*Server, func(tb testing.TB)) {
	db, err := sqldb.NewDB("user", "user", "127.0.0.1:3306", "mydb")
	if err != nil {
		log.Fatal(err)
	}
	dom := domain.NewDomain(db)
	server := NewServer(dom, 8000, "secret")

	db.DropTablesIfExist()
	sqldb, _ := db.GormDB().DB()
	db.CreateTables()
	return server, func(tb testing.TB) {
		db.DropTablesIfExist()
		sqldb.Close()
	}
}

func getUserDomain(c echo.Context) (*domain.User, error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, errors.New("no key called user")
	}
	claims, ok := user.Claims.(*JwtUserClaims)
	if !ok {
		return nil, errors.New("no jwt claims")
	}
	return &domain.User{
		Username: claims.Username,
		ID:       claims.ID,
		IsAdmin:  claims.Admin,
	}, nil
}
