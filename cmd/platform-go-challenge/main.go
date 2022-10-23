package main

import (
	"log"
	"net/http"
	"time"

	"platform-go-challenge/config"
	http_asset "platform-go-challenge/http/asset"
	http_user "platform-go-challenge/http/user"
	"platform-go-challenge/internal/app/assets"
	"platform-go-challenge/internal/app/dashboards"
	"platform-go-challenge/internal/app/users"
	challenge_sql "platform-go-challenge/internal/repositories/sql"

	"github.com/joho/godotenv"

	"github.com/gorilla/mux"
)

func main() {
	log.Print("app starting ...")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.NewConfig()
	conn := setupConnections(cfg)
	services := setupServices(conn)
	r := createRouter(services, *cfg)
	for conn.mysql.Ping() != nil {
		log.Print("waiting for db ...")
		time.Sleep(2 * time.Second)
	}

	up, down := DBdata()
	defer func() {
		for _, q := range down {
			_, _ = conn.mysql.Exec(q)
		}
	}()
	for _, q := range up {
		_, err = conn.mysql.Exec(q)
		if err != nil {
			log.Print(err.Error())
			return
		}
	}

	_ = http.ListenAndServe(":5000", r)
}

func createRouter(services services, cfg config.Config) *mux.Router {
	r := mux.NewRouter()
	jwt := cfg.JWT

	uh := http_user.NewUserHandler(services.userService, jwt.JWTSecret)
	ah := http_asset.NewAssetHandler(services.assetService)

	// USER endpoints
	r.HandleFunc("/user/{user_id}/dashboard/list", uh.GetDashboard).Methods(http.MethodGet)
	r.HandleFunc("/user/{user_id}/token/get", uh.GetToken).Methods(http.MethodGet)

	// USER endpoints with token
	r.Handle("/user/{user_id}/dashboard/asset/add", uh.Middleware(http.HandlerFunc(uh.AddToDashboard))).Methods(http.MethodPut)
	r.Handle("/user/{user_id}/dashboard/asset/remove", uh.Middleware(http.HandlerFunc(uh.RemoveFromDashboard))).Methods(http.MethodPut)
	r.Handle("/user/{user_id}/dashboard/asset/edit", uh.Middleware(http.HandlerFunc(uh.EditDescription))).Methods(http.MethodPatch)

	// ASSET endpoints
	r.HandleFunc("/asset/list", ah.ListAssets).Methods(http.MethodGet)

	return r
}

type services struct {
	userService  *users.Service
	assetService *assets.Service
}
type connections struct {
	mysql *challenge_sql.Client
}

func setupConnections(config *config.Config) connections {
	db, _ := challenge_sql.New(config.MySQL.CDN())

	return connections{
		mysql: db,
	}
}

func setupServices(connections connections) services {
	dashboardsRepo := challenge_sql.NewDashboardsRepository(connections.mysql)
	dashboardService := dashboards.NewDashboardService(dashboardsRepo)
	userService := users.NewUsersService(dashboardService)

	assetRepo := challenge_sql.NewAssetsRepository(connections.mysql)
	assetService := assets.NewAssetService(assetRepo)

	return services{
		userService:  userService,
		assetService: assetService,
	}
}

func DBdata() ([]string, []string) {
	chartsUp := `INSERT INTO challenge.charts
	(id,	title, 		x_axis, 	y_axis, 	data) VALUES
	(1,		"chart_1",	"x axis",	"y axis", 	'{"data":1}'),
	(2,		"chart_2",	"x axis",	"y axis", 	'{"data":2}'),
	(3,		"chart_3",	"x axis",	"y axis", 	'{"data":3}');`

	audiencesUp := `INSERT INTO challenge.audiences
	(id, gender,   country_of_birth, 	age_group, 			hours_spent_online, number_of_purchases_last_month) VALUES
	(1,	 'male',   'gr', 				'young-adults', 	10, 				5),
	(2,	 'female', 'gr', 				'young-adults', 	10, 				5),
	(3,	 'male',   'de', 				'teenagers', 		25, 				2);
	`

	insightsUp := `INSERT INTO challenge.insights
	(id, title) VALUES
	(1,  '40% of millenials spend more than 3hours on social media daily'),
	(2,  '60% of teenagers spend more than 6hours on social media daily'),
	(3,  '10% of seniors spend less than 3hours on social media weekly');
	`

	usersUp := `INSERT INTO challenge.users
	(id, username) VALUES
	(1,'user_1'),
	(2,'user_2'),
	(3,'user_3');
	`

	dashboardsUp := `INSERT INTO challenge.dashboards
	(id,user_id) VALUES
	(1,	1),
	(2,	2),
	(3,	3);`

	d2aUp := `INSERT INTO challenge.dashboards2assets
	(dashboard_id, 	asset_id, 	asset_type, description) VALUES
	(1, 			1, 			'chart', 		'chart description 1'),
	(1, 			2, 			'chart', 		'chart description 2'),
	(1, 			1, 			'audience', 	'audience description 1'),
	(1, 			1, 			'insight', 		'insight description 1'),
	(2, 			1, 			'chart', 		'chart description 1 alt');
	`
	audiencesDown := "DELETE FROM  challenge.audiences;"
	chartsDown := "DELETE FROM  challenge.charts;"
	insightsDown := `DELETE FROM  challenge.insights;`
	usersDown := `DELETE FROM  challenge.users;`
	dashboardsDown := "DELETE FROM  challenge.dashboards;"
	d2aDown := "DELETE FROM  challenge.dashboards2assets;"

	up := []string{chartsUp, audiencesUp, insightsUp, usersUp, dashboardsUp, d2aUp}
	down := []string{chartsDown, audiencesDown, insightsDown, dashboardsDown, usersDown, d2aDown}

	return up, down
}
