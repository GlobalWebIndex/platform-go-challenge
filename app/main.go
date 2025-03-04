package main

import (
	"github.com/gin-gonic/gin"

	"app/api"
	"app/configs"
	"app/repositories"
	entityServices "app/services"
	apiServices "app/api/services"
)

func main() {
    cfg := configs.NewConfigs()
	repository, _ := repositories.NewMongoRepository(cfg.DbClient, cfg.DbName)
	entityService := entityServices.NewGService(repository)
    apiService := apiServices.NewXService(entityService)

    router := gin.Default()
    handler := api.NewApiHandler(apiService)
    handler.SetupRoutes(router)

    router.Run(":8080")
}