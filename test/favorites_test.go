package tests

import (
    "github.com/dintzeler/platform-go-challenge/controllers"
	"github.com/dintzeler/platform-go-challenge/models"
	"github.com/dintzeler/platform-go-challenge/storage"
    "net/http"
    "net/http/httptest"
    "testing"
    "encoding/json"
    "strings"
    "os"
    "path/filepath"
    "runtime"
	"math/rand/v2"
	"fmt"
	"github.com/dintzeler/platform-go-challenge/middleware"
)

var authToken string

func init() {
    // Get the directory of this test file
    _, filename, _, _ := runtime.Caller(0)
    testDir := filepath.Dir(filename)
    
    // Go up one level to the project root
    projectRoot := filepath.Dir(testDir)
    
    // Set the absolute path to test_data.json
    dataFile := filepath.Join(projectRoot, "test_data.json")
    os.Setenv("DATA_FILE", dataFile)
    
    // Verify the file exists
    if _, err := os.Stat(dataFile); os.IsNotExist(err) {
        panic("test_data.json not found at: " + dataFile)
    }

	token, err := loginAndGetToken()
	if err != nil {
		panic("Failed to login and get token: " + err.Error())
	}
	authToken = token
}

func callFavoritesWithAuth(w http.ResponseWriter, req *http.Request) {
    middleware.JWTMiddleware(controllers.FavoritesHandler)(w, req)
}

func loginAndGetToken() (string, error) {
    loginBody := `{"email": "user1@example.com", "password": "hashed_password_1"}`
    req, _ := http.NewRequest("POST", "/login", strings.NewReader(loginBody))
    req.Header.Set("Content-Type", "application/json")
    
    w := httptest.NewRecorder()
    controllers.LoginHandler(w, req)
    
    var response map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &response)
    
    return response["token"].(string), nil
}

func TestGetFavorites(t *testing.T) {
	t.Run("Get favorites without token", func(t *testing.T) {
		testNoToken(t, "GET")
	})

	t.Run("Get favorites with invalid token", func(t *testing.T) {
		testInvalidToken(t, "GET")
	})

	t.Run("Get favorites successfully", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/favorites", nil)
		req.Header.Set("Authorization", "Bearer "+authToken)
 		callFavoritesWithAuth(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check JSON response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}

		if charts, ok := response["charts"].([]interface{}); !ok {
			t.Errorf("Response charts is missing or not an array")
		} else if len(charts) != 1 {
			t.Errorf("Expected 1 chart, got %d", len(charts))
		}else{
			chart := charts[0].(map[string]interface{})
			testChart(t, chart)
		}

		if insights, ok := response["insights"].([]interface{}); !ok {
			t.Errorf("Response insights is missing or not an array")
		} else if len(insights) != 1 {
			t.Errorf("Expected 1 insight, got %d", len(insights))
		}else{
			insight := insights[0].(map[string]interface{})
			testInsight(t, insight)
		}

		if audiences, ok := response["audiences"].([]interface{}); !ok {
			t.Errorf("Response audiences is missing or not an array")
		} else if len(audiences) != 2 {
			t.Errorf("Expected 2 audiences, got %d", len(audiences))
		}else{
			for _, a := range audiences {
				audience := a.(map[string]interface{})
				testAudience(t, audience)
			}
		}

	})	
}

func TestAddFavorite(t *testing.T) {
	t.Run("Get favorites without token", func(t *testing.T) {
		testNoToken(t, "POST")
	})
	
	t.Run("Add favorite with invalid token", func(t *testing.T) {
		testInvalidToken(t, "POST")
	})

	t.Run("Add favorite without providing request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/favorites", strings.NewReader(""))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST_BODY")
	})

	t.Run("Add favorite without providing asset_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_type": "chart"}`
		req, _ := http.NewRequest("POST", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "asset_id is a required field", "MISSING_ASSET_ID")
	})

	t.Run("Add favorite without providing asset_type", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": 1}`
		req, _ := http.NewRequest("POST", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "asset_type is a required field", "MISSING_ASSET_TYPE")
	})

	t.Run("Add favorite with negative asset_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": -5, "asset_type": "chart"}`
		req, _ := http.NewRequest("POST", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "asset_id must be a positive integer", "INVALID_ASSET_ID")
	})

	t.Run("Add favorite with invalid asset_type", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": 1, "asset_type": "invalid_type"}`
		req, _ := http.NewRequest("POST", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "Invalid asset_type (type must be 'chart', 'insight', or 'audience')", "INVALID_ASSET_TYPE")
	})

	t.Run("Add favorite for non-existing asset", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": 9999, "asset_type": "chart"}`
		req, _ := http.NewRequest("POST", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusNotFound, "Asset does not exist", "ASSET_NOT_FOUND")
	})

	t.Run("Add favorite that already exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": 1, "asset_type": "chart"}`
		req, _ := http.NewRequest("POST", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check JSON response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}

		if assetID := int(response["asset_id"].(float64)); assetID != 1 {
			t.Errorf("Response asset_id: %v", assetID)
		}
		if assetType := response["asset_type"]; assetType != "chart" {
			t.Errorf("Response asset_type: %v", assetType)
		}
		if userID := int(response["user_id"].(float64)); userID != 1 {
			t.Errorf("Response user_id: %v", userID)
		}
	})
	
	t.Run("Add favorite successfully", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": 2, "asset_type": "insight"}`
		req, _ := http.NewRequest("POST", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		// Check JSON response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}

		if assetID := int(response["asset_id"].(float64)); assetID != 2 {
			t.Errorf("Response asset_id: %v", assetID)
		}

		if assetType := response["asset_type"]; assetType != "insight" {
			t.Errorf("Response asset_type: %v", assetType)
		}

		if userID := int(response["user_id"].(float64)); userID != 1 {
			t.Errorf("Response user_id: %v", userID)
		}

		// Verify favorite is persisted in data file
		verifyFavoriteInDataFile(t, 1, 2, "insight")
	})
}

