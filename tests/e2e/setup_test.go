package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	testDB     *gorm.DB
	testRouter *gin.Engine
	testServer *httptest.Server
)

// GraphQLRequest represents a GraphQL query request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   json.RawMessage `json:"data"`
	Errors []struct {
		Message string `json:"message"`
		Path    []any  `json:"path,omitempty"`
	} `json:"errors,omitempty"`
}

// SetupTestDB initializes a test database
func SetupTestDB() *gorm.DB {
	testDBURL, isSet := os.LookupEnv("TEST_DB_URL")
	if !isSet {
		// Load .env file from project root
		_ = godotenv.Load("../../.env")

		testDBURL = os.Getenv("TEST_DB_URL")
	}

	database, err := gorm.Open(postgres.Open(testDBURL), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect to test database: %v", err))
	}

	// Auto-migrate the schema
	database.AutoMigrate(
		&models.Audience{},
		&models.Chart{},
		&models.Insight{},
		&models.UserStar{},
	)

	return database
}

// SetupTestRouter creates a test Gin router with GraphQL endpoint
func SetupTestRouter(database *gorm.DB) *gin.Engine {
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

// ExecuteGraphQL executes a GraphQL query against the test server
func ExecuteGraphQL(t *testing.T, query string, variables map[string]interface{}) *GraphQLResponse {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("failed to marshal request: %v", err)
	}

	req := httptest.NewRequest("POST", "/graphql", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testRouter.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status code: got %d, want %d. Body: %s", w.Code, http.StatusOK, w.Body.String())
	}

	body, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var resp GraphQLResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	return &resp
}

// CleanupTestData removes all data from test tables
func CleanupTestData(database *gorm.DB) {
	database.Exec("DELETE FROM user_favourites")
	database.Exec("DELETE FROM insights")
	database.Exec("DELETE FROM charts")
	database.Exec("DELETE FROM audiences")
}

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	// Setup
	testDB = SetupTestDB()
	db.GormDB = testDB // Set global DB for any code that uses it
	testRouter = SetupTestRouter(testDB)

	// Run tests
	code := m.Run()

	// Teardown
	sqlDB, _ := testDB.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}

	os.Exit(code)
}

// SeedTestData creates sample data for testing
func SeedTestData(t *testing.T, database *gorm.DB) (audienceID, chartID, insightID uint) {
	// Create test audience
	audience := models.Audience{
		Gender:        "Male",
		BirthCountry:  "USA",
		AgeGroup:      "25-34",
		DailyHours:    5,
		NoOfPurchases: 10,
	}
	if err := database.Create(&audience).Error; err != nil {
		t.Fatalf("failed to create test audience: %v", err)
	}

	// Create test chart
	chart := models.Chart{
		Title:      "Sales Chart",
		XAxisTitle: "Month",
		YAxisTitle: "Revenue",
	}
	if err := database.Create(&chart).Error; err != nil {
		t.Fatalf("failed to create test chart: %v", err)
	}

	// Create test insight
	insight := models.Insight{
		Text: "Revenue increased by 20% this quarter",
	}
	if err := database.Create(&insight).Error; err != nil {
		t.Fatalf("failed to create test insight: %v", err)
	}

	return audience.ID, chart.ID, insight.ID
}
