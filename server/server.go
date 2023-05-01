package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
	logger     *log.Logger
}

func NewServer(
	httpServer *http.Server,
	router *mux.Router,
	logger *log.Logger,
) *Server {
	return &Server{
		httpServer: httpServer,
		router:     router,
		logger:     logger,
	}
}

func (s *Server) Start() {
	s.registerRoutes()
	if err := http.ListenAndServe(":8080", s.router); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Server is running at port 8080")
}