func TestUpdateFavorite(t *testing.T) {
	t.Run("Get favorites without token", func(t *testing.T) {
		testNoToken(t, "PATCH")
	})

	t.Run("Update favorite with invalid token", func(t *testing.T) {
		testInvalidToken(t, "PATCH")
	})

	t.Run("Update favorite without providing request body", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PATCH", "/favorites", strings.NewReader(""))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST_BODY")
	})

	t.Run("Update favorite without providing asset_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_type": "chart"}`
		req, _ := http.NewRequest("PATCH", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}

		// Check JSON response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}

		if message := response["message"]; message != "asset_id is a required field"  {
			t.Errorf("Response message: %v", message)
		}

		if errorCode := response["error_code"]; errorCode != "MISSING_ASSET_ID" {
			t.Errorf("Response error_code: %v", errorCode)
		}
	})

	t.Run("Update favorite without providing asset_type", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": 1}`
		req, _ := http.NewRequest("PATCH", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "asset_type is a required field", "MISSING_ASSET_TYPE")
	})

	t.Run("Update favorite with negative asset_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": -5, "asset_type": "chart"}`
		req, _ := http.NewRequest("PATCH", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "asset_id must be a positive integer", "INVALID_ASSET_ID")
	})

	t.Run("Update favorite with missing description", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": 1, "asset_type": "chart"}`
		req, _ := http.NewRequest("PATCH", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "description is a required field", "MISSING_DESCRIPTION")
	})

	t.Run("Update favorite with invalid asset_type", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": 1, "asset_type": "invalid_type", "description": "New description"}`
		req, _ := http.NewRequest("PATCH", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "Invalid asset_type (type must be 'chart', 'insight', or 'audience')", "INVALID_ASSET_TYPE")
	})

	t.Run("Update favorite for non-existing asset", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": 9999, "asset_type": "chart", "description": "New description"}`
		req, _ := http.NewRequest("PATCH", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusNotFound, "Asset does not exist", "ASSET_NOT_FOUND")
	})

	t.Run("Update asset that is not favorited", func(t *testing.T) {
		w := httptest.NewRecorder()
		reqBody := `{"asset_id": 2, "asset_type": "chart", "description": "Updated description"}`
		req, _ := http.NewRequest("PATCH", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusForbidden, "User does not have this asset as favorite", "ASSET_NOT_FAVORITE")
	})

	t.Run("Update favorite successfully", func(t *testing.T) {
		w := httptest.NewRecorder()
		UpdatedDescription := fmt.Sprintf("Updated description %s", generateRandomNDigitNumber(5))
		reqBody := fmt.Sprintf(`{"asset_id": 1, "asset_type": "chart", "description": "%s"}`, UpdatedDescription)
		req, _ := http.NewRequest("PATCH", "/favorites", strings.NewReader(reqBody))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Check JSON response
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		if err != nil {
			t.Fatalf("Failed to parse JSON response: %v", err)
		}

		if assetID := int(response["asset_id"].(float64)); assetID != 1 {
			t.Errorf("Response asset_id: %v", assetID)
		}

		if assetType := response["asset_type"]; assetType != "chart" {
			t.Errorf("Response asset_type: %v", assetType)
		}

		if userID := int(response["user_id"].(float64)); userID != 1 {
			t.Errorf("Response user_id: %v", userID)
		}

		if description := response["updated_description"]; description != UpdatedDescription {
			t.Errorf("Response description: %v", description)
		}

		//verify description changed
		verifyDescriptionChanged(t, 1, 1, "chart", UpdatedDescription)
	})
}

