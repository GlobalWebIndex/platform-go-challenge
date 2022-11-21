package favourite

import (
	"net/http"
	"platform-go-challenge/internal/app/favouriteasset"

	"github.com/labstack/echo/v4"
)

// Router infrastructure definition.
type Router struct {
	asSvc favouriteasset.Service
}

// NewRouter returns an HTTP component to serve all the routes for the user favourites.
func NewRouter(asSvc favouriteasset.Service) *Router {
	return &Router{
		asSvc: asSvc,
	}
}

// AppendRoutes adds favourite assets routes to router.
func (r *Router) AppendRoutes(e *echo.Echo) {
	e.GET("/users/:user_id/favourites", r.GetUserFavourites)
	e.POST("/users/:user_id/favourites/:asset_id", r.AddAssetToFavourites)
	e.PUT("/users/:user_id/favourites/:favourite_asset_id", r.EditFavouriteAsset)
	e.DELETE("/users/:user_id/favourites/:favourite_asset_id", r.RemoveAssetFromFavourites)
}

// GetUserFavourites gets the favourite assets of a user.
func (r *Router) GetUserFavourites(c echo.Context) error {
	userID := c.Param("user_id")

	assets, err := r.asSvc.GetFavouriteAssets(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, assets)
}

// AddAssetToFavourites gets the favourite assets of a user.
func (r *Router) AddAssetToFavourites(c echo.Context) error {
	//TODO: Validate inputs
	userID := c.Param("user_id")
	assetID := c.Param("asset_id")

	res, err := r.asSvc.AddAssetToFavourites(c.Request().Context(), userID, assetID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, res)
}

// EditFavouriteAsset edits a favourite asset.
func (r *Router) EditFavouriteAsset(c echo.Context) error {
	// TODO: Validate inputs
	userID := c.Param("user_id")
	fAssetID := c.Param("favourite_asset_id")
	body := new(favouriteasset.EditFavouriteAsset)
	if err := c.Bind(body); err != nil {
		return err
	}
	err := r.asSvc.EditFavouriteAsset(c.Request().Context(), userID, fAssetID, *body)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

// RemoveAssetFromFavourites removes an asset from user's favourites.
func (r *Router) RemoveAssetFromFavourites(c echo.Context) error {
	userID := c.Param("user_id")
	fAssetID := c.Param("favourite_asset_id")
	err := r.asSvc.RemoveAssetFromFavourites(c.Request().Context(), userID, fAssetID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}
