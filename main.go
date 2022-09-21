package main

import (
	"log"
	"os"
	"platform-go-challenge/domain"
	"platform-go-challenge/httpapi"
	"platform-go-challenge/sqldb"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dbHost := os.Getenv("MYSQL_HOST")
	dbUsername := os.Getenv("MYSQL_USERNAME")
	dbPassword := os.Getenv("MYSQL_PASSWORD")
	dbDB := os.Getenv("MYSQL_DB")
	portStr := os.Getenv("PORT")
	port, _ := strconv.Atoi(portStr)
	secret := os.Getenv("JWT_SECRET")

	db, err := sqldb.NewDB(dbUsername, dbPassword, dbHost, dbDB)
	if err != nil {
		log.Fatal(err)
	}
	db.CreateTables()
	dom := domain.NewDomain(db)
	server := httpapi.NewServer(dom, port, secret)
	server.Run()
}
