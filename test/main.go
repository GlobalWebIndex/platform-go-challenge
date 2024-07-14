package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var baseURL string

// Helper function to log requests and responses
func logRequestResponse(method, url string, requestBody, responseBody []byte, statusCode int) {
	log.Println("---------------------------------------------------------------------")
	log.Printf("Request: %s %s\n", method, url)
	if len(requestBody) > 0 {
		log.Printf("Request Body: %s\n", string(requestBody))
	}
	log.Printf("Response Status: %d\n", statusCode)
	if len(responseBody) > 0 {
		log.Printf("Response Body: %s\n", string(responseBody))
	}
	log.Println("---------------------------------------------------------------------")
}

func main() {
	log.Println("Start Tests")

	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")
	baseURL = fmt.Sprintf("http://%s:%s", appHost, appPort)

	log.Printf("Base URL: %s\n", baseURL)

	if err := EndToEndTest(); err != nil {
		log.Fatalf("End-to-end test failed: %v", err)
	}
	log.Println("All tests passed successfully!")
}

func EndToEndTest() error {
	// Wait for the server to start, needs fixing
	time.Sleep(5 * time.Second)

	// Ping the server
	if err := pingServer(); err != nil {
		return fmt.Errorf("ping test failed: %v", err)
	}

	log.Println("Ping Test success!")

	// Add a user
	userID, err := addUser()
	if err != nil {
		return fmt.Errorf("add user test failed: %v", err)
	}

	log.Println("Add user Test success!")

	// Add assets
	assetIDs, err := addAssets()
	if err != nil {
		return fmt.Errorf("add assets test failed: %v", err)
	}

	log.Println("Add asset Test success!")

	// Add favorites
	if err := addFavorites(userID, assetIDs); err != nil {
		return fmt.Errorf("add favorites test failed: %v", err)
	}

	log.Println("Add favorites Test success!")

	// Check detailed favorites
	if err := checkDetailedFavorites(userID); err != nil {
		return fmt.Errorf("check detailed favorites test failed: %v", err)
	}

	log.Println("Get favorites Test success!")

	if err := addAndDeleteAssetTest(); err != nil {
		return fmt.Errorf("add and Delete Asset test failed: %v", err)
	}

	log.Println("Delete assets Test success!")

	if err := modifyAssetTest(); err != nil {
		return fmt.Errorf("modify Asset test failed: %v", err)
	}

	log.Println("Modify assets Test success!")

	return nil
}

// Helper functions

func pingServer() error {
	resp, err := http.Get(baseURL + "/ping")
	if err != nil {
		return fmt.Errorf("failed to ping server: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logRequestResponse("GET", baseURL+"/ping", nil, body, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status OK, got %v", resp.Status)
	}

	bodyStr := string(body)
	parts := strings.SplitN(bodyStr, "-", 2)
	status := strings.TrimSpace(parts[0])

	if status != "OK" {
		return fmt.Errorf("expected body to start with 'OK', got %s", status)
	}

	return nil
}

