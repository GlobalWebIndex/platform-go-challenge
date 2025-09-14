package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"platform-go-challenge/app/services"
	"platform-go-challenge/persistence"
	"platform-go-challenge/presentation"
	"syscall"
	"time"

	"github.com/arangodb/go-driver"
	arangohttp "github.com/arangodb/go-driver/http"
	"github.com/gorilla/mux"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	log.Println("Server is up and running...")

	log.Println("Waiting 10 seconds for the ArangoDB container to be ready...")
	time.Sleep(10 * time.Second)

	arangoURL := getEnv("ARANGO_URL", "")
	arangoUsername := getEnv("ARANGO_USERNAME", "root")
	arangoPassword := getEnv("ARANGO_PASSWORD", "")
	arangoDB := getEnv("ARANGO_DB", "")
	serverOnPort := getEnv("SERVE_ON_PORT", "8000")

	connection, err := arangohttp.NewConnection(
		arangohttp.ConnectionConfig{
			Endpoints: []string{arangoURL},
		},
	)
	if err != nil {
		log.Fatalf("Failed to create HTTP connection: %v", err)
	}

	client, err := driver.NewClient(
		driver.ClientConfig{
			Connection:     connection,
			Authentication: driver.BasicAuthentication(arangoUsername, arangoPassword),
		},
	)
	if err != nil {
		log.Fatalf("Failed to create the client: %v", err)
	}

	db, err := client.Database(context.Background(), arangoDB)
	if err != nil {
		log.Printf("Error getting database %s: %v", arangoDB, err)
		if driver.IsNotFoundGeneral(err) {
			db, err = client.CreateDatabase(context.Background(), arangoDB, nil)
			if err != nil {
				log.Fatalf("Failed to create database: %v", err)
			}
		}
		log.Println("Database 'platform-db' was created successfully")
	}

	collectionRepository := persistence.NewCollectionRepository(db)

	userCollection := collectionRepository.GetOrCreate("user", "document")
	assetCollection := collectionRepository.GetOrCreate("asset", "document")
	favouriteCollection := collectionRepository.GetOrCreate("favourite", "edge")

	userRepo := persistence.NewUserRepository(userCollection)
	assetRepo := persistence.NewAssetRepository(assetCollection)
	_, _ = userRepo.CreateUsers()
	_, _ = assetRepo.CreateAssets()

	r := mux.NewRouter()

	favouriteRepo := persistence.NewFavouriteRepository(db, map[string]driver.Collection{
		"asset":     assetCollection,
		"user":      userCollection,
		"favourite": favouriteCollection})
	favouriteService := services.NewFavouriteService(favouriteRepo, userRepo, assetRepo)

	favouriteRoutes := presentation.CreateRoutes(favouriteService)
	for routePath, routeDefinition := range favouriteRoutes {
		log.Printf("adding %s route with methods %v\n", routePath, routeDefinition.Methods)
		r.Handle(routePath, routeDefinition.HandlerFunc).Methods(routeDefinition.Methods...)
	}

	http.Handle("/", r)

	// Run our server in a goroutine so that it doesn't block.
	server := http.Server{
		Addr:    ":" + serverOnPort,
		Handler: r,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT)

	// Block until we receive our signal.
	<-stop
	log.Println("Shutting down server...")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server gracefully stopped.")
}
