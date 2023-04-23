package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
	"platform2.0-go-challenge/controllers"
	"platform2.0-go-challenge/models"
	"platform2.0-go-challenge/utils"
)

func init() {
	gotenv.Load()
}

func InitializeData() {
	utils.ConnectPostgresDB()

	utils.DB.AutoMigrate(&models.User{}, &models.Chart{}, &models.Insight{}, &models.Audience{})

}

func InittializeRouter() {
	controller := controllers.Controller{}
	chartcontroller := controllers.ChartController{}
	insightcontroller := controllers.InsightController{}
	audiencecontroller := controllers.AudienceController{}
	usercontroller := controllers.UserController{}

	router := mux.NewRouter()

	router.HandleFunc("/api/assets/{user_id}", utils.AuthorizationToken(controller.GetUserAssets(utils.DB))).Methods("GET")
	router.HandleFunc("/api/assets/charts/{user_id}", utils.AuthorizationToken(chartcontroller.AddChart(utils.DB))).Methods("POST")
	router.HandleFunc("/api/assets/insights/{user_id}", utils.AuthorizationToken(insightcontroller.AddInsight(utils.DB))).Methods("POST")
	router.HandleFunc("/api/assets/audiences/{user_id}", utils.AuthorizationToken(audiencecontroller.AddAudience(utils.DB))).Methods("POST")
	//Add or remove from favourites
	router.HandleFunc("/api/assets/charts/{user_id}", utils.AuthorizationToken(chartcontroller.UpdateChart(utils.DB))).Methods("PUT")
	router.HandleFunc("/api/assets/insights/{user_id}", utils.AuthorizationToken(insightcontroller.UpdateInsight(utils.DB))).Methods("PUT")
	router.HandleFunc("/api/assets/audiences/{user_id}", utils.AuthorizationToken(audiencecontroller.UpdateAudience(utils.DB))).Methods("PUT")

	router.HandleFunc("/api/users/signup", usercontroller.AddUser(utils.DB)).Methods("POST")
	router.HandleFunc("/api/users/login", usercontroller.LoginUser(utils.DB)).Methods("POST")

	fmt.Println("Server is running at port 8000")
	log.Fatal(http.ListenAndServe(os.Getenv("PORT"), handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}))(router)))

}

func main() {
	InitializeData()
	InittializeRouter()
}
