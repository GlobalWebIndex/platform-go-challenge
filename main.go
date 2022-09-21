package main

import (
	"log"
	"platform-go-challenge/domain"
	"platform-go-challenge/httpapi"
	"platform-go-challenge/sqldb"
)

func main() {
	db, err := sqldb.NewDB("user", "user", "127.0.0.1:3306", "mydb")
	if err != nil {
		log.Fatal(err)
	}
	db.CreateTables()
	dom := domain.NewDomain(db)
	server := httpapi.NewServer(dom, 8000, "secret")
	server.Run()
}
