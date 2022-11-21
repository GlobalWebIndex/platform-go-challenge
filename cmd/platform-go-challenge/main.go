package main

import (
	"context"
	"platform-go-challenge/internal/app/favouriteasset"
	"platform-go-challenge/internal/config"
	"platform-go-challenge/internal/infra/http/router/favourite"
	mongorepo "platform-go-challenge/internal/infra/repository/mongo"
	assetrepo "platform-go-challenge/internal/infra/repository/mongo/favouriteasset"
	redisrepo "platform-go-challenge/internal/infra/repository/redis"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	e := echo.New()

	if err := run(e); err != nil {
		e.Logger.Fatalf("Error: %v")
	}
}

func run(e *echo.Echo) error {
	ctx := context.Background()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	err := godotenv.Load("../../.env")
	if err != nil {
		e.Logger.Warnf("no .env file exists: %v", err)
	}

	// Setup app config
	cfg := config.New()

	// set up Mongo.
	mongoClient, err := setUpMongoDB(ctx, cfg)
	if err != nil {
		return err
	}

	// set up Redis.
	var redisClient *redis.Client
	if cfg.App.WithCache == "true" {
		redisClient, err = setUpRedisClient(ctx, cfg)
		if err != nil {
			return err
		}
	}

	// Repo creation.
	ar, err := getFavouriteAssetRepo(mongoClient, redisClient, cfg)
	if err != nil {
		return err
	}

	// Service creation.
	assetSrv := favouriteasset.NewService(ar)

	// Add service to router & add routes.
	favourite.NewRouter(assetSrv).AppendRoutes(e)

	// start the app.
	return e.Start(":" + cfg.App.Port)
}

func setUpMongoDB(ctx context.Context, cfg *config.Config) (*mongo.Database, error) {
	return mongorepo.New(ctx, mongorepo.Config(cfg.Mongo))
}

func setUpRedisClient(ctx context.Context, cfg *config.Config) (*redis.Client, error) {
	return redisrepo.New(ctx, redisrepo.Config(cfg.Redis))
}

// getFavouriteAssetRepo uses the decorator pattern in order to instantiate the asset repo with cache or not
// depending on the config.
func getFavouriteAssetRepo(md *mongo.Database, rc *redis.Client, cfg *config.Config) (assetrepo.Repo, error) {
	var ar assetrepo.Repo
	ar, err := assetrepo.NewRepository(md)
	if cfg.App.WithCache == "true" {
		ar, err = assetrepo.NewCacheRepository(ar, rc)
	}

	return ar, err
}
