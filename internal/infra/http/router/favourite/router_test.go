package favourite_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"platform-go-challenge/internal/app/favouriteasset"
	"platform-go-challenge/internal/infra/http/router/favourite"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRouter_GetUserFavourites(t *testing.T) {
	favouriteAssetID := "6378a0ae5893f7f7cc2770a9"
	assetID := "6378a00d5893f7f7cc2770a3"
	userID := "6377c9c68c9f75c3d2a8a5ba"

	favouriteAsset := &[]favouriteasset.FavouriteAsset{
		{
			ID:          &favouriteAssetID,
			Description: "desc",
			Asset: &favouriteasset.Asset{
				ID:   assetID,
				Type: "Audience",
				Data: "test data",
			},
		},
	}
	favouriteAssetRes := &favouriteasset.GetFavouriteAssetsRes{
		FavouriteAssets: favouriteAsset,
	}
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc *favouriteasset.SvcMock
		echoCtx echo.Context
		expRes  string
		expErr  error
	}{
		"Should succeed on GetFavouriteAssets call ": {
			mockSvc: func() *favouriteasset.SvcMock {
				mockSvc := &favouriteasset.SvcMock{}
				mockSvc.On("GetFavouriteAssets", context.TODO(), userID).
					Return(favouriteAssetRes, nil)

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				c := e.NewContext(req, rec)
				c.SetPath("/users/:user_id/favourites")
				c.SetParamNames("user_id")
				c.SetParamValues(userID)

				return c
			}(),
			expRes: func() string {
				b, err := json.Marshal(favouriteAssetRes)
				assert.Nil(t, err)

				return string(b) + "\n"
			}(),
		},
		"Should return error on GetFavouriteAssets error ": {
			mockSvc: func() *favouriteasset.SvcMock {
				mockSvc := &favouriteasset.SvcMock{}
				mockSvc.On("GetFavouriteAssets", context.TODO(), userID).
					Return(&favouriteasset.GetFavouriteAssetsRes{}, errors.New("random error"))

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				c := e.NewContext(req, rec)
				c.SetPath("/users/:user_id/favourites")
				c.SetParamNames("user_id")
				c.SetParamValues(userID)

				return c
			}(),
			expErr: errors.New("random error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			r := favourite.NewRouter(tt.mockSvc)
			err := r.GetUserFavourites(tt.echoCtx)

			if tt.expErr == nil {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.Equal(t, tt.expRes, rec.Body.String())
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expErr, err)
			}
		})
	}
}

func TestRouter_AddAssetToFavourites(t *testing.T) {
	favouriteAssetID := "6378a0ae5893f7f7cc2770a9"
	assetID := "6378a00d5893f7f7cc2770a3"
	userID := "6377c9c68c9f75c3d2a8a5ba"
	addFavouriteAssetRes := &favouriteasset.AddFavouriteAssetRes{
		ID: favouriteAssetID,
	}
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc *favouriteasset.SvcMock
		echoCtx echo.Context
		expRes  string
		expErr  error
	}{
		"Should succeed on AddAssetToFavourites call ": {
			mockSvc: func() *favouriteasset.SvcMock {
				mockSvc := &favouriteasset.SvcMock{}
				mockSvc.On("AddAssetToFavourites", context.TODO(), userID, assetID).
					Return(addFavouriteAssetRes, nil)

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				c := e.NewContext(req, rec)
				c.SetPath("/users/:user_id/favourites/:asset_id")
				c.SetParamNames("user_id", "asset_id")
				c.SetParamValues(userID, assetID)

				return c
			}(),
			expRes: func() string {
				b, err := json.Marshal(addFavouriteAssetRes)
				assert.Nil(t, err)

				return string(b) + "\n"
			}(),
		},
		"Should return error on AddAssetToFavourites error ": {
			mockSvc: func() *favouriteasset.SvcMock {
				mockSvc := &favouriteasset.SvcMock{}
				mockSvc.On("AddAssetToFavourites", context.TODO(), userID, assetID).
					Return(&favouriteasset.AddFavouriteAssetRes{}, errors.New("random error"))

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPost, "/", nil)
				c := e.NewContext(req, rec)
				c.SetPath("/users/:user_id/favourites/:asset_id")
				c.SetParamNames("user_id", "asset_id")
				c.SetParamValues(userID, assetID)

				return c
			}(),
			expErr: errors.New("random error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			r := favourite.NewRouter(tt.mockSvc)
			err := r.AddAssetToFavourites(tt.echoCtx)

			if tt.expErr == nil {
				assert.Equal(t, http.StatusCreated, rec.Code)
				assert.Equal(t, tt.expRes, rec.Body.String())
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expErr, err)
			}
		})
	}
}

