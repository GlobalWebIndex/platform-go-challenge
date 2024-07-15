package handlers_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ceciivanov/platform-go-challenge/internal/handlers"
	"github.com/ceciivanov/platform-go-challenge/internal/repository"
	"github.com/ceciivanov/platform-go-challenge/internal/service"
	"github.com/gorilla/mux"
)

const (
	numUsers     = 1000
	numFavorites = 100
)

// Each benchmark test runs on its own instance of the repository,
// also the timer is reset after initializing data to remove the setup from the overall benchmark time.

func BenchmarkGetUserFavorites(b *testing.B) {
	// Initialize user repository with sample data
	repo := repository.NewInMemoryUserRepository()
	repo.GenerateSampleUsers(numUsers, numFavorites)

	// Initialize user service and handler
	userService := service.NewUserService(repo)
	userHandler := handlers.NewUserHandler(userService)

	// Create a new router and assign the handler
	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		// Select a random user ID within the range of numUsers
		userID := rand.Intn(numUsers) + 1 // Adding 1 because user IDs start from 1

		// Create a new HTTP request with the payload
		req, err := http.NewRequest("GET", fmt.Sprintf("/users/%d/favorites", userID), nil)
		if err != nil {
			b.Fatal(err)
		}

		rr := httptest.NewRecorder()
		start := time.Now()
		r.ServeHTTP(rr, req)
		elapsed := time.Since(start)

		// Check the status code
		if status := rr.Code; status != http.StatusOK {
			b.Errorf("handler returned wrong status code, got: %v expected: %v", status, http.StatusOK)
		}

		// Log the time taken for the request
		b.Logf("Request took %v", elapsed)
	}
}

func BenchmarkAddUserFavorite(b *testing.B) {
	// Initialize user repository with sample data
	repo := repository.NewInMemoryUserRepository()
	repo.GenerateSampleUsers(numUsers, numFavorites)

	// Initialize user service and handler
	userService := service.NewUserService(repo)
	userHandler := handlers.NewUserHandler(userService)

	// Create a new router and assign the handler
	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Select a random user ID within the range of numUsers
		userID := rand.Intn(numUsers) + 1 // Adding 1 because user IDs start from 1

		assetID := 1001 + i

		// Prepare JSON payload for the asset you want to add as favorite
		payload := []byte(fmt.Sprintf(`{
			"id": %d,
			"type": "Insight",
			"description": "Sample Insight for testing",
			"text": "Testing Insight"
		}`, assetID))

		// Create a new HTTP request with the payload
		req, err := http.NewRequest("POST", fmt.Sprintf("/users/%d/favorites", userID), bytes.NewBuffer(payload))
		if err != nil {
			b.Fatal(err)
		}

		// Set the request Content-Type header
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		start := time.Now()
		r.ServeHTTP(rr, req)
		elapsed := time.Since(start)

		// Check the status code
		if status := rr.Code; status != http.StatusCreated {
			b.Errorf("handler returned wrong status code, got: %v expected: %v", status, http.StatusCreated)
		}

		// Log the time taken for the request
		b.Logf("Request took %v", elapsed)
	}
}

func BenchmarkDeleteUserFavorite(b *testing.B) {

	// Initialize user repository with sample data
	repo := repository.NewInMemoryUserRepository()
	repo.GenerateSampleUsers(numUsers, numFavorites)

	// Initialize user service and handler
	userService := service.NewUserService(repo)
	userHandler := handlers.NewUserHandler(userService)

	// Create a new router and assign the handler
	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)

	for i := 0; i < b.N; i++ {
		// Select a random user ID within the range of numUsers
		userID := rand.Intn(numUsers) + 1 // Adding 1 because user IDs start from 1

		assetID := i + 1 // asset IDs starts from 1 to numFavorites
		if assetID > numFavorites {
			continue
		}

		// Create a new HTTP request
		req, err := http.NewRequest("DELETE", fmt.Sprintf("/users/%d/favorites/%d", userID, assetID), nil)
		if err != nil {
			b.Fatal(err)
		}

		rr := httptest.NewRecorder()
		start := time.Now()
		r.ServeHTTP(rr, req)
		elapsed := time.Since(start)

		// Check the status code
		if status := rr.Code; status != http.StatusNoContent {
			b.Errorf("handler returned wrong status code, got: %v expected: %v", status, http.StatusNoContent)
		}

		// Log the time taken for the request
		b.Logf("Request took %v", elapsed)
	}
}

func BenchmarkEditUserFavorite(b *testing.B) {
	// Initialize a new repository for each iteration
	repo := repository.NewInMemoryUserRepository()
	repo.GenerateSampleUsers(numUsers, numFavorites)

	// Create a new router and assign the handler
	userService := service.NewUserService(repo)
	userHandler := handlers.NewUserHandler(userService)

	r := mux.NewRouter()
	userHandler.RegisterRoutes(r)

	// Reset benchmark timer before starting
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Select a random user ID within the range of numUsers
		userID := rand.Intn(numUsers) + 1 // Adding 1 because user IDs start from 1

		// Retrieve user's favorites to find a matching asset type
		user, ok := repo.Users[userID]
		if !ok {
			b.Fatal("User not found")
		}

		// Find only assets of type 'Insight' to update
		var assetID int
		for _, fav := range user.Favourites {
			// Example: Update only 'Insight' type assets
			if fav.GetType() == "Insight" {
				assetID = fav.GetID()
				break
			}
		}

		// Prepare JSON payload for the asset update
		payload := []byte(fmt.Sprintf(`{
			"id": %d,
			"type": "Insight",
			"description": "Updated Insight",
			"text": "Updated Insight Text"
		}`, assetID))

		// Create a new HTTP request with the payload
		req, err := http.NewRequest("PUT", fmt.Sprintf("/users/%d/favorites/%d", userID, assetID), bytes.NewBuffer(payload))
		if err != nil {
			b.Fatal(err)
		}

		// Set the request Content-Type header
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		start := time.Now()
		r.ServeHTTP(rr, req)
		elapsed := time.Since(start)

		// Check the status code
		if status := rr.Code; status != http.StatusOK {
			b.Errorf("handler returned wrong status code, got: %v expected: %v", status, http.StatusOK)
		}

		// Log the time taken for the request
		b.Logf("Request took %v", elapsed)
	}
}
