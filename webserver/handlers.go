package webserver

import (
	"bytes"
	"challenge/cache"
	"challenge/middleware"
	"challenge/models"
	"challenge/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// SigninExistingUser logs in existing user
func SigninExistingUser(w http.ResponseWriter, r *http.Request) {
	var existingUser models.User
	err := json.NewDecoder(r.Body).Decode(&existingUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body")
		return
	}
	// Get the existing entry present in the database for the given username
	found, encPass, err := storage.Db.FetchUserPass(existingUser.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// If password is empty, then user with that ID was not found
	if !found { // encPass == ""
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User with that id was not found")
		return
	}
	// Compare the stored hashed password, with the hashed version of the password that was received
	if err = bcrypt.CompareHashAndPassword([]byte(encPass), []byte(existingUser.Password)); err != nil {
		// If the two passwords don't match, return a 401 status
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Incorrect password for this user")
		return
	}
	// Credentials have been verified, generate token
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims (username and expiry time)
	claims := &models.Claims{
		ID: existingUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newCookie := "token=" + tokenString
	found, err = storage.Db.UpdateUserToken(existingUser.ID, newCookie)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User was not found")
		return
	}
	// Set the client cookie for "token" and an expiry time (same as token)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	w.WriteHeader(http.StatusOK)
}

// SignoutExistingUser logs out a user that has logged in
func SignoutExistingUser(w http.ResponseWriter, r *http.Request) {
	c := r.Header.Get("Cookie")
	if c == "" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	// immediately clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
}

// GetBuildVersion retrives the build/git information of the binary
func GetBuildVersion(w http.ResponseWriter, r *http.Request) {
	payload := make(map[string]string)
	payload["version"] = Version
	payload["buildTime"] = BuildTime
	payload["commitHash"] = CommitHash
	payload["buildUser"] = BuildUser
	w.Header().Set("Content-Type", "application/json")
	jsonDoc, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonDoc)
}

// CreateUserHandler manually creates a new user (requires Privileges)
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// ID and password are required fields
	if newUser.ID == "" || newUser.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing required fields in request body")
		return
	}
	// Check if user with that username already exists
	exists, err := storage.Db.CheckUserExistence(newUser.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// User already exists
	if exists {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "User with that username already exists")
		return
	}
	// Encrypt the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 8)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)
	// Add new user to storage
	result, err := storage.Db.CreateNewUser(newUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !result {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "User with that username already exists")
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "New user registered")
}

// DeleteUserHandler manually deletes a user (requires Privileges)
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	// Check if asset exists for requested user
	if userID == "" {
		parts := strings.Split(r.URL.String(), "/")
		userID = parts[2]
	}
	found, err := storage.Db.RemoveUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error removing user")
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User with that ID was not found")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User removed from favourites")
}

// GetFavouritesHandler returns a list of all the user's favourites assets
func GetFavouritesHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	if userID == "" {
		parts := strings.Split(r.URL.String(), "/")
		userID = parts[2]
	}
	// Check if user with `user_id` exists
	userExists, err := storage.Db.CheckUserExistence(userID) // exists
	if err != nil {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// User with `user_id` does not exist
	if !userExists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User with that id does not exist in storage")
		return
	}

	// User with `user_id` exists, so load the favourites
	// Check if Cookie is set
	_, exists := r.Header["Cookie"]
	if !exists {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Check if response is cached
	cookie := r.Header["Cookie"]
	token := strings.Split(cookie[0], "=")[1]
	set, cachedResp, err := cache.Cache.Load(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if set {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(cachedResp)
		return
	}

	// Favourites are not cached so retrieve it from storage
	found, res, err := storage.Db.RetrieveUserAssets(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User does not exist")
		return
	}
	// Serialize and cache the response
	jsonDoc, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = cache.Cache.Store(token, jsonDoc, true)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonDoc)
}

// AddFavouriteHandler allows the user to add an asset to their favourites
func AddFavouriteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	if userID == "" {
		parts := strings.Split(r.URL.String(), "/")
		userID = parts[2]
	}
	var asset models.Asset
	err := json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body")
		return
	}
	// Check if asset exists for requested user
	assetExists, err := storage.Db.CheckUserAssetExistence(userID, asset.ID)
	if err != nil {
		// BASED ON ERROR reply USER NOT FOUND or INTERNAL SERVER ERROR
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error checking user's asset existence")
		return
	}
	// Asset with that `asset_id` already exists for this user
	if assetExists {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Asset with that id already exists for this user")
		return
	}
	// Asset with that `asset_id` does not exist for this user, add it
	res, err := storage.Db.AddAssetToUser(userID, asset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error adding asset to user")
		return
	}
	if !res {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "user with that ID not found")
		return
	}

	// Invalidate existing cache entry since the favourites list changed
	cookie := r.Header["Cookie"]
	token := strings.Split(cookie[0], "=")[1]
	err = cache.Cache.Evict(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error evicting cache entry")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Asset added to user's favourites")
}

