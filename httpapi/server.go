package httpapi

import (
	"fmt"
	"platform-go-challenge/domain"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	domain domain.IDomain
	port   int
	secret string
}

func NewServer(domain domain.IDomain, port int, secret string) *Server {
	return &Server{domain: domain, port: port, secret: secret}
}

func (s *Server) Run() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	a := e.Group("/auth")
	a.POST("/login", s.loginUserHandler)
	a.POST("/users", s.createUserHandler)

	config := middleware.JWTConfig{
		Claims:     &JwtUserClaims{},
		SigningKey: []byte(s.secret),
	}

	r := e.Group("/api/v1")

	r.Use(middleware.JWTWithConfig(config))
	r.POST("/admin/:assetType", s.addAssetHandler)
	r.PUT("/admin/:assetType/:id", s.updateAssetHandler)
	r.DELETE("/admin/:assetType/:id", s.deleteAssetHandler)

	r.GET("/me", s.meHandler)
	r.GET("/me/favourites", s.listMyFavourites)

	r.GET("/assets", s.listAssetsHandler)

	r.GET("/:assetType/:id", s.getAssetHandler)
	r.PUT("/:assetType/:id/favourite", s.favourAnAssetHandler)
	r.DELETE("/:assetType/:id/favourite", s.favourAnAssetHandler)

	e.Logger.Fatal(e.Start(fmt.Sprint(":", s.port)))
}
