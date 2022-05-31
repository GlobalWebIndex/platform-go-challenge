package main

import (
	"log"
	"platform-go-challenge/domain"
	"platform-go-challenge/httpapi"
	"platform-go-challenge/sqldb"
)

func main() {
	db, err := sqldb.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	dom := domain.NewDomain(db)
	server := httpapi.NewServer(dom, 8000, "secret")
	server.Run()
}
