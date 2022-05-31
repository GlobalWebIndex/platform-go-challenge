package httpapi

import (
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
	return &Server{domain: domain, port: port}
}

func (s *Server) Run() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	a := e.Group("/auth")
	a.GET("/login", s.loginUser)
	a.POST("/users", s.createUser)

	r := e.Group("/api/v1")
	r.Use(middleware.JWT([]byte(s.secret)))
	r.PUT("/users/assets/:id", s.favourAnAsset)
	r.GET("/assets", s.listAssets)
	r.POST("/assets", s.addAsset)
	r.GET("/assets/:id", s.getAsset)
	r.PUT("/assets/:id", s.updateAsset)
	r.DELETE("/assets/:id", s.deleteAsset)

	e.Logger.Fatal(e.Start(":8080"))
}
