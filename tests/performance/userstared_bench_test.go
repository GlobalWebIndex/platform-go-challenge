package performance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"platform-go-challenge/db"
	"platform-go-challenge/graph"
	"platform-go-challenge/graph/resolvers"
	"platform-go-challenge/models"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	benchDB     *gorm.DB
	benchRouter *gin.Engine
)

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

func setupBenchDB() *gorm.DB {
	testDBURL, isSet := os.LookupEnv("TEST_DB_URL")
	if !isSet {
		// Load .env file from project root
		_ = godotenv.Load("../../.env")

		testDBURL = os.Getenv("TEST_DB_URL")
	}

	database, err := gorm.Open(postgres.Open(testDBURL), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to benchmark database: %v", err))
	}

	database.AutoMigrate(
		&models.Audience{},
		&models.Chart{},
		&models.Insight{},
		&models.UserStar{},
	)

	return database
}

func setupBenchRouter(database *gorm.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	resolver := &resolvers.Resolver{
		DB: database,
	}

	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	router.POST("/graphql", func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	})

	return router
}

func seedBenchmarkData(database *gorm.DB, numUsers, itemsPerUser int) {
	// Clean existing data
	database.Exec("DELETE FROM user_favourites")
	database.Exec("DELETE FROM insights")
	database.Exec("DELETE FROM charts")
	database.Exec("DELETE FROM audiences")

	// Create a pool of assets
	audiences := make([]models.Audience, itemsPerUser)
	for i := 0; i < itemsPerUser; i++ {
		audiences[i] = models.Audience{
			Gender:        "Male",
			BirthCountry:  "USA",
			AgeGroup:      "25-34",
			DailyHours:    5,
			NoOfPurchases: 10,
		}
		database.Create(&audiences[i])
	}

	charts := make([]models.Chart, itemsPerUser)
	for i := 0; i < itemsPerUser; i++ {
		charts[i] = models.Chart{
			Title:      fmt.Sprintf("Chart %d", i),
			XAxisTitle: "X",
			YAxisTitle: "Y",
		}
		database.Create(&charts[i])
	}

	insights := make([]models.Insight, itemsPerUser)
	for i := 0; i < itemsPerUser; i++ {
		insights[i] = models.Insight{
			Text: fmt.Sprintf("Insight %d", i),
		}
		database.Create(&insights[i])
	}

	// Create favourites for each user
	for userID := 1; userID <= numUsers; userID++ {
		for i := 0; i < itemsPerUser; i++ {
			database.Create(&models.UserStar{
				UserID:  uint(userID),
				Type: models.AssetTypeAudience,
				AssetID: audiences[i].ID,
			})
			database.Create(&models.UserStar{
				UserID:  uint(userID),
				Type: models.AssetTypeChart,
				AssetID: charts[i].ID,
			})
			database.Create(&models.UserStar{
				UserID:  uint(userID),
				Type: models.AssetTypeInsight,
				AssetID: insights[i].ID,
			})
		}
	}
}

func executeGraphQLBench(b *testing.B, query string, variables map[string]interface{}) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		b.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest("POST", "/graphql", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	benchRouter.ServeHTTP(w, req)

	if w.Code != 200 {
		body, _ := io.ReadAll(w.Body)
		b.Fatalf("unexpected status code: got %d, want 200. Body: %s", w.Code, body)
	}
}

// BenchmarkUserStared_SmallDataset benchmarks with 5 items per user
func BenchmarkUserStared_SmallDataset(b *testing.B) {
	if benchDB == nil {
		benchDB = setupBenchDB()
		db.GormDB = benchDB
		benchRouter = setupBenchRouter(benchDB)
	}

	seedBenchmarkData(benchDB, 10, 5)

	query := `
		query GetUserStared($userID: ID!) {
			userstared(userID: $userID) {
				userid
				audience {
					id
					gender
					birthcountry
					agegroup
					dailyhours
					noofpurchases
				}
				chart {
					id
					title
					xaxistitle
					yaxistitle
				}
				insight {
					id
					text
				}
			}
		}
	`

	variables := map[string]interface{}{
		"userID": "1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executeGraphQLBench(b, query, variables)
	}
}

