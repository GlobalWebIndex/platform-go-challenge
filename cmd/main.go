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
	wallet := repository.NewAlgoHandler()

	adminService := service.NewAdminService(dbHandler)
	userService := service.NewUserService(dbHandler)
	businessService := service.NewBusinessService(dbHandler)
	ownershipService := service.NewOwnershipService(dbHandler)
	authService := service.NewAuthService(dbHandler, tokenManager)
	productService := service.NewProductService(dbHandler)
	walletService := service.NewWalletService(wallet)
	notifyService := service.NewNotifyService()

	// Interceptors
	grpcOpts := app.GrpcInterceptor()
	httpOpts := app.HttpInterceptor()

	// Starting gRPC server
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:8900")
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("start server")

		grpcServer := grpc.NewServer(grpcOpts)
		desc.RegisterMicroserviceServer(grpcServer, app.NewMicroservice(
			adminService,
			userService,
			businessService,
			ownershipService,
			authService,
			tokenManager,
			productService,
			walletService,
			notifyService,
		))

		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// Starting HTTP server
	mux := runtime.NewServeMux(httpOpts)
	err = desc.RegisterMicroserviceHandlerServer(context.Background(), mux, app.NewMicroservice(
		adminService,
		userService,
		businessService,
		ownershipService,
		authService,
		tokenManager,
		productService,
		walletService,
		notifyService,
	))
	if err != nil {
		log.Println("cannot register this service")
	}
	log.Fatalln(http.ListenAndServe("0.0.0.0:8901", addCORSHeaders(mux)))
}

// addCORSHeaders is a middleware function that adds the necessary CORS headers to the HTTP response.
func addCORSHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		handler.ServeHTTP(w, r)
	})
}
func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
}
