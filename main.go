package main

import (
	"platform-go-challenge/api"
	"platform-go-challenge/db"
	"platform-go-challenge/graph"
	"platform-go-challenge/graph/resolvers"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func graphqlHandler(resolver graph.ResolverRoot) gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	db.InitDB()
	router := gin.Default()

	// GraphQL resolver
	resolver := &resolvers.Resolver{
		DB: db.GormDB,
	}

	//REST routes
	// Audience routes
	router.POST("/audience", api.CreateAudience)
	router.GET("/audiences", api.GetAudiences)
	router.GET("/audience/:id", api.GetAudience)
	router.PUT("/audience/:id", api.UpdateAudience)
	router.DELETE("/audience/:id", api.DeleteAudience)

	// Chart routes
	router.POST("/chart", api.CreateChart)
	router.GET("/charts", api.GetCharts)
	router.GET("/chart/:id", api.GetChart)
	router.PUT("/chart/:id", api.UpdateChart)
	router.DELETE("/chart/:id", api.DeleteChart)

	// Insight routes
	router.POST("/insight", api.CreateInsight)
	router.GET("/insights", api.GetInsights)
	router.GET("/insight/:id", api.GetInsight)
	router.PUT("/insight/:id", api.UpdateInsight)
	router.DELETE("/insight/:id", api.DeleteInsight)

	// UserStar routes
	router.POST("/userstar", api.CreateUserStar)
	router.GET("/userstars", api.GetUserStars)
	router.GET("/userstar/:id", api.GetUserStar)
	router.PUT("/userstar/:id", api.UpdateUserStar)
	router.DELETE("/userstar/:id", api.DeleteUserStar)

	// GraphQL routes
	router.POST("/graphql", graphqlHandler(resolver))
	router.GET("/graphql", playgroundHandler())

	router.Run(":8080")
}
