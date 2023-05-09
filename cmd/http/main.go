package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/loukaspe/platform-go-challenge/internal/repositories"
	"github.com/loukaspe/platform-go-challenge/pkg/logger"
	"github.com/loukaspe/platform-go-challenge/pkg/server"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func main() {
	getEnv()

	logger := logger.NewLogger(context.Background())
	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    os.Getenv("SERVER_ADDR"),
		Handler: router,
	}
	db := getDB()

	server := server.NewServer(db, router, httpServer, logger)

	server.Run()
}

func getDB() *gorm.DB {
	dbDsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s TimeZone=Europe/Athens",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)
	db, err := gorm.Open(postgres.Open(dbDsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to database: ", err)
	}

	// Drops added in order to start with clean DB on App start for
	// assessment reasons
	db.Migrator().DropTable("users")
	db.Migrator().DropTable("charts")
	db.Migrator().DropTable("insights")
	db.Migrator().DropTable("audiences")
	db.Migrator().DropTable("users_charts")
	db.Migrator().DropTable("users_insights")
	db.Migrator().DropTable("users_audiences")

	err = db.AutoMigrate(&repositories.User{})
	if err != nil {
		log.Fatal("cannot migrate user table")
	}

	return db
}

func getEnv() {
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	}
}