// BenchmarkUserStared_MediumDataset benchmarks with 50 items per user
func BenchmarkUserStared_MediumDataset(b *testing.B) {
	if benchDB == nil {
		benchDB = setupBenchDB()
		db.GormDB = benchDB
		benchRouter = setupBenchRouter(benchDB)
	}

	seedBenchmarkData(benchDB, 10, 50)

	query := `
		query GetUserStared($userID: ID!) {
			userstared(userID: $userID) {
				userid
				audience { id gender }
				chart { id title }
				insight { id text }
			}
		}
	`

	variables := map[string]interface{}{
		"userID": "1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executeGraphQLBench(b, query, variables)
	}
}

// BenchmarkUserStared_LargeDataset benchmarks with 200 items per user
func BenchmarkUserStared_LargeDataset(b *testing.B) {
	if benchDB == nil {
		benchDB = setupBenchDB()
		db.GormDB = benchDB
		benchRouter = setupBenchRouter(benchDB)
	}

	seedBenchmarkData(benchDB, 10, 200)

	query := `
		query GetUserStared($userID: ID!) {
			userstared(userID: $userID) {
				userid
				audience { id gender }
				chart { id title }
				insight { id text }
			}
		}
	`

	variables := map[string]interface{}{
		"userID": "1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executeGraphQLBench(b, query, variables)
	}
}

// BenchmarkUserStared_MinimalFields benchmarks with minimal field selection
func BenchmarkUserStared_MinimalFields(b *testing.B) {
	if benchDB == nil {
		benchDB = setupBenchDB()
		db.GormDB = benchDB
		benchRouter = setupBenchRouter(benchDB)
	}

	seedBenchmarkData(benchDB, 10, 50)

	query := `
		query GetUserStared($userID: ID!) {
			userstared(userID: $userID) {
				userid
				audience { id }
				chart { id }
				insight { id }
			}
		}
	`

	variables := map[string]interface{}{
		"userID": "1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executeGraphQLBench(b, query, variables)
	}
}

// BenchmarkUserStared_AllFields benchmarks with all fields selected
func BenchmarkUserStared_AllFields(b *testing.B) {
	if benchDB == nil {
		benchDB = setupBenchDB()
		db.GormDB = benchDB
		benchRouter = setupBenchRouter(benchDB)
	}

	seedBenchmarkData(benchDB, 10, 50)

	query := `
		query GetUserStared($userID: ID!) {
			userstared(userID: $userID) {
				userid
				audience {
					id
					gender
					birthcountry
					agegroup
					dailyhours
					noofpurchases
				}
				chart {
					id
					title
					xaxistitle
					yaxistitle
				}
				insight {
					id
					text
				}
			}
		}
	`

	variables := map[string]interface{}{
		"userID": "1",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executeGraphQLBench(b, query, variables)
	}
}

// BenchmarkUserStared_Parallel benchmarks concurrent requests
func BenchmarkUserStared_Parallel(b *testing.B) {
	if benchDB == nil {
		benchDB = setupBenchDB()
		db.GormDB = benchDB
		benchRouter = setupBenchRouter(benchDB)
	}

	seedBenchmarkData(benchDB, 100, 50)

	query := `
		query GetUserStared($userID: ID!) {
			userstared(userID: $userID) {
				userid
				audience { id gender }
				chart { id title }
				insight { id text }
			}
		}
	`

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		userID := 1
		for pb.Next() {
			variables := map[string]interface{}{
				"userID": fmt.Sprintf("%d", userID),
			}
			executeGraphQLBench(b, query, variables)
			userID++
			if userID > 100 {
				userID = 1
			}
		}
	})
}
