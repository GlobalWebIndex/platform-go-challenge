package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/artemis13/platform-go-challenge/handlers"
	"github.com/artemis13/platform-go-challenge/models"
	"github.com/artemis13/platform-go-challenge/storage"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestData() {
	storage.ClearUsers()
	storage.AddUser(&models.User{
		ID: 1,
		Favorites: []models.Asset{
			{ID: "1", Description: "Favorite 1"},
			{ID: "2", Description: "Favorite 2"},
			{ID: "3", Description: "Favorite 3"},
			{ID: "4", Description: "Favorite 4"},
			{ID: "5", Description: "Favorite 5"},
			{ID: "6", Description: "Favorite 6"},
			{ID: "7", Description: "Favorite 7"},
			{ID: "8", Description: "Favorite 8"},
			{ID: "9", Description: "Favorite 9"},
			{ID: "10", Description: "Favorite 10"},
			{ID: "11", Description: "Favorite 11"},
		},
	})
}

func TestGetUserFavorites(t *testing.T) { //test with no pager
	e := echo.New()
	setupTestData()

	req := httptest.NewRequest(http.MethodGet, "/users/1/favorites", nil)
	req.Header.Set(echo.HeaderAuthorization, "gwi-token-12345")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id/favorites")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handlers.GetUserFavorites(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		expected := `{"page":1,"limit":10,"total":11,"favorites":[{"ID":"1","Type":0,"Description":"Favorite 1","Chart":null,"Insight":null,"Audience":null},{"ID":"2","Type":0,"Description":"Favorite 2","Chart":null,"Insight":null,"Audience":null},{"ID":"3","Type":0,"Description":"Favorite 3","Chart":null,"Insight":null,"Audience":null},{"ID":"4","Type":0,"Description":"Favorite 4","Chart":null,"Insight":null,"Audience":null},{"ID":"5","Type":0,"Description":"Favorite 5","Chart":null,"Insight":null,"Audience":null},{"ID":"6","Type":0,"Description":"Favorite 6","Chart":null,"Insight":null,"Audience":null},{"ID":"7","Type":0,"Description":"Favorite 7","Chart":null,"Insight":null,"Audience":null},{"ID":"8","Type":0,"Description":"Favorite 8","Chart":null,"Insight":null,"Audience":null},{"ID":"9","Type":0,"Description":"Favorite 9","Chart":null,"Insight":null,"Audience":null},{"ID":"10","Type":0,"Description":"Favorite 10","Chart":null,"Insight":null,"Audience":null}]}`

		assert.JSONEq(t, expected, rec.Body.String())
	}
}

func TestGetUserFavoritesPager(t *testing.T) { //test with 1 pager
	//to be implemented
}
func TestAddUserFavorite(t *testing.T) {
	e := echo.New()
	setupTestData()

	assetJSON := `{"Type":0,"Description":"New Favorite"}`
	req := httptest.NewRequest(http.MethodPost, "/users/1/favorites", strings.NewReader(assetJSON))
	req.Header.Set(echo.HeaderAuthorization, "gwi-token-12345")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id/favorites")
	c.SetParamNames("id")
	c.SetParamValues("1")

	if assert.NoError(t, handlers.AddUserFavorite(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		var asset models.Asset
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &asset))
		assert.NotEmpty(t, asset.ID)
		assert.Equal(t, "New Favorite", asset.Description)
	}
}

func TestRemoveUserFavorite(t *testing.T) {
	e := echo.New()
	setupTestData()

	req := httptest.NewRequest(http.MethodDelete, "/users/1/favorites/1", nil)
	req.Header.Set(echo.HeaderAuthorization, "gwi-token-12345")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id/favorites/:asset_id")
	c.SetParamNames("id", "asset_id")
	c.SetParamValues("1", "1")

	if assert.NoError(t, handlers.RemoveUserFavorite(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

func TestEditUserFavorite(t *testing.T) {
	e := echo.New()
	setupTestData()

	newDescription := `{"description":"Updated Favorite"}`
	req := httptest.NewRequest(http.MethodPut, "/users/1/favorites/1", strings.NewReader(newDescription))
	req.Header.Set(echo.HeaderAuthorization, "gwi-token-12345")
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:id/favorites/:asset_id")
	c.SetParamNames("id", "asset_id")
	c.SetParamValues("1", "1")

	if assert.NoError(t, handlers.EditUserFavorite(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var asset models.Asset
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &asset))
		assert.Equal(t, "Updated Favorite", asset.Description)
	}
}
