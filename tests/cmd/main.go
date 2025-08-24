/*
 * GWI API
 */

package main

import (
	"context"
	"gwiEng/internal/api/router"
	"gwiEng/internal/config"
	"log"
	"net/http"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("Main is starting...")

	log.Printf("GWI listening in 0.0.0.0:8082")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	mongoClient := config.InitMongoClient(cfg.Database.MongoURI)
	defer func(){
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Printf("Error disconnecting MongoDB: %s", err)
		}
	}()
	logger := config.InitLogger(cfg.Logs.LogsDir)

	router := router.NewRouter(mongoClient, logger)

	//log.Fatal(http.ListenAndServe(":8082", router))
	if cfg.Server.TLS{
		log.Fatal(http.ListenAndServeTLS(":8082", "app/cert/server.crt", "app/cert/server.key", router))
	}else{
		log.Fatal(http.ListenAndServe(":8082", router))
	}
}