func TestDeleteFavorite(t *testing.T) {
	t.Run("Get favorites without token", func(t *testing.T) {
		testNoToken(t, "DELETE")
	})

	t.Run("Delete favorite with invalid Token", func(t *testing.T) {
		testInvalidToken(t, "DELETE")
	})

	t.Run("Delete favorite without providing asset_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/favorites", strings.NewReader(""))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "asset_id is a required field", "MISSING_ASSET_ID")
	})

	t.Run("Delete favorite without providing asset_type", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/favorites?asset_id=1", strings.NewReader(""))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "asset_type is a required field", "MISSING_ASSET_TYPE")
	})

	t.Run("Delete favorite with negative asset_id", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/favorites?asset_id=-5&asset_type=chart", strings.NewReader(""))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "asset_id must be a positive integer", "INVALID_ASSET_ID")
	})

	t.Run("Delete favorite with invalid asset_type", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/favorites?asset_id=1&asset_type=invalid_type", strings.NewReader(""))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusBadRequest, "Invalid asset_type (type must be 'chart', 'insight', or 'audience')", "INVALID_ASSET_TYPE")
	})

	t.Run("Delete favorite for non-existing asset", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/favorites?asset_id=9999&asset_type=chart", strings.NewReader(""))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusNotFound, "Asset does not exist", "ASSET_NOT_FOUND")
	})

	t.Run("Delete favorite that does not exist", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/favorites?asset_id=2&asset_type=chart", strings.NewReader(""))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		testErrorResponse(t, w, http.StatusNotFound, "Favorite does not exist", "FAVORITE_NOT_FOUND")
	})

	t.Run("Delete favorite successfully", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/favorites?asset_id=2&asset_type=insight", strings.NewReader(""))
		req.Header.Set("Authorization", "Bearer "+authToken)
		callFavoritesWithAuth(w, req)

		if w.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", w.Code)
		}

		// Verify favorite is removed from data file
		verifyFavoriteRemovedFromDataFile(t, 1, 2, "insight")
	})
}


func verifyFavoriteInDataFile(t *testing.T, userID, assetID int, assetType models.AssetType) {
    dataFile := os.Getenv("DATA_FILE")
    if dataFile == "" {
        t.Fatal("DATA_FILE environment variable not set")
    }

    dataStore, err := storage.LoadData(dataFile)
    if err != nil {
        t.Fatalf("Failed to load test data: %v", err)
    }

    for _, fav := range dataStore.Favorites {
        if fav.UserID == userID && fav.AssetID == assetID && fav.AssetType == assetType {
            return
        }
    }

    t.Errorf("Favorite not found in data file: user_id=%d, asset_id=%d, asset_type=%s", userID, assetID, assetType)
}

func verifyDescriptionChanged(t *testing.T, userID, assetID int, assetType, expectedDescription string) {
    dataFile := os.Getenv("DATA_FILE")
    if dataFile == "" {
        t.Fatal("DATA_FILE environment variable not set")
    }

    data, err := storage.LoadData(dataFile)
    if err != nil {
        t.Fatalf("Failed to load test data: %v", err)
    }

	for _, ch := range data.Charts {
		if ch.ID == assetID {
			description := ch.Description
			if description != expectedDescription {
				t.Errorf("Description mismatch: expected '%s', got '%s'", expectedDescription, description)
			}
			return
		}
	}

	t.Errorf("Asset not found in data file: asset_id=%d, asset_type=%s", assetID, assetType)
	
}

func verifyFavoriteRemovedFromDataFile(t *testing.T, userID, assetID int, assetType models.AssetType) {
	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		t.Fatal("DATA_FILE environment variable not set")
	}

	data, err := storage.LoadData(dataFile)
	if err != nil {
		t.Fatalf("Failed to load test data file: %v", err)
	}

	for _, fav := range data.Favorites {
		if fav.UserID == userID && fav.AssetID == assetID && fav.AssetType == assetType {
			t.Errorf("Favorite still found in data file after deletion: user_id=%d, asset_id=%d, asset_type=%s", userID, assetID, assetType)
			return
		}
	}
}

func testNoToken(t *testing.T, method string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, "/favorites", nil)
	callFavoritesWithAuth(w, req)

	testErrorResponse(t, w, http.StatusUnauthorized, "Missing authorization header", "MISSING_AUTH_HEADER")
}

func testInvalidToken(t *testing.T, method string) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, "/favorites", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	callFavoritesWithAuth(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}

	testErrorResponse(t, w, http.StatusUnauthorized, "Invalid or expired token", "INVALID_TOKEN")
}

func generateRandomNDigitNumber(n int) string {
	randomNDigitNumber := ""
	for i := 0; i < n; i++ {
		digit := rand.IntN(10)
		randomNDigitNumber += fmt.Sprintf("%d", digit)
	}
	return randomNDigitNumber
}

