package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"

	"ownify_api/internal/app"
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"
	"ownify_api/internal/service"
	desc "ownify_api/pkg"

	"ownify_api/internal/config"

	"github.com/ipinfo/go-ipinfo/ipinfo"
)

const (
	DB_OWNIFY   = "database.ownify.dbname"
	DB_USER_LOG = "database.log.dbname"
)

func main() {
	// Ownify DB
	ownifydb, err := repository.NewDB(DB_OWNIFY)
	if err != nil {
		log.Fatalf("[ERR] cannot create database %s", err)
		return
	}

	err = ownifydb.Ping()
	if err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}

	// User Log DB
	logdb, err := repository.NewDB(DB_USER_LOG)
	if err != nil {
		log.Fatalf("[ERR] cannot create database %s", err)
		return
	}

	err = logdb.Ping()
	if err != nil {
		log.Fatalf("cannot ping db: %v", err)
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

	stripeKey := viper.Get("stripe.secret.key").(string)

	// JWT
	signedKeyJWT := viper.Get("jwt.signedKey").(string)
	tokenManager := service.NewTokenManager(signedKeyJWT)

	// Register all services
	//dbHandler := repository.NewDBHandler(db)
	ownifyDBHandler := repository.NewDBHandler(ownifydb)
	logDBHandler := repository.NewDBHandler(logdb)

	wallet := repository.NewAlgoHandler()

	adminService := service.NewAdminService(ownifyDBHandler)
	userService := service.NewUserService(ownifyDBHandler)
	businessService := service.NewBusinessService(ownifyDBHandler)
	ownershipService := service.NewOwnershipService(ownifyDBHandler)
	authService := service.NewAuthService(ownifyDBHandler, tokenManager)
	productService := service.NewProductService(ownifyDBHandler)
	walletService := service.NewWalletService(wallet)
	notifyService := service.NewNotifyService()
	logService := service.NewloggerService(logDBHandler)
	licenseService := service.NewLicenseService(ownifyDBHandler)
	paymentService := service.NewPaymentService(ownifyDBHandler, stripeKey)

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
			logService,
			licenseService,
			paymentService,
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
		logService,
		licenseService,
		paymentService,
	))
	if err != nil {
		log.Println("cannot register this service")
	}
	//log.Fatalln(http.ListenAndServe(":8901", setUserLog(addCORSHeaders(mux), logService)))
	log.Fatalln(http.ListenAndServe(":8901", addCORSHeaders(mux)))
}

// addCORSHeaders is a middleware function that adds the necessary CORS headers to the HTTP response.
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
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
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

// Add User Session analysis
func setUserLog(handler http.Handler, logger service.LoggerService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logUserActivity(w, r, logger)
		handler.ServeHTTP(w, r)
	})
}

func logUserActivity(w http.ResponseWriter, r *http.Request, logger service.LoggerService) {

	validIP, _, _ := net.SplitHostPort(getIP(r))
	ip := net.ParseIP(validIP)
	details, err := ipinfo.GetInfo(ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if strings.Contains(r.RequestURI, "product/verify") || strings.Contains(r.RequestURI, "ownership/") {
		logger.LogUserActivity(dto.NewUserLogByIPInfo(details))
	}

	// Print the user's location information
	//fmt.Println(w, "Your IP address is %s\n", ip.String())
	fmt.Println(w, "Your location is %s, %s, %s\n", details.City, details.Region, details.Country)
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.RemoteAddr
		}
	}
	return ip
}
