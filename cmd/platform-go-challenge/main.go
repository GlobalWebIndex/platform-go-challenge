package main

import (
	"fmt"
	"log"
	"net/http"

	"platform-go-challenge/config"
	http_user "platform-go-challenge/http/user"
	"platform-go-challenge/internal/app/dashboards"
	"platform-go-challenge/internal/app/users"
	challenge_sql "platform-go-challenge/internal/repositories/sql"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("App starting....")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.NewConfig()
	r := createRouter(cfg)
	fmt.Println(cfg.MySQL.CDN())
	_ = http.ListenAndServe(":5000", r)
}

func createRouter(config *config.Config) *mux.Router {
	r := mux.NewRouter()

	db, _ := challenge_sql.New(config.MySQL.CDN())
	dashboardsRepo := challenge_sql.NewDashboardsRepository(db)
	dashboardService := dashboards.NewDashboardService(dashboardsRepo)
	userService := users.NewUsersService(dashboardService)

	h := http_user.NewUserHandler(userService)

	r.HandleFunc("/user/dashboard/{user_id}", h.GetDashboard).Methods(http.MethodGet)
	r.HandleFunc("/user/dashboard/{user_id}/asset/add", h.AddToDashboard).Methods(http.MethodPut)
	r.HandleFunc("/user/dashboard/{user_id}/asset/remove", h.RemoveFromDashboard).Methods(http.MethodPut)
	r.HandleFunc("/user/dashboard/{user_id}/asset/edit", h.EditDescription).Methods(http.MethodPatch)
	// 	// r.HandleFunc("/categories/{id}/read", h.ReadCategory).Methods(http.MethodGet)
	// 	// r.Handle("/categories/{id}/delete", h.Middleware(http.HandlerFunc(h.DeleteCategory))).Methods(http.MethodDelete)
	// 	// r.Handle("/categories/create", h.Middleware(http.HandlerFunc(h.CreateCategory))).Methods(http.MethodPost)
	// 	// r.Handle("/categories/update", h.Middleware(http.HandlerFunc(h.UpdateCategory))).Methods(http.MethodPut)

	// 	// r.HandleFunc("/products/list", h.ListProducts).Methods(http.MethodGet)
	// 	// r.HandleFunc("/products/{id}/read", h.ReadProduct).Methods(http.MethodGet)
	// 	// r.Handle("/products/{id}/delete", h.Middleware(http.HandlerFunc(h.DeleteProduct))).Methods(http.MethodDelete)
	// 	// r.Handle("/products/create", h.Middleware(http.HandlerFunc(h.CreateProduct))).Methods(http.MethodPost)
	// 	// r.Handle("/products/update", h.Middleware(http.HandlerFunc(h.UpdateProduct))).Methods(http.MethodPut)

	// 	r.HandleFunc("/token", h.GetToken).Methods(http.MethodGet)

	return r
}