func testChart(t *testing.T, chart map[string]interface{}) {
	chartDatabase := findChartInTestData(t, int(chart["id"].(float64)))

	if chartID := int(chart["id"].(float64)); chartID != chartDatabase.ID {
		t.Errorf("Chart ID mismatch: expected %d, got %d", chartDatabase.ID, chartID)
	}
	if title := chart["title"].(string); title != chartDatabase.Title {
		t.Errorf("Chart Title mismatch: expected %s, got %s", chartDatabase.Title, title)
	}
	if description := chart["description"].(string); description != chartDatabase.Description {
		t.Errorf("Chart Description mismatch: expected %s, got %s", chartDatabase.Description, description)
	}
	for i, point := range chart["data"].([]interface{}) {
		dataPoint := point.(map[string]interface{})
		expectedPoint := chartDatabase.Data[i]
		if dataPoint["x"] != expectedPoint.X || dataPoint["y"] != expectedPoint.Y {
			t.Errorf("Chart data point mismatch at index %d: expected (%v, %v), got (%v, %v)", i, expectedPoint.X, expectedPoint.Y, dataPoint["x"], dataPoint["y"])
		}
	}
}

func testInsight(t *testing.T, insight map[string]interface{}) {
	insightDatabase := findInsightInTestData(t, int(insight["id"].(float64)))

	if insightID := int(insight["id"].(float64)); insightID != insightDatabase.ID {
		t.Errorf("Insight ID mismatch: expected %d, got %d", insightDatabase.ID, insightID)
	}

	if title := insight["text"].(string); title != insightDatabase.Text {
		t.Errorf("Insight Text mismatch: expected %s, got %s", insightDatabase.Text, title)
	}

	if content := insight["description"].(string); content != insightDatabase.Description {
		t.Errorf("Insight Description mismatch: expected %s, got %s", insightDatabase.Description, content)
	}
}

func testAudience(t *testing.T, audience map[string]interface{}) {
	audienceDatabase := findAudienceInTestData(t, int(audience["id"].(float64)))

	if audienceID := int(audience["id"].(float64)); audienceID != audienceDatabase.ID {
		t.Errorf("Audience ID mismatch: expected %d, got %d", audienceDatabase.ID, audienceID)
	}

	if gender := audience["gender"].(string); gender != audienceDatabase.Gender {
		t.Errorf("Audience Gender mismatch: expected %s, got %s", audienceDatabase.Gender, gender)
	}

	if age_group := audience["age_group"].(string); age_group != audienceDatabase.AgeGroup {
		t.Errorf("Audience Age group mismatch: expected %s, got %s", audienceDatabase.AgeGroup, age_group)
	}

	if hours_social_media_daily := int(audience["hours_social_media_daily"].(float64)); hours_social_media_daily != audienceDatabase.HoursSocialMediaDaily {
		t.Errorf("Audience HoursSocialMediaDaily mismatch: expected %d, got %d", audienceDatabase.HoursSocialMediaDaily, hours_social_media_daily)
	}

	if number_of_purchases_last_month := int(audience["number_of_purchases_last_month"].(float64)); number_of_purchases_last_month != audienceDatabase.NumberOfPurchasesLastMonth {
		t.Errorf("Audience NumberOfPurchasesLastMonth mismatch: expected %d, got %d", audienceDatabase.NumberOfPurchasesLastMonth, number_of_purchases_last_month)
	}

	if description := audience["description"].(string); description != audienceDatabase.Description {
		t.Errorf("Audience Description mismatch: expected %s, got %s", audienceDatabase.Description, description)
	}
}

func findChartInTestData(t *testing.T, chartID int) *models.Chart {
	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		t.Fatal("DATA_FILE environment variable not set")
	}


    dataStore, err := storage.LoadData(dataFile)
    if err != nil {
        t.Fatalf("Failed to load test data: %v", err)
    }

	for _, chart := range dataStore.Charts {
		if chart.ID == chartID {
			return &chart
		}
	}
	return nil
}

func findInsightInTestData(t *testing.T, insightID int) *models.Insight {
	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		t.Fatal("DATA_FILE environment variable not set")
	}

	dataStore, err := storage.LoadData(dataFile)
	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	for _, insight := range dataStore.Insights {
		if insight.ID == insightID {
			return &insight
		}
	}
	return nil
}

func findAudienceInTestData(t *testing.T, audienceID int) *models.Audience {
	dataFile := os.Getenv("DATA_FILE")
	if dataFile == "" {
		t.Fatal("DATA_FILE environment variable not set")
	}

	dataStore, err := storage.LoadData(dataFile)
	if err != nil {
		t.Fatalf("Failed to load test data: %v", err)
	}

	for _, audience := range dataStore.Audiences {
		if audience.ID == audienceID {
			return &audience
		}
	}

	return nil
}
