package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	"ownify_api/internal/app"
	"ownify_api/internal/domain"
	"ownify_api/internal/repository"
	"ownify_api/internal/service"
	desc "ownify_api/pkg"
)

func main() {
	// DB
	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("[ERR] cannot create database %s", err)
		return
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	//prepare algorand client
	algoClient, algoIndexer, err := repository.NewAlgoClient(domain.MainNet)
	if err != nil {
		log.Fatalf("cannot initialize algorand client: %v", err)
	}

	algoTestClient, algoTestIndexer, err := repository.NewAlgoClient(domain.TestNet)
	if err != nil {
		log.Fatalf("cannot initialize algorand client: %v", err)
	}

	// preparing config file
	viper.AddConfigPath("../config")
	viper.SetConfigName("config")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalln("cannot read from a config")
	}

	// JWT
	signedKeyJWT := viper.Get("jwt.signedKey").(string)
	tokenManager := service.NewTokenManager(signedKeyJWT)

	// Register all services
	//dbHandler := repository.NewDBHandler(db)
	dbHandler := repository.NewDBHandler(db)
	wallet := repository.NewAlgoHandler(algoClient, algoIndexer, algoTestClient, algoTestIndexer)

	userService := service.NewUserService(dbHandler)
	authService := service.NewAuthService(dbHandler, tokenManager)
	productService := service.NewProductService(dbHandler)
	walletService := service.NewWalletService(wallet)

	// Interceptors
	grpcOpts := app.GrpcInterceptor()
	httpOpts := app.HttpInterceptor()

	// Starting gRPC server
	go func() {
		listener, err := net.Listen("tcp", ":8081")
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("start server")

		grpcServer := grpc.NewServer(grpcOpts)
		desc.RegisterMicroserviceServer(grpcServer, app.NewMicroservice(
			userService,
			authService,
			tokenManager,
			productService,
			walletService,
		))

		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// Starting HTTP server
	mux := runtime.NewServeMux(httpOpts)
	err = desc.RegisterMicroserviceHandlerServer(context.Background(), mux, app.NewMicroservice(
		userService,
		authService,
		tokenManager, productService, walletService))
	if err != nil {
		log.Println("cannot register this service")
	}
	log.Fatalln(http.ListenAndServe(":8080", mux))
}
