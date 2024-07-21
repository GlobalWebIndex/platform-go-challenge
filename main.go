// Asset API.
//
// Documentation of an awesome API.
//
//	BasePath: /api/v1
//	Version: 1.0.0
//	Host: localhost:6060
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- jwt
//
// swagger:meta
package main

import (
	"favourites/database"
	"favourites/handlers"
	"favourites/middleware"
	"favourites/utils"
	"github.com/gin-gonic/gin"
	_ "github.com/go-openapi/swag"
	"net/http"
)

func init() {
	utils.Load()
}

func setupRouter() *gin.Engine {

	db := utils.GetDB()
	//utils.CloseClientDB()
	r := gin.Default()
	var (
		// Collections.
		chartsCollection     = db.Collection("charts")
		insightsCollection   = db.Collection("insights")
		audiencesCollection  = db.Collection("audiences")
		favouritesCollection = db.Collection("favourites")
		usersCollection      = db.Collection("users")

		// Services.
		chartService     = database.NewChartService(chartsCollection)
		insightService   = database.NewInsightService(insightsCollection)
		audienceService  = database.NewAudienceService(audiencesCollection)
		favouriteService = database.NewFavouriteService(favouritesCollection)
		assetService     = database.NewAssetService(chartService, insightService, audienceService)
		userService      = database.NewUserService(usersCollection)
	)
	appGroup := r.Group("/api/v1/")
	{

		usersGroup := appGroup.Group("/users")
		{
			userHandler := handlers.NewUserHandler(userService)
			usersGroup.POST("/login", userHandler.Login)
			usersGroup.POST("/logout", userHandler.LogOut)
			usersGroup.POST("/signup", userHandler.SignUp)
			usersGroup.GET("/:username", userHandler.GetByUsername)

			favouriteGroup := usersGroup.Group("/favourites")
			{
				favouriteGroup.Use(middleware.IsAuthorized())
				favouriteHandler := handlers.NewFavouriteHandler(favouriteService)
				favouriteGroup.GET("/", favouriteHandler.GetAll)
				favouriteGroup.GET("/:id", favouriteHandler.Get)
				favouriteGroup.POST("/add", favouriteHandler.Add)
				favouriteGroup.POST("/delete", favouriteHandler.Remove)
			}

		}

		adminGroup := appGroup.Group("/admin")
		{
			adminGroup.Use(middleware.IsAdmin())
			userHandler := handlers.NewUserHandler(userService)
			adminGroup.GET("/users", userHandler.GetAll)
			adminGroup.POST("/add-users-bulk", userHandler.AddAll)
			adminGroup.POST("/add-charts-bulk", handlers.NewChartHandler(chartService).AddAll)
			adminGroup.POST("/add-insights-bulk", handlers.NewInsightHandler(insightService).AddAll)
			adminGroup.POST("/add-audiences-bulk", handlers.NewAudienceHandler(audienceService).AddAll)
		}

		assetGroup := appGroup.Group("/assets")
		{
			assetHandler := handlers.NewAssetHandler(assetService)
			assetGroup.GET("/", assetHandler.GetAll)
		}

		insightGroup := appGroup.Group("/insights")
		{
			insightHandler := handlers.NewInsightHandler(insightService)
			insightGroup.GET("/", insightHandler.GetAll)
			insightGroup.GET("/:id", insightHandler.Get)
		}

		audienceGroup := appGroup.Group("/audiences")
		{
			audienceHandler := handlers.NewAudienceHandler(audienceService)
			audienceGroup.GET("/", audienceHandler.GetAll)
			audienceGroup.GET("/:id", audienceHandler.Get)
		}

		chartGroup := appGroup.Group("/charts")
		{
			chartHandler := handlers.NewChartHandler(chartService)
			chartGroup.GET("/", chartHandler.GetAll)
			chartGroup.GET("/:id", chartHandler.Get)
		}

	}

	// Ping Health Check
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": "pong"})
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in localhost:8080
	r.Run(":8080")
}
