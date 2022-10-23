package main

import (
	"fmt"
	"log"
	"net/http"

	"platform-go-challenge/config"
	http_asset "platform-go-challenge/http/asset"
	http_user "platform-go-challenge/http/user"
	"platform-go-challenge/internal/app/assets"
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
	services := setupServices(cfg)
	r := createRouter(services)
	fmt.Println(cfg.MySQL.CDN())
	_ = http.ListenAndServe(":5000", r)
}

func createRouter(services services) *mux.Router {
	r := mux.NewRouter()

	uh := http_user.NewUserHandler(services.userService)
	ah := http_asset.NewAssetHandler(services.assetService)

	r.HandleFunc("/user/dashboard/{user_id}", uh.GetDashboard).Methods(http.MethodGet)
	r.HandleFunc("/user/dashboard/{user_id}/asset/add", uh.AddToDashboard).Methods(http.MethodPut)
	r.HandleFunc("/user/dashboard/{user_id}/asset/remove", uh.RemoveFromDashboard).Methods(http.MethodPut)
	r.HandleFunc("/user/dashboard/{user_id}/asset/edit", uh.EditDescription).Methods(http.MethodPatch)

	r.HandleFunc("/asset/list", ah.ListAssets).Methods(http.MethodGet)

	return r
}

type services struct {
	userService  *users.Service
	assetService *assets.Service
}

func setupServices(config *config.Config) services {
	db, _ := challenge_sql.New(config.MySQL.CDN())
	dashboardsRepo := challenge_sql.NewDashboardsRepository(db)
	dashboardService := dashboards.NewDashboardService(dashboardsRepo)
	userService := users.NewUsersService(dashboardService)

	assetRepo := challenge_sql.NewAssetsRepository(db)
	assetService := assets.NewAssetService(assetRepo)

	return services{
		userService:  userService,
		assetService: assetService,
	}
}
