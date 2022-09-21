package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func createUserFromHttpTest(t *testing.T, server *Server, e *echo.Echo, submitInputJson string) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(submitInputJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expected := `{"status":"success"}`
	if assert.NoError(t, server.createUserHandler(c)) {
		assert.Equal(t, http.StatusOK, rec.Code, "Failed with "+submitInputJson)
		assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()), "Failed with "+submitInputJson)
	}
}

func loginFromHttpTest(t *testing.T, server *Server, e *echo.Echo, loginInput string) ResponseLogin {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(loginInput))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.NoError(t, server.loginUserHandler(c), "Failed with "+loginInput)
	assert.Equal(t, http.StatusOK, rec.Code, "Failed with "+loginInput)
	res := ResponseLogin{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	assert.NoError(t, err, "Failed with "+loginInput)
	assert.NotNil(t, res.Token, "Failed with "+loginInput)
	assert.NotNil(t, res.ExpiresAt, "Failed with "+loginInput)
	return res
}

func TestCreateUserSuccess(t *testing.T) {
	server, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	e := echo.New()

	users := []string{
		`{"username":"admin", "password":"pass", "isAdmin":true}`,
		`{"username":"user", "password":"pass", "isAdmin":false}`,
	}
	for _, v := range users {
		createUserFromHttpTest(t, server, e, v)
	}
}

func TestLoginUserSuccess(t *testing.T) {
	server, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	e := echo.New()
	users := map[int][]string{
		1: {`{"username":"admin", "password":"pass", "isAdmin":true}`, `{"username":"admin", "password":"pass"}`},
		2: {`{"username":"user", "password":"pass", "isAdmin":false}`, `{"username":"user", "password":"pass"}`},
	}

	config := middleware.JWTConfig{
		Claims:     &JwtUserClaims{},
		ContextKey: "user",
		SigningKey: []byte(server.secret),
	}
	e.Use(middleware.JWTWithConfig(config))
	for k, v := range users {
		userInput := v[0]
		loginInput := v[1]
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(userInput))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := `{"status":"success"}`
		if assert.NoError(t, server.createUserHandler(c)) {
			assert.Equal(t, http.StatusOK, rec.Code, k)
			assert.Equal(t, expected, strings.TrimSpace(rec.Body.String()), k)
		}

		res := loginFromHttpTest(t, server, e, loginInput)
		req = httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "bearer "+*res.Token)
		rec = httptest.NewRecorder()
		c = e.NewContext(req, rec)
		meHandler := middleware.JWTWithConfig(config)(server.meHandler)
		assert.NoError(t, meHandler(c), k)
		assert.Equal(t, http.StatusOK, rec.Code, k)
		fmt.Println(rec.Body.String())
	}
}

func TestUnauthorizedUserForAdminHandlers(t *testing.T) {
	server, teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	e := echo.New()
	config := middleware.JWTConfig{
		Claims:     &JwtUserClaims{},
		ContextKey: "user",
		SigningKey: []byte(server.secret),
	}
	e.Use(middleware.JWTWithConfig(config))

	submitInputJson := `{"username":"user", "password":"pass", "isAdmin":false}`
	loginInput := `{"username":"user", "password":"pass"}`
	createUserFromHttpTest(t, server, e, submitInputJson)
	res := loginFromHttpTest(t, server, e, loginInput)

	arr := []struct {
		handler func(c echo.Context) error
		name    string
	}{
		{server.addAssetHandler, "addAssetHandler"},
		{server.deleteAssetHandler, "deleteAssetHandler"},
		{server.updateAssetHandler, "updateAssetHandler"},
		{server.listUserFavouriteAssetsHandler, "listUserFavouriteAssetsHandler"},
	}
	for _, v := range arr {

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set(echo.HeaderAuthorization, "bearer "+*res.Token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		h := middleware.JWTWithConfig(config)(v.handler)
		he, ok := h(c).(*echo.HTTPError)
		if assert.True(t, ok, v.name) {
			assert.Error(t, he, v.name)
			assert.Equal(t, http.StatusUnauthorized, he.Code, v.name)
		}

	}

}
