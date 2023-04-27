package webserver

import (
	"bytes"
	"challenge/cache"
	"challenge/models"
	"challenge/storage"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestSigninExistingUser tests all the possible scenarios for user signin //
func TestSigninExistingUser(t *testing.T) {
	// Run the test scenarios
	for scenario, fn := range map[string]func(t *testing.T){
		"invalid request body":           testSigninInvalidBody,
		"error retrieving user password": testSigninFetchPassword,
		"user not found":                 testSigninUserNotFound,
		"incorrect user password":        testSigninIncorrectPassword,
		"valid signin":                   testSigninValid,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// SigninInvalidBody (test scenario)
func testSigninInvalidBody(t *testing.T) {
	// Signup a user
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id":       "testuser",
		"password": "testpassword",
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)
	requestBody2 := []byte("test")
	req, err = http.NewRequest("POST", "/signin", bytes.NewBuffer(requestBody2))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr = httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler = http.HandlerFunc(SigninExistingUser)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	storage.Db.Release()
}

// SigninFetchPassword (test scenario)
func testSigninFetchPassword(t *testing.T) {
	// Signup a user
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id":       "testuser",
		"password": "testpassword",
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)
	// Release the Db to cause the retrieval error
	storage.Db.Release()
	// Create a new request with a JSON-encoded user in the body
	requestBody2 := map[string]string{
		"id":       "testuser",
		"password": "testPassword",
	}
	requestBodyJson, _ = json.Marshal(requestBody2)
	req, err = http.NewRequest("POST", "/signin", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr = httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler2 := http.HandlerFunc(SigninExistingUser)
	handler2.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

// SigninUserNotFound (test scenario)
func testSigninUserNotFound(t *testing.T) {
	// Signup a user
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id":       "testuser",
		"password": "testpassword",
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)
	requestBody2 := map[string]string{
		"id":       "testuser2",
		"password": "testpassword",
	}
	requestBodyJson, _ = json.Marshal(requestBody2)
	req, err = http.NewRequest("POST", "/signin", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr = httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler2 := http.HandlerFunc(SigninExistingUser)
	handler2.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	// Check the response body is correct
	expected := "User with that id was not found"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	storage.Db.Release()
}

// SigninIncorrectPassword (test scenario)
func testSigninIncorrectPassword(t *testing.T) {
	// Signup a user
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id":       "testuser",
		"password": "testpassword",
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)
	requestBody2 := map[string]string{
		"id":       "testuser",
		"password": "testpassword2",
	}
	requestBodyJson, _ = json.Marshal(requestBody2)
	req, err = http.NewRequest("POST", "/signin", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr = httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler2 := http.HandlerFunc(SigninExistingUser)
	handler2.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// Check the response body is correct
	expected := "Incorrect password for this user"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	storage.Db.Release()
}

// SigninValid (test scenario)
func testSigninValid(t *testing.T) {
	// Signup a user
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id":       "testuser",
		"password": "testpassword",
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)
	requestBody2 := map[string]string{
		"id":       "testuser",
		"password": "testpassword",
	}
	requestBodyJson, _ = json.Marshal(requestBody2)
	req, err = http.NewRequest("POST", "/signin", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr = httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler2 := http.HandlerFunc(SigninExistingUser)
	handler2.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	storage.Db.Release()
}

// TestSignoutExistingUser tests all the possible scenarios for user signout //
func TestSignoutExistingUser(t *testing.T) {
	// Initialize a user
	MockValidUser("testuser", "testpassword")
	// Run the test scenarios
	for scenario, fn := range map[string]func(t *testing.T){
		"invalid signout": testSignoutInvalid,
		"valid signout":   testSignoutValid,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
	storage.Db.Release()
}

// SignoutInvalid (test scenario)
func testSignoutInvalid(t *testing.T) {
	req, err := http.NewRequest("GET", "/signout", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(SignoutExistingUser)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusNotAcceptable {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotAcceptable)
	}
}

// SignoutValid (test scenario)
func testSignoutValid(t *testing.T) {
	req, err := http.NewRequest("GET", "/signout", nil)
	if err != nil {
		t.Fatal(err)
	}
	userExists, cookie, err := storage.Db.FetchUserToken("testuser")
	if err != nil {
		t.Fatal(err)
	}
	if !userExists {
		t.Fatalf("user was not found in storage")
	}

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	req.Header.Set("Cookie", cookie)
	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(SignoutExistingUser)
	handler.ServeHTTP(rr, req)
	// Check the response status code is correct
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

// TestGetBuildVersion tests the retrieval of build version information //
func TestGetBuildVersion(t *testing.T) {
	// Run the test scenarios
	for scenario, fn := range map[string]func(t *testing.T){
		"vcs info available": testGetBuildVersionAvailable,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testGetBuildVersionAvailable (test scenario)
func testGetBuildVersionAvailable(t *testing.T) {
	// Create a new request with
	req, err := http.NewRequest("GET", "/version", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(GetBuildVersion)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("failed to retrieve build version")
	}
}

// TestCreateUserHandler tests all the possible scenarios for manual user //
// creation //
func TestCreateUserHandler(t *testing.T) {
	// Run the test scenarios
	for scenario, fn := range map[string]func(t *testing.T){
		"invalid request body":    testCreateUserInvalidBody,
		"missing required fields": testCreateUserMissingCredentials,
		"user check error":        testCreateUserCheckError,
		"user already exists":     testCreateUserExists,
		"password hashing error":  testCreateUserHashingError,
		"error creating new user": testCreateUserError,
		"new user created":        testCreateUserSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testCreateUserInvalidBody (test scenario)
func testCreateUserInvalidBody(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := []byte("test")
	req, err := http.NewRequest("POST", "/users/testuser/create", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	storage.Db.Release()
}

// func testCreateUserMissingCredentials (test scenario)
func testCreateUserMissingCredentials(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id": "testuser",
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/users/testuser/create", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Check the response body is correct
	expected := "Missing required fields in request body"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	storage.Db.Release()
}

// testCreateUserCheckError (test scenario)
func testCreateUserCheckError(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id":       "testuser",
		"password": "testPassword",
	}
	if err := storage.Db.Release(); err != nil {
		t.Fatalf("failed to release storage resource")
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/users/testuser/create", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

// testCreateUserExists (test scenario)
func testCreateUserExists(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id":       "testuser",
		"password": "testPassword",
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/users/testuser/create", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)
	// Check the response status code is correct
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	req, err = http.NewRequest("POST", "/users/testuser/create", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr = httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler = http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
	if err := storage.Db.Release(); err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testCreateUserHashingError (test scenario)
func testCreateUserHashingError(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id":       "testuser",
		"password": "testPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPassword",
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/users/testuser/create", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)
	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	if err := storage.Db.Release(); err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testCreateUserError (test scenario)
func testCreateUserError(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id":       "testuser",
		"password": "testPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPasswordtestPassword",
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/users/testuser/create", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)
	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	if err := storage.Db.Release(); err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testCreateUserSuccess (test scenario)
func testCreateUserSuccess(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id":       "testuser",
		"password": "testPassword",
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/users/testuser/create", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)
	// Check the response status code is correct
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	if err := storage.Db.Release(); err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestDeleteUserHandler tests all the possible scenarios for manual user //
// deletion //
func TestDeleteUserHandler(t *testing.T) {
	// Run the test scenarios
	for scenario, fn := range map[string]func(t *testing.T){
		"user check error":     testDeleteUserCheckError,
		"user does not exists": testDeleteUserNotExist,
		"user deleted":         testDeleteUserSuccess,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testDeleteUserCheckError (test scenario)
func testDeleteUserCheckError(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	if err := storage.Db.Release(); err != nil {
		t.Fatalf("failed to release storage resource")
	}
	req, err := http.NewRequest("DELETE", "/users/testuser/delete", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the DeleteUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(DeleteUserHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

// testDeleteUserNotExist (test scenario)
func testDeleteUserNotExist(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	req, err := http.NewRequest("POST", "/users/testuser/delete", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the DeleteUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(DeleteUserHandler)
	handler.ServeHTTP(rr, req)
	// Check the response status code is correct
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}

	if err := storage.Db.Release(); err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testDeleteUserSuccess (test scenario)
func testDeleteUserSuccess(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	newUser := models.User{
		ID:       "testuser",
		Password: "testpassword",
	}
	res, err := storage.Db.CreateNewUser(newUser)
	if err != nil {
		t.Fatalf("failed to create example user")
	}
	if !res {
		t.Fatalf("failed to add example user")
	}
	req, err := http.NewRequest("DELETE", "/users/testuser/delete", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the DeleteUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(DeleteUserHandler)
	handler.ServeHTTP(rr, req)
	// Check the response status code is correct
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	if err := storage.Db.Release(); err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestGetFavouritesHandler tests all the possible scenarios for user //
// favourites retrieval //
func TestGetFavouritesHandler(t *testing.T) {
	// Run the test scenarios
	for scenario, fn := range map[string]func(t *testing.T){
		"user check error":     testFavouritesUserCheckError,
		"missing cookie":       testFavouritesCookieMissing,
		"cache failed load":    testFavouritesCacheLoadError,
		"cache succesful load": testFavouritesCacheLoad,
		"storage load error":   testFavouritesStorageLoadError,
		"cache failed store":   testFavouritesCacheStoreError,
		"storage load":         testFavouritesStorageLoad,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// FavouritesUserCheckError (test scenario)
func testFavouritesUserCheckError(t *testing.T) {
	MockValidUser("testuser", "testpassword")
	storage.Db.Release()
	// Create a new request with
	req, err := http.NewRequest("GET", "/users/testuser/favourites", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(GetFavouritesHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

// FavouritesCookieMissing (test scenario)
func testFavouritesCookieMissing(t *testing.T) {
	MockValidUser("testuser", "testpassword")
	// Create a new request with
	req, err := http.NewRequest("GET", "/users/testuser/favourites", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	userExists, tok, err := storage.Db.FetchUserToken("testuser")
	if err != nil {
		t.Fatal(err)
	}
	if !userExists {
		t.Fatalf("user does not exist in storage")
	}
	rr.Header().Set("Cookie", tok)
	rr.Header().Del("cookie")
	// Call the GetFavouritesHandler function with the new request and response recorder
	handler := http.HandlerFunc(GetFavouritesHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
	storage.Db.Release()
}

// FavouritesCacheLoadError (test scenario)
func testFavouritesCacheLoadError(t *testing.T) {
	MockValidUser("testuser", "testpassword")
	cache.Cache.Release()
	_, tok, _ := storage.Db.FetchUserToken("testuser")

	// Create a new request with
	req, err := http.NewRequest("GET", "/users/testuser/favourites", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	req.Header.Set("Cookie", tok)
	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(GetFavouritesHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	storage.Db.Release()
}

// FavouritesCacheLoad (test scenario)
func testFavouritesCacheLoad(t *testing.T) {
	MockValidUser("testuser", "testpassword")
	userExists, cookie, err := storage.Db.FetchUserToken("testuser")
	if err != nil {
		t.Fatal(err)
	}
	if !userExists {
		t.Fatalf("user does not exist in storage")
	}
	parts := strings.Split(cookie, "=")
	tok := parts[1]
	err = cache.Cache.Store(tok, []byte("exampleResponse"), false)
	if err != nil {
		t.Fatal(err)
	}
	// Create a new request with
	req, err := http.NewRequest("GET", "/users/testuser/favourites", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	req.Header.Set("Cookie", cookie)
	req.Header.Set("app-user", "testuser")
	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(GetFavouritesHandler)
	handler.ServeHTTP(rr, req)
	// Check the response status code is correct
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	expected := "exampleResponse"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	cache.Cache.Release()
	storage.Db.Release()
}

// FavouritesStorageLoadError (test scenario)
func testFavouritesStorageLoadError(t *testing.T) {
	MockValidUser("testuser", "testpassword")
	storage.Db.Release()
	// Create a new request with
	req, err := http.NewRequest("GET", "/users/testuser/favourites", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(GetFavouritesHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	cache.Cache.Release()
}

// FavouritesCacheStoreError (test scenario)
func testFavouritesCacheStoreError(t *testing.T) {
	MockValidUser("testuser", "testpassword")
	cache.Cache.Release()
	userExists, cookie, err := storage.Db.FetchUserToken("testuser")
	if err != nil {
		t.Fatal(err)
	}
	if !userExists {
		t.Fatalf("user does not exist in storage")
	}
	// Create a new request with
	req, err := http.NewRequest("GET", "/users/testuser/favourites", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	req.Header.Set("Cookie", cookie)
	req.Header.Set("app-user", "testuser")
	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(GetFavouritesHandler)
	handler.ServeHTTP(rr, req)
	storage.Db.Release()
}

// FavouritesStorageLoad (test scenario)
func testFavouritesStorageLoad(t *testing.T) {
	MockValidUser("testuser", "testpassword")
	storage.Db.AddAssetToUser("testuser", models.Asset{
		ID:          "asset1",
		Type:        "chart",
		Description: "description1",
		AssetData:   nil,
		Added:       time.Now(),
		Modified:    time.Now(),
	})

	userExists, tok, err := storage.Db.FetchUserToken("testuser")
	if err != nil {
		t.Fatal(err)
	}
	if !userExists {
		t.Fatalf("user does not exist in storage")
	}
	// Create a new request with
	req, err := http.NewRequest("GET", "/users/testuser/favourites", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	req.Header.Set("Cookie", tok)
	req.Header.Set("app-user", "testuser")

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(GetFavouritesHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	storage.Db.Release()
	cache.Cache.Release()
}

// TestAddFavouriteHandler tests all the possible scenarios for adding //
// new favourite assets //
func TestAddFavouriteHandler(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"invalid request body":        testAddFavouriteInvalidBody,
		"asset existence check error": testAddFavouriteAssetCheckError,
		"asset already exists":        testAddFavouriteAssetExists,
		"adding asset error":          testAddFavouriteAssetAddError,
		"adding asset success":        testAddFavouriteAssetAdd,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testAddFavouriteInvalidBody  (test scenario)
func testAddFavouriteInvalidBody(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := []byte("test")
	req, err := http.NewRequest("POST", "/users/testuser/favourites", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the AddFavouriteHandler function with the new request and response recorder
	handler := http.HandlerFunc(AddFavouriteHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testAddFavouriteAssetCheckError  (test scenario)
func testAddFavouriteAssetCheckError(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	newAsset := models.Asset{
		ID:          "asset1",
		Type:        "chart",
		Description: "example description",
		AssetData:   nil,
		Added:       time.Now(),
		Modified:    time.Now(),
	}
	jsonDoc, err := json.Marshal(newAsset)
	if err != nil {
		t.Fatalf("failed to serialize asset")
	}
	req, err := http.NewRequest("POST", "/users/testuser/favourites", bytes.NewBuffer(jsonDoc))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the AddFavouriteHandler function with the new request and response recorder
	handler := http.HandlerFunc(AddFavouriteHandler)
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

// testAddFavouriteAssetExists  (test scenario)
func testAddFavouriteAssetExists(t *testing.T) {
	// Create example user
	MockValidUser("testuser", "testpassword")
	// Add example asset to user
	newAsset := models.Asset{
		ID:          "asset1",
		Type:        "chart",
		Description: "example description",
		AssetData:   nil,
		Added:       time.Now(),
		Modified:    time.Now(),
	}
	_, tok, _ := storage.Db.FetchUserToken("testuser")
	// Create a new request with a JSON-encoded user in the body
	jsonDoc, err := json.Marshal(newAsset)
	if err != nil {
		t.Fatalf("failed to serialize asset")
	}
	req, err := http.NewRequest("POST", "/users/testuser/favourites", bytes.NewBuffer(jsonDoc))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	req.Header.Add("Cookie", tok)
	// Call the AddFavouriteHandler function with the new request and response recorder
	handler := http.HandlerFunc(AddFavouriteHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	// Second request
	req3, err := http.NewRequest("POST", "/users/testuser/favourites", bytes.NewBuffer(jsonDoc))
	if err != nil {
		t.Fatal(err)
	}
	rr3 := httptest.NewRecorder()
	req3.Header.Set("Content-Type", "application/json")
	req3.Header.Add("Cookie", tok)
	handler3 := http.HandlerFunc(AddFavouriteHandler)
	handler3.ServeHTTP(rr3, req3)
	// Check the response status code is correct
	if status := rr3.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusConflict)
	}
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testAddFavouriteAssetAddError  (test scenario)
func testAddFavouriteAssetAddError(t *testing.T) {
	// Create example user
	MockValidUser("testuser", "testpassword")
	// Add example asset to user
	newAsset := models.Asset{
		ID:          "asset1",
		Type:        "chart",
		Description: "example description",
		AssetData:   nil,
		Added:       time.Now(),
		Modified:    time.Now(),
	}
	_, tok, _ := storage.Db.FetchUserToken("testuser")
	// Create a new request with a JSON-encoded user in the body
	jsonDoc, err := json.Marshal(newAsset)
	if err != nil {
		t.Fatalf("failed to serialize asset")
	}
	req, err := http.NewRequest("POST", "/users/testuser/favourites", bytes.NewBuffer(jsonDoc))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	req.Header.Add("Cookie", tok)
	// Release storage to cause error
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	// Call the AddFavouriteHandler function with the new request and response recorder
	handler := http.HandlerFunc(AddFavouriteHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

// testAddFavouriteAssetAdd  (test scenario)
func testAddFavouriteAssetAdd(t *testing.T) {
	// Create example user
	MockValidUser("testuser", "testpassword")
	// Add example asset to user
	newAsset := models.Asset{
		ID:          "asset1",
		Type:        "chart",
		Description: "example description",
		AssetData:   nil,
		Added:       time.Now(),
		Modified:    time.Now(),
	}
	_, tok, _ := storage.Db.FetchUserToken("testuser")
	// Create a new request with a JSON-encoded user in the body
	jsonDoc, err := json.Marshal(newAsset)
	if err != nil {
		t.Fatalf("failed to serialize asset")
	}
	req, err := http.NewRequest("POST", "/users/testuser/favourites", bytes.NewBuffer(jsonDoc))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	req.Header.Add("Cookie", tok)
	// Call the AddFavouriteHandler function with the new request and response recorder
	handler := http.HandlerFunc(AddFavouriteHandler)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// TestEditAssetDescriptionHandler tests all the possible scenarios for updating //
// a favourited asset's description //
func TestEditAssetDescriptionHandler(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"invalid request body":             testEditDescriptionInvalidBody,
		"asset existence check error":      testEditDescriptionAssetExistenceCheckError,
		"asset not found":                  testEditDescriptionAssetNotFound,
		"update asset description error":   testEditDescriptionUpdateError,
		"update asset description success": testEditDescriptionUpdate,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testEditDescriptionInvalidBody  (test scenario)
func testEditDescriptionInvalidBody(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	requestBody := []byte("test")
	req, err := http.NewRequest("POST", "/users/testuser/favourites", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the EditAssetDescriptionHandler function with the new request and response recorder
	handler := http.HandlerFunc(EditAssetDescriptionHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testEditDescriptionAssetExistenceCheckError  (test scenario)
func testEditDescriptionAssetExistenceCheckError(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create a new request with a JSON-encoded user in the body
	newAsset := map[string]string{
		"description": "new description",
	}
	jsonDoc, err := json.Marshal(newAsset)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users/testuser/favourites/asset1", bytes.NewBuffer(jsonDoc))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	// Release resource to trigger error on request
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	// Call the EditAssetDescriptionHandler function with the new request and response recorder
	handler := http.HandlerFunc(EditAssetDescriptionHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

// testEditDescriptionAssetNotFound  (test scenario)
func testEditDescriptionAssetNotFound(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create example user
	MockValidUser("testuser", "testpass")
	// Create a new request with a JSON-encoded user in the body
	newAsset := map[string]string{
		"description": "new description",
	}
	jsonDoc, err := json.Marshal(newAsset)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users/testuser/favourites/asset1", bytes.NewBuffer(jsonDoc))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the EditAssetDescriptionHandler function with the new request and response recorder
	handler := http.HandlerFunc(EditAssetDescriptionHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	// Release resource
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	err = cache.Cache.Release()
	if err != nil {
		t.Fatalf("failed to release cache resource")
	}
}

// testEditDescriptionUpdateError  (test scenario)
func testEditDescriptionUpdateError(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create example user
	MockValidUser("testuser", "testpass")
	// Create example asset
	newAsset := models.Asset{}
	storage.Db.AddAssetToUser("testuser", newAsset)
	// Create a new request with a JSON-encoded user in the body
	newAssetEx := map[string]string{
		"description": "new description",
	}
	jsonDoc, err := json.Marshal(newAssetEx)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users/testuser/favourites/asset1", bytes.NewBuffer(jsonDoc))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Release resource to trigger error on request
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}

	// Call the EditAssetDescriptionHandler function with the new request and response recorder
	handler := http.HandlerFunc(EditAssetDescriptionHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
	err = cache.Cache.Release()
	if err != nil {
		t.Fatalf("failed to release cache resource")
	}
}

// testEditDescriptionUpdate  (test scenario)
func testEditDescriptionUpdate(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create example user
	MockValidUser("testuser", "testpass")
	// Create a new request with a JSON-encoded user in the body
	newAssetEx := map[string]string{
		"description": "new description",
	}
	jsonDoc, err := json.Marshal(newAssetEx)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", "/users/testuser/favourites/assetExample", bytes.NewBuffer(jsonDoc))
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Release resource to trigger error on request
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}

	// Call the EditAssetDescriptionHandler function with the new request and response recorder
	handler := http.HandlerFunc(EditAssetDescriptionHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			rr.Body.String(), "Asset removed from favourites")
	}
	err = cache.Cache.Release()
	if err != nil {
		t.Fatalf("failed to release cache resource")
	}
}

// TestRemoveFavouriteHandler tests all the possible scenarios for removing //
// a favourited asset from a user //
func TestRemoveFavouriteHandler(t *testing.T) {
	for scenario, fn := range map[string]func(t *testing.T){
		"asset existence check error": testRemoveFavouriteExistenceCheckError,
		"asset not found":             testRemoveFavouriteNotFound,
		"asset removal error":         testRemoveFavouriteRemovalError,
		"asset removal success":       testRemoveFavouriteRemoval,
	} {
		t.Run(scenario, func(t *testing.T) {
			fn(t)
		})
	}
}

// testRemoveFavouriteExistenceCheckError  (test scenario)
func testRemoveFavouriteExistenceCheckError(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	req, err := http.NewRequest("DELETE", "/users/testuser/favourites/asset1", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()
	// Release resource to trigger error on request
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
	// Call the EditAssetDescriptionHandler function with the new request and response recorder
	handler := http.HandlerFunc(RemoveFavouriteHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

// testRemoveFavouriteNotFound  (test scenario)
func testRemoveFavouriteNotFound(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create example user
	MockValidUser("testuser", "testpass")
	req, err := http.NewRequest("DELETE", "/users/testuser/favourites/asset1", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the RemoveFavouriteHandler function with the new request and response recorder
	handler := http.HandlerFunc(RemoveFavouriteHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
	// Release resource
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}

// testRemoveFavouriteRemovalError  (test scenario)
func testRemoveFavouriteRemovalError(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create example user
	MockValidUser("testuser", "testpass")
	req, err := http.NewRequest("DELETE", "/users/testuser/favourites/assetExample", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Release resource
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}

	// Call the RemoveFavouriteHandler function with the new request and response recorder
	handler := http.HandlerFunc(RemoveFavouriteHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}
}

// testRemoveFavouriteRemoval  (test scenario)
func testRemoveFavouriteRemoval(t *testing.T) {
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	// Create example user
	MockValidUser("testuser", "testpass")

	_, tok, _ := storage.Db.FetchUserToken("testuser")
	req, err := http.NewRequest("DELETE", "/users/testuser/favourites/assetExample", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", tok)
	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the RemoveFavouriteHandler function with the new request and response recorder
	handler := http.HandlerFunc(RemoveFavouriteHandler)
	handler.ServeHTTP(rr, req)

	// Check the response status code is correct
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	res, err := storage.Db.CheckUserAssetExistence("testuser", "assetExample")
	if err != nil {
		t.Fatalf("failed to check asset existence")
	}
	if res {
		t.Fatalf("asset still exists")
	}

	// Release resource
	err = storage.Db.Release()
	if err != nil {
		t.Fatalf("failed to release storage resource")
	}
}