// EditAssetDescription allows the user to edit the description of an asset
func EditAssetDescriptionHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	if userID == "" {
		parts := strings.Split(r.URL.String(), "/")
		userID = parts[2]
	}
	assetID := params["asset_id"]
	if assetID == "" {
		parts := strings.Split(r.URL.String(), "/")
		assetID = parts[3]
	}
	var description struct {
		Description string `json:"description"`
	}
	err := json.NewDecoder(r.Body).Decode(&description)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error decoding request body")
		return
	}
	// Check if asset exists for requested user
	assetExists, err := storage.Db.CheckUserAssetExistence(userID, assetID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error checking user's asset existence")
		return
	}
	// Asset with that `asset_id`  does not exist for this user
	if !assetExists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Asset with that id not found for this user")
		return
	}
	res, err := storage.Db.UpdateUserAssetDescription(userID, assetID, description.Description)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error updating asset description")
		return
	}
	if !res {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Asset with that ID was not found")
		return
	}
	// Invalidate existing cache entry since the favourites list changed
	cookie := r.Header["Cookie"]
	token := strings.Split(cookie[0], "=")[1]
	err = cache.Cache.Evict(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error evicting cache entry")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Asset description updated")
}

// RemoveFavouriteHandler allows the user to remove an asset from their favourites
func RemoveFavouriteHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["user_id"]
	if userID == "" {
		parts := strings.Split(r.URL.String(), "/")
		userID = parts[2]
	}
	assetID := params["asset_id"]
	if assetID == "" {
		parts := strings.Split(r.URL.String(), "/")
		assetID = parts[4]
	}
	// Check if asset exists for requested user
	assetExists, err := storage.Db.CheckUserAssetExistence(userID, assetID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error checking user's asset existence")
		return
	}
	// Asset with that `asset_id` already exists for this user
	if !assetExists {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Asset with that id not found for this user")
		return
	}
	found, err := storage.Db.RemoveUserAsset(userID, assetID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error removing asset from user")
		return
	}
	if !found {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Asset with that ID was not found for this user")
		return
	}
	// Invalidate existing cache entry since the favourites list changed
	cookie := r.Header["Cookie"]
	token := strings.Split(cookie[0], "=")[1]
	err = cache.Cache.Evict(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error evicting cache entry")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Asset removed from favourites")
}

// MISCELLANEOUS FUNCTIONS //

// MockValidUser performs the signup and signin operation for an example valid
// user
func MockValidUser(id string, pass string) {
	// Signup a user
	// Create a connection to storage
	storage.ConnectToStorage("memory")
	cache.ConnectToCache("memory")

	// Create a new request with a JSON-encoded user in the body
	requestBody := map[string]string{
		"id":       id,
		"password": pass,
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest("POST", "/signup", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr := httptest.NewRecorder()

	// Call the CreateUserHandler function with the new request and response recorder
	handler := http.HandlerFunc(CreateUserHandler)
	handler.ServeHTTP(rr, req)
	requestBody2 := map[string]string{
		"id":       id,
		"password": pass,
	}
	requestBodyJson, _ = json.Marshal(requestBody2)
	req, err = http.NewRequest("POST", "/signin", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder to capture the handler's response
	rr = httptest.NewRecorder()

	// Call the AddNewUser function with the new request and response recorder
	handler2 := http.HandlerFunc(SigninExistingUser)
	handler2.ServeHTTP(rr, req)

	// Get the token
	_, tok, _ := storage.Db.FetchUserToken("testuser")
	// Add an example asset
	requestBody3 := models.Asset{
		ID:          "assetExample",
		Type:        "chart",
		Description: "example description",
		AssetData:   nil,
		Added:       time.Now(),
		Modified:    time.Now(),
	}
	requestBodyJson, _ = json.Marshal(requestBody3)
	req, err = http.NewRequest("POST", "/users/"+id+"/favourites", bytes.NewBuffer(requestBodyJson))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Set the request Content-Type header to JSON
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", tok)
	// Create a response recorder to capture the handler's response
	rr = httptest.NewRecorder()

	// Call the AddFavouriteHandler function with the new request and response recorder
	handler3 := http.HandlerFunc(AddFavouriteHandler)
	handler3.ServeHTTP(rr, req)
}
