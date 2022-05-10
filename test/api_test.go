package gwi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	gwi "github.com/josedelrio85/platform-go-challenge/internal"
)

func TestGetFavsFromUser(t *testing.T) {

	assert := assert.New(t)

	uuidval := uuid.New().String()

	tests := []struct {
		Description      string
		Url              string
		Input            map[string]string
		Payload          interface{}
		ExpectedResponse int
		ExpectedSucceed  bool
	}{
		{
			Description:      "when there is no uuid",
			Url:              "foobar",
			ExpectedResponse: http.StatusUnprocessableEntity,
			ExpectedSucceed:  false,
		},
		{
			Description:      "when there is uuid, but it is not valid",
			Url:              "user/foobar/fav/list",
			ExpectedResponse: http.StatusUnprocessableEntity,
			ExpectedSucceed:  false,
		},
		{
			Description: "when there is uuid, and it is valid",
			Url:         fmt.Sprintf("user/%s/fav/list", uuidval),
			Input: map[string]string{
				"userid": uuidval,
			},
			ExpectedResponse: http.StatusNoContent,
			ExpectedSucceed:  true,
		},
	}

	mer, err := gwi.NewMemoryRepository()
	assert.NoError(err)

	handler := gwi.Handler{
		Repo: mer,
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			}))
			defer func() { testServer.Close() }()

			jsonStr, err := json.Marshal(test.Payload)
			assert.NoError(err)

			url := fmt.Sprintf("%s/%s", testServer.URL, test.Url)
			log.Println(url)
			req := httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonStr))
			req = mux.SetURLVars(req, test.Input)
			res := httptest.NewRecorder()

			h := handler.GetFavsFromUser()
			h.ServeHTTP(res, req)
			assert.Equal(res.Result().StatusCode, test.ExpectedResponse)
		})
	}
}

func TestAddNewFav(t *testing.T) {

	assert := assert.New(t)

	uuidval := uuid.New().String()

	tests := []struct {
		Description      string
		Url              string
		Input            map[string]string
		Payload          interface{}
		ExpectedResponse int
		ExpectedSucceed  bool
	}{
		{
			Description:      "when there is no uuid",
			Url:              "foobar",
			ExpectedResponse: http.StatusUnprocessableEntity,
			ExpectedSucceed:  false,
		},
		{
			Description:      "when there is uuid, but it is not valid",
			Url:              "user/foobar/fav/add",
			ExpectedResponse: http.StatusUnprocessableEntity,
			ExpectedSucceed:  false,
		},
		{
			Description: "when there is uuid, and it is valid",
			Url:         fmt.Sprintf("user/%s/fav/add", uuidval),
			Input: map[string]string{
				"userid": uuidval,
			},
			Payload: gwi.Asset{
				AssetType: gwi.TypeInsight,
				Insight: gwi.Insight{
					Id:          uuid.New().String(),
					Description: "",
				},
			},
			ExpectedResponse: http.StatusOK,
			ExpectedSucceed:  true,
		},
	}

	mer, err := gwi.NewMemoryRepository()
	assert.NoError(err)

	handler := gwi.Handler{
		Repo: mer,
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			}))
			defer func() { testServer.Close() }()

			jsonStr, err := json.Marshal(test.Payload)
			assert.NoError(err)

			url := fmt.Sprintf("%s/%s", testServer.URL, test.Url)
			log.Println(url)
			req := httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonStr))
			req = mux.SetURLVars(req, test.Input)
			res := httptest.NewRecorder()

			h := handler.AddNewFav()
			h.ServeHTTP(res, req)
			assert.Equal(res.Result().StatusCode, test.ExpectedResponse)
		})
	}
}

func TestApiUpdateFav(t *testing.T) {

	assert := assert.New(t)

	uuidval := uuid.New().String()

	tests := []struct {
		Description      string
		Url              string
		Input            map[string]string
		Payload          interface{}
		ExpectedResponse int
		ExpectedSucceed  bool
	}{
		{
			Description:      "when there is no uuid",
			Url:              "foobar",
			ExpectedResponse: http.StatusUnprocessableEntity,
			ExpectedSucceed:  false,
		},
		{
			Description:      "when there is uuid, but it is not valid",
			Url:              "user/foobar/fav/add",
			ExpectedResponse: http.StatusUnprocessableEntity,
			ExpectedSucceed:  false,
		},
		{
			Description: "when updated entity does not exist",
			Url:         fmt.Sprintf("user/%s/fav/add", uuidval),
			Input: map[string]string{
				"userid": uuidval,
			},
			Payload: gwi.Asset{
				AssetType: gwi.TypeInsight,
				Insight: gwi.Insight{
					Id:          uuid.New().String(),
					Description: "",
				},
			},
			ExpectedResponse: http.StatusNoContent,
			ExpectedSucceed:  true,
		},
	}

	mer, err := gwi.NewMemoryRepository()
	assert.NoError(err)

	handler := gwi.Handler{
		Repo: mer,
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			}))
			defer func() { testServer.Close() }()

			jsonStr, err := json.Marshal(test.Payload)
			assert.NoError(err)

			url := fmt.Sprintf("%s/%s", testServer.URL, test.Url)
			log.Println(url)
			req := httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonStr))
			req = mux.SetURLVars(req, test.Input)
			res := httptest.NewRecorder()

			h := handler.UpdateFav()
			h.ServeHTTP(res, req)
			assert.Equal(res.Result().StatusCode, test.ExpectedResponse)
		})
	}
}

func TestDeleteFav(t *testing.T) {

	assert := assert.New(t)

	uuidval := uuid.New().String()

	tests := []struct {
		Description      string
		Url              string
		Input            map[string]string
		Payload          interface{}
		ExpectedResponse int
		ExpectedSucceed  bool
	}{
		{
			Description:      "when there is no uuid",
			Url:              "foobar",
			ExpectedResponse: http.StatusUnprocessableEntity,
			ExpectedSucceed:  false,
		},
		{
			Description:      "when there is uuid, but it is not valid",
			Url:              "user/foobar/fav/add",
			ExpectedResponse: http.StatusUnprocessableEntity,
			ExpectedSucceed:  false,
		},
		{
			Description: "when deleted entity does not exist",
			Url:         fmt.Sprintf("user/%s/fav/add", uuidval),
			Input: map[string]string{
				"userid": uuidval,
			},
			Payload: gwi.Asset{
				AssetType: gwi.TypeInsight,
				Insight: gwi.Insight{
					Id:          uuid.New().String(),
					Description: "",
				},
			},
			ExpectedResponse: http.StatusNoContent,
			ExpectedSucceed:  true,
		},
	}

	mer, err := gwi.NewMemoryRepository()
	assert.NoError(err)

	handler := gwi.Handler{
		Repo: mer,
	}

	for _, test := range tests {
		t.Run(test.Description, func(t *testing.T) {

			testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			}))
			defer func() { testServer.Close() }()

			jsonStr, err := json.Marshal(test.Payload)
			assert.NoError(err)

			url := fmt.Sprintf("%s/%s", testServer.URL, test.Url)
			log.Println(url)
			req := httptest.NewRequest(http.MethodGet, url, bytes.NewBuffer(jsonStr))
			req = mux.SetURLVars(req, test.Input)
			res := httptest.NewRecorder()

			h := handler.DeleteFav()
			h.ServeHTTP(res, req)
			assert.Equal(res.Result().StatusCode, test.ExpectedResponse)
		})
	}
}