func addUser() (int, error) {
	username := fmt.Sprintf("testuser_%d", rand.Intn(10000))
	user := map[string]string{
		"username": username,
		"email":    username + "@nero.com",
	}

	userJSON, _ := json.Marshal(user)

	resp, err := http.Post(baseURL+"/user/add", "application/json", bytes.NewBuffer(userJSON))
	if err != nil {
		return 0, fmt.Errorf("failed to add user: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logRequestResponse("POST", baseURL+"/user/add", userJSON, body, resp.StatusCode)

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("expected status Created, got %v", resp.Status)
	}

	var createdUser map[string]interface{}
	json.Unmarshal(body, &createdUser)
	return int(createdUser["id"].(float64)), nil
}

func addAssets() (map[string]int, error) {
	assetIDs := make(map[string]int)
	assetTypes := []string{"CHART", "INSIGHT", "AUDIENCE"}
	for _, assetType := range assetTypes {
		var payload map[string]interface{}
		switch assetType {
		case "CHART":
			payload = map[string]interface{}{
				"type": "CHART",
				"payload": map[string]interface{}{
					"title":        "Test Chart",
					"x_axis_title": "X Axis",
					"y_axis_title": "Y Axis",
					"data_points":  []map[string]float64{{"x": 1, "y": 2}, {"x": 2, "y": 4}},
				},
			}
		case "INSIGHT":
			payload = map[string]interface{}{
				"type":    "INSIGHT",
				"payload": map[string]string{"text": "Test Insight"},
			}
		case "AUDIENCE":
			payload = map[string]interface{}{
				"type": "AUDIENCE",
				"payload": map[string]interface{}{
					"gender":                      "Male",
					"birth_country":               "USA",
					"age_group":                   "25-34",
					"daily_hours_on_social_media": 2.5,
					"purchases_last_month":        3,
				},
			}
		}

		assetJSON, _ := json.Marshal(payload)
		resp, err := http.Post(baseURL+"/assets/add", "application/json", bytes.NewBuffer(assetJSON))
		if err != nil {
			return nil, fmt.Errorf("failed to add %s: %v", assetType, err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		logRequestResponse("POST", baseURL+"/assets/add", assetJSON, body, resp.StatusCode)

		if resp.StatusCode != http.StatusCreated {
			return nil, fmt.Errorf("expected status Created, got %v", resp.Status)
		}

		var createdAsset map[string]interface{}
		json.Unmarshal(body, &createdAsset)
		assetIDs[assetType] = int(createdAsset["id"].(float64))
	}
	return assetIDs, nil
}

func addFavorites(userID int, assetIDs map[string]int) error {
	favorites := make([]map[string]interface{}, 0)
	for assetType, assetID := range assetIDs {
		favorites = append(favorites, map[string]interface{}{
			"asset_id":   assetID,
			"asset_type": assetType,
		})
	}

	payload := map[string]interface{}{
		"user_id":   userID,
		"favorites": favorites,
	}

	favoritesJSON, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/user/favorite/add", "application/json", bytes.NewBuffer(favoritesJSON))
	if err != nil {
		return fmt.Errorf("failed to add favorites: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logRequestResponse("POST", baseURL+"/user/favorite/add", favoritesJSON, body, resp.StatusCode)

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("expected status Created, got %v", resp.Status)
	}

	return nil
}

func checkDetailedFavorites(userID int) error {

	url := fmt.Sprintf("%s/user/favorites?user_id=%d", baseURL, userID)
	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("failed to get detailed favorites: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	logRequestResponse("GET", url, nil, body, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status OK, got %v", resp.Status)
	}
	var detailedFavorites map[string]interface{}
	if err := json.Unmarshal(body, &detailedFavorites); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	log.Printf("detailedFavorites: %+v", detailedFavorites)

	favorites, ok := detailedFavorites["favorites"].([]interface{})
	if !ok {
		return fmt.Errorf("expected favorites to be an array, got %T", detailedFavorites["favorites"])
	}
	if len(favorites) != 3 {
		return fmt.Errorf("expected 3 favorites, got %d", len(favorites))
	}

	for i, fav := range favorites {
		log.Printf("Favorite %d: %+v", i, fav)
		favorite, ok := fav.(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected favorite to be a map, got %T", fav)
		}

		log.Printf("Favorite %d after type assertion: %+v", i, favorite)

		asset, ok := favorite["asset"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("expected asset to be a map, got %T", favorite["asset"])
		}

		assetType, ok := asset["asset_type"].(string)
		if !ok {
			return fmt.Errorf("expected asset_type to be a string, got %T", asset["asset_type"])
		}

		switch assetType {
		case "CHART":
			title, ok := asset["title"].(string)
			if !ok {
				return fmt.Errorf("expected chart title to be a string, got %T", asset["title"])
			}
			if title != "Test Chart" {
				return fmt.Errorf("expected chart title 'Test Chart', got %s", title)
			}
			// You can add more checks for other chart properties if needed
		case "INSIGHT":
			text, ok := asset["text"].(string)
			if !ok {
				return fmt.Errorf("expected insight text to be a string, got %T", asset["text"])
			}
			if text != "Test Insight" {
				return fmt.Errorf("expected insight text 'Test Insight', got %s", text)
			}
		case "AUDIENCE":
			gender, ok := asset["gender"].(string)
			if !ok {
				return fmt.Errorf("expected gender to be a string, got %T", asset["gender"])
			}
			birthCountry, ok := asset["birth_country"].(string)
			if !ok {
				return fmt.Errorf("expected birth_country to be a string, got %T", asset["birth_country"])
			}
			if gender != "Male" || birthCountry != "USA" {
				return fmt.Errorf("unexpected audience data: gender=%s, birth_country=%s", gender, birthCountry)
			}
			// You can add more checks for other audience properties if needed
		default:
			return fmt.Errorf("unexpected asset type: %s", assetType)
		}
	}

	log.Println("Detailed favorites check passed successfully!")
	return nil
}

func addAndDeleteAssetTest() error {
	log.Println("Starting Add and Delete Asset Test")

	assetTypes := []string{"CHART", "INSIGHT", "AUDIENCE"}

	for _, assetType := range assetTypes {
		log.Printf("Testing asset type: %s", assetType)

		// Add an asset
		assetID, err := addSingleAsset(assetType)
		if err != nil {
			return fmt.Errorf("failed to add %s asset: %v", assetType, err)
		}
		log.Printf("Successfully added %s asset with ID: %d", assetType, assetID)

		// Delete the asset
		if err := deleteAsset(assetID, assetType); err != nil {
			return fmt.Errorf("failed to delete %s asset: %v", assetType, err)
		}
		log.Printf("Successfully deleted %s asset with ID: %d", assetType, assetID)

		// Try to get the deleted asset (should fail)
		if err := getAsset(assetID, assetType); err == nil {
			return fmt.Errorf("%s asset %d still exists after deletion", assetType, assetID)
		} else {
			log.Printf("%s asset %d correctly not found after deletion", assetType, assetID)
		}

		log.Printf("Add and Delete %s Asset Test passed successfully!", assetType)
	}

	log.Println("All Add and Delete Asset Tests passed successfully!")
	return nil
}

func addSingleAsset(assetType string) (int, error) {
	var payload map[string]interface{}
	switch assetType {
	case "CHART":
		payload = map[string]interface{}{
			"type": "CHART",
			"payload": map[string]interface{}{
				"title":        "Test Chart",
				"x_axis_title": "X Axis",
				"y_axis_title": "Y Axis",
				"data_points":  []map[string]float64{{"x": 1, "y": 2}, {"x": 2, "y": 4}},
			},
		}
	case "INSIGHT":
		payload = map[string]interface{}{
			"type": "INSIGHT",
			"payload": map[string]string{
				"text": "Test Insight",
			},
		}
	case "AUDIENCE":
		payload = map[string]interface{}{
			"type": "AUDIENCE",
			"payload": map[string]interface{}{
				"gender":                      "Female",
				"birth_country":               "Canada",
				"age_group":                   "35-44",
				"daily_hours_on_social_media": 3.5,
				"purchases_last_month":        5,
			},
		}
	default:
		return 0, fmt.Errorf("unsupported asset type: %s", assetType)
	}

	assetJSON, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/assets/add", "application/json", bytes.NewBuffer(assetJSON))
	if err != nil {
		return 0, fmt.Errorf("failed to add asset: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logRequestResponse("POST", baseURL+"/assets/add", assetJSON, body, resp.StatusCode)

	if resp.StatusCode != http.StatusCreated {
		return 0, fmt.Errorf("expected status Created, got %v", resp.Status)
	}

	var createdAsset map[string]interface{}
	json.Unmarshal(body, &createdAsset)
	return int(createdAsset["id"].(float64)), nil
}

func deleteAsset(assetID int, assetType string) error {
	url := fmt.Sprintf("%s/assets/delete?id=%d&type=%s", baseURL, assetID, assetType)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send delete request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logRequestResponse("DELETE", url, nil, body, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status OK, got %v", resp.Status)
	}

	return nil
}

func getAsset(assetID int, assetType string) error {
	url := fmt.Sprintf("%s/assets/get?id=%d&type=%s", baseURL, assetID, assetType)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logRequestResponse("GET", url, nil, body, resp.StatusCode)

	if resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("asset not found")
	} else {
		return fmt.Errorf("unexpected status: %v", resp.Status)
	}
}

func modifyAssetTest() error {
	log.Println("Starting Modify Asset Test")

	assetTypes := []string{"CHART", "INSIGHT", "AUDIENCE"}

	for _, assetType := range assetTypes {
		log.Printf("Testing asset type: %s", assetType)

		// Add an asset
		assetID, err := addSingleAsset(assetType)
		if err != nil {
			return fmt.Errorf("failed to add %s asset: %v", assetType, err)
		}
		log.Printf("Successfully added %s asset with ID: %d", assetType, assetID)

		// Modify the asset
		if err := modifyAsset(assetType, assetID); err != nil {
			return fmt.Errorf("failed to modify %s asset: %v", assetType, err)
		}
		log.Printf("Successfully modified %s asset with ID: %d", assetType, assetID)

		// Verify the modification
		if err := verifyAsset(assetType, assetID, true); err != nil {
			return fmt.Errorf("failed to verify modified %s asset: %v", assetType, err)
		}
		log.Printf("Successfully verified modified %s asset with ID: %d", assetType, assetID)

		log.Printf("Modify %s Asset Test passed successfully!", assetType)
	}

	log.Println("All Modify Asset Tests passed successfully!")
	return nil
}

func modifyAsset(assetType string, assetID int) error {
	var payload map[string]interface{}
	switch assetType {
	case "CHART":
		payload = map[string]interface{}{
			"type": "CHART",
			"payload": map[string]interface{}{
				"id":           assetID,
				"title":        "Modified Test Chart",
				"x_axis_title": "Modified X Axis",
				"y_axis_title": "Modified Y Axis",
				"data_points":  []map[string]float64{{"x": 3, "y": 6}, {"x": 4, "y": 8}},
			},
		}
	case "INSIGHT":
		payload = map[string]interface{}{
			"type": "INSIGHT",
			"payload": map[string]interface{}{
				"id":   assetID,
				"text": "Modified Test Insight",
			},
		}
	case "AUDIENCE":
		payload = map[string]interface{}{
			"type": "AUDIENCE",
			"payload": map[string]interface{}{
				"id":                          assetID,
				"gender":                      "Male",
				"birth_country":               "UK",
				"age_group":                   "45-54",
				"daily_hours_on_social_media": 4.5,
				"purchases_last_month":        7,
			},
		}
	default:
		return fmt.Errorf("unsupported asset type: %s", assetType)
	}

	assetJSON, _ := json.Marshal(payload)
	req, err := http.NewRequest("PUT", baseURL+"/assets/modify", bytes.NewBuffer(assetJSON))
	if err != nil {
		return fmt.Errorf("failed to create modify request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send modify request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logRequestResponse("PUT", baseURL+"/assets/modify", assetJSON, body, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status OK, got %v", resp.Status)
	}

	return nil
}

func verifyAsset(assetType string, assetID int, isModified bool) error {

	url := fmt.Sprintf("%s/assets/get?id=%d&type=%s", baseURL, assetID, assetType)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to get asset: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	logRequestResponse("GET", url, nil, body, resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status OK, got %v", resp.Status)
	}

	var asset map[string]interface{}
	if err := json.Unmarshal(body, &asset); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	switch assetType {
	case "CHART":
		expectedTitle := "Test Chart for Modification"
		expectedXAxis := "X Axis"
		expectedYAxis := "Y Axis"
		if isModified {
			expectedTitle = "Modified Test Chart"
			expectedXAxis = "Modified X Axis"
			expectedYAxis = "Modified Y Axis"
		}
		if asset["title"] != expectedTitle ||
			asset["x_axis_title"] != expectedXAxis ||
			asset["y_axis_title"] != expectedYAxis {
			return fmt.Errorf("chart verification failed")
		}
	case "INSIGHT":
		expectedText := "Test Insight for Modification"
		if isModified {
			expectedText = "Modified Test Insight"
		}
		if asset["text"] != expectedText {
			return fmt.Errorf("insight verification failed")
		}
	case "AUDIENCE":
		expectedGender := "Female"
		expectedCountry := "Canada"
		expectedAgeGroup := "35-44"
		if isModified {
			expectedGender = "Male"
			expectedCountry = "UK"
			expectedAgeGroup = "45-54"
		}
		if asset["gender"] != expectedGender ||
			asset["birth_country"] != expectedCountry ||
			asset["age_group"] != expectedAgeGroup {
			return fmt.Errorf("audience verification failed")
		}
	default:
		return fmt.Errorf("unsupported asset type: %s", assetType)
	}

	return nil
}
