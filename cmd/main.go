package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"

	"gwi_api/internal/app"
	"gwi_api/internal/config"
	"gwi_api/internal/repository"
	"gwi_api/internal/service"
	desc "gwi_api/pkg"
	appruntime "runtime"
)

func main() {
	// Ownify DB
	gwidb, err := repository.NewDB()
	if err != nil {
		log.Fatalf("[ERR] cannot create database %s", err)
		return
	}

	// preparing config file
	configPath, err := config.GetConfigPath()
	if err != nil {
		log.Fatalln("cannot determine config path")
	}
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalln("cannot read from a config")
	}

	//stripeKey := viper.Get("stripe.secret.key").(string)

	// JWT
	signedKeyJWT := viper.Get("jwt.signedKey").(string)
	tokenManager := service.NewTokenManager(signedKeyJWT)

	// Register all services
	//dbHandler := repository.NewDBHandler(db)
	gwiDB := repository.NewDBHandler(gwidb)
	//logDBHandler := repository.NewDBHandler(logdb)

	userService := service.NewUserService(gwiDB)
	authService := service.NewAuthService(gwiDB, tokenManager)

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
			authService,
			userService,
		))

		err = grpcServer.Serve(listener)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	// Starting HTTP server
	grpcMux := runtime.NewServeMux(httpOpts)

	err = desc.RegisterMicroserviceHandlerServer(context.Background(), grpcMux, app.NewMicroservice(
		authService,
		userService,
	))
	if err != nil {
		log.Println("cannot register this service")
	}

	// Adjust the path to Swagger UI files
	_, b, _, _ := appruntime.Caller(0)
	basepath := filepath.Dir(b)
	swaggerDir := filepath.Join(basepath, "..", "static", "swagger-ui")

	standardMux := http.NewServeMux()
	standardMux.Handle("/api/swagger-ui/", http.StripPrefix("/api/swagger-ui", http.FileServer(http.Dir(swaggerDir))))
	standardMux.HandleFunc("/api/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(basepath, "..", "pkg", "microservice.swagger.json"))
	})

	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/swagger-ui") || r.URL.Path == "/api/swagger.json" {
			standardMux.ServeHTTP(w, r)
			return
		}

		grpcMux.ServeHTTP(w, r)
	})

	log.Fatalln(http.ListenAndServe(":8901", addCORSHeaders(mux)))
}

// addCORSHeaders is a middleware function that adds the necessary CORS headers to the HTTP response.
func addCORSHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r)
		// Preflight request. Reply successfully:
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers to the response
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE, PATCH")
}

func rateLimitMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !limiter.Allow() {
				http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
