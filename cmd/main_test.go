package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"svenvermeulen/platform-go-challenge/pkg/model"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHappyPath(t *testing.T) {
	router := SetupRouter()
	userId := uuid.New()

	// GIVEN a user with a favourite item
	f := model.UserFavouriteShort{
		Description:  "favourite description",
		ResourceType: "chart",
		Id:           uuid.New(),
	}

	bodyBytes, err := json.Marshal(f)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewBuffer(bodyBytes)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/favourites", body)
	req.Header.Add("Authorization", getAuthToken(userId))
	require.NoError(t, err, "error creating http request")
	router.ServeHTTP(w, req)
	require.Equal(t, 201, w.Code)

	// WHEN I retrieve the user's favourites
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/favourites", nil)
	req.Header.Add("Authorization", getAuthToken(userId))
	require.NoError(t, err, "error creating http request")
	router.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	result := []model.UserFavourite{}
	if err := json.NewDecoder(w.Result().Body).Decode(&result); err != nil {
		t.Fatalf("could not decode response body: %s", err)
	}

	// THEN I get a list of audiences, insights and charts
	assert.Equal(t, 1, len(result))
	assert.Equal(t, f.Id, result[0].Chart.Id)
}

func TestDeletion(t *testing.T) {
	router := SetupRouter()
	userId := uuid.New()

	var uuidToDelete uuid.UUID

	// GIVEN a user with a few favourite items
	for i := range 2 {
		f := model.UserFavouriteShort{
			Description:  fmt.Sprintf("favourite %d description", i),
			ResourceType: "chart",
			Id:           uuid.New(),
		}
		uuidToDelete = f.Id

		bodyBytes, err := json.Marshal(f)
		if err != nil {
			log.Fatal(err)
		}

		body := bytes.NewBuffer(bodyBytes)

		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/favourites", body)
		req.Header.Add("Authorization", getAuthToken(userId))
		require.NoError(t, err, "error creating http request")
		router.ServeHTTP(w, req)
		require.Equal(t, 201, w.Code)
	}
	// WHEN I delete one favourite
	w := httptest.NewRecorder()
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/favourites/%v", uuidToDelete), nil)
	req.Header.Add("Authorization", getAuthToken(userId))
	require.NoError(t, err, "error creating http request")
	router.ServeHTTP(w, req)
	require.Equal(t, 204, w.Code)

	// WHEN I retrieve the user's favourites
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/favourites", nil)
	req.Header.Add("Authorization", getAuthToken(userId))
	require.NoError(t, err, "error creating http request")
	router.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	result := []model.UserFavourite{}
	if err := json.NewDecoder(w.Result().Body).Decode(&result); err != nil {
		t.Fatalf("could not decode response body: %s", err)
	}

	// THEN I get a list 1 favourite
	assert.Equal(t, 1, len(result))
}

func TestUpdate(t *testing.T) {
	router := SetupRouter()
	userId := uuid.New()

	// GIVEN a user with a favourite item
	f := model.UserFavouriteShort{
		Description:  "initial description",
		ResourceType: "chart",
		Id:           uuid.New(),
	}

	bodyBytes, err := json.Marshal(f)
	if err != nil {
		log.Fatal(err)
	}

	body := bytes.NewBuffer(bodyBytes)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/favourites", body)
	req.Header.Add("Authorization", getAuthToken(userId))
	require.NoError(t, err, "error creating http request")
	router.ServeHTTP(w, req)
	require.Equal(t, 201, w.Code)

	// WHEN I update the description
	w = httptest.NewRecorder()
	f.Description = "NEW DESCRIPTION"
	bodyBytes, err = json.Marshal(f)
	if err != nil {
		log.Fatal(err)
	}

	body = bytes.NewBuffer(bodyBytes)

	req, err = http.NewRequest("PATCH", fmt.Sprintf("/favourites/%v", f.Id), body)
	req.Header.Add("Authorization", getAuthToken(userId))
	require.NoError(t, err, "error creating http request")
	router.ServeHTTP(w, req)
	require.Equal(t, 204, w.Code)

	// WHEN I retrieve the user's favourites
	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/favourites", nil)
	req.Header.Add("Authorization", getAuthToken(userId))
	require.NoError(t, err, "error creating http request")
	router.ServeHTTP(w, req)
	require.Equal(t, 200, w.Code)

	result := []model.UserFavourite{}
	if err := json.NewDecoder(w.Result().Body).Decode(&result); err != nil {
		t.Fatalf("could not decode response body: %s", err)
	}

	// THEN I get a list with 1 favourite, with the correct description
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "NEW DESCRIPTION", result[0].Description)
}

func getAuthToken(userId uuid.UUID) string {
	claims := jwt.MapClaims{
		"subject": "favouritessvc",
		"userid":  userId.String(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte("12345678123456781234567812345678"))

	if err != nil {
		log.Errorf("Error generating token: %v\n", err)
		return ""
	}

	log.Debugf("Signed token: %v", signed)
	return fmt.Sprintf("Bearer %s", signed)
}
