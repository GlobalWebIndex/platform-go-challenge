package server

func (s *Server) registerRoutes() {
	usersRouter := s.router.PathPrefix("/users").Subrouter()
	usersRouter.HandleFunc("", GetAllUsersHandler).Methods("GET")
	usersRouter.HandleFunc("", AddUserHandler).Methods("POST")
	usersRouter.HandleFunc("/{user_id}", GetUserByIdHandler).Methods("GET")
	usersRouter.HandleFunc("/{user_id}/favorites", GetFavoriteForUserByIdHandler).Methods("GET")
	usersRouter.HandleFunc("/{user_id}/favorites", AddFavoriteForUserByIdHandler).Methods("POST")
}