func TestRouter_EditFavouriteAsset(t *testing.T) {
	favouriteAssetID := "6378a0ae5893f7f7cc2770a9"
	userID := "6377c9c68c9f75c3d2a8a5ba"
	editFavouriteAsset := favouriteasset.EditFavouriteAsset{
		Description: "test description",
	}
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc *favouriteasset.SvcMock
		echoCtx echo.Context
		expErr  error
	}{
		"Should succeed on EditFavouriteAsset call ": {
			mockSvc: func() *favouriteasset.SvcMock {
				mockSvc := &favouriteasset.SvcMock{}
				mockSvc.On("EditFavouriteAsset", context.TODO(), userID, favouriteAssetID, editFavouriteAsset).
					Return(nil)

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"description": "test description"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath("/users/:user_id/favourites/:favourite_asset_id")
				c.SetParamNames("user_id", "favourite_asset_id")
				c.SetParamValues(userID, favouriteAssetID)

				return c
			}(),
		},
		"Should return error on EditFavouriteAsset error ": {
			mockSvc: func() *favouriteasset.SvcMock {
				mockSvc := &favouriteasset.SvcMock{}
				mockSvc.On("EditFavouriteAsset", context.TODO(), userID, favouriteAssetID, editFavouriteAsset).
					Return(errors.New("random error"))

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"description": "test description"}`))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				c := e.NewContext(req, rec)
				c.SetPath("/users/:user_id/favourites/:favourite_asset_id")
				c.SetParamNames("user_id", "favourite_asset_id")
				c.SetParamValues(userID, favouriteAssetID)

				return c
			}(),
			expErr: errors.New("random error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			r := favourite.NewRouter(tt.mockSvc)
			err := r.EditFavouriteAsset(tt.echoCtx)
			if tt.expErr == nil {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expErr, err)
			}

		})
	}
}

func TestRouter_RemoveAssetFromFavourites(t *testing.T) {
	favouriteAssetID := "6378a0ae5893f7f7cc2770a9"
	userID := "6377c9c68c9f75c3d2a8a5ba"
	rec := httptest.NewRecorder()

	tests := map[string]struct {
		mockSvc *favouriteasset.SvcMock
		echoCtx echo.Context
		expErr  error
	}{
		"Should succeed on RemoveAssetFromFavourites call ": {
			mockSvc: func() *favouriteasset.SvcMock {
				mockSvc := &favouriteasset.SvcMock{}
				mockSvc.On("RemoveAssetFromFavourites", context.TODO(), userID, favouriteAssetID).
					Return(nil)

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodDelete, "/", nil)
				c := e.NewContext(req, rec)
				c.SetPath("/users/:user_id/favourites/:favourite_asset_id")
				c.SetParamNames("user_id", "favourite_asset_id")
				c.SetParamValues(userID, favouriteAssetID)

				return c
			}(),
		},
		"Should return error on RemoveAssetFromFavourites error ": {
			mockSvc: func() *favouriteasset.SvcMock {
				mockSvc := &favouriteasset.SvcMock{}
				mockSvc.On("RemoveAssetFromFavourites", context.TODO(), userID, favouriteAssetID).
					Return(errors.New("random error"))

				return mockSvc
			}(),
			echoCtx: func() echo.Context {
				e := echo.New()
				req := httptest.NewRequest(http.MethodDelete, "/", nil)
				c := e.NewContext(req, rec)
				c.SetPath("/users/:user_id/favourites/:favourite_asset_id")
				c.SetParamNames("user_id", "favourite_asset_id")
				c.SetParamValues(userID, favouriteAssetID)

				return c
			}(),
			expErr: errors.New("random error"),
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {

			r := favourite.NewRouter(tt.mockSvc)
			err := r.RemoveAssetFromFavourites(tt.echoCtx)
			if tt.expErr == nil {
				assert.Equal(t, http.StatusOK, rec.Code)
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expErr, err)
			}

		})
	}
}
