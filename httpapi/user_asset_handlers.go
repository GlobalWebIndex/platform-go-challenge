package httpapi

import (
	"errors"
	"net/http"
	"platform-go-challenge/domain"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) getAssetHandler(c echo.Context) error {
	user, err := getUserDomain(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
			"error":  err.Error(),
		})
	}
	idStr := c.Param("id")
	assetId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseStatus{
			Status: FailureStatus,
			Error:  "asset ID not a number",
		})
	}
	at := c.Param("assetType")
	var assetType domain.AssetType
	switch at {
	case AssetTypeInsights:
		assetType = domain.InsightAssetType
	case AssetTypeCharts:
		assetType = domain.ChartAssetType
	case AssetTypeAudiences:
		assetType = domain.AudienceAssetType
	}
	asset, err := s.domain.GetAsset(c.Request().Context(), user, uint(assetId), assetType)
	if err != nil {
		return c.JSON(http.StatusNotFound, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, asset)
}

// @Summary      List of assets
// @Description  Get list of assets based on the asset type, the number of assets in the page and the last ID to start counting
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        query  body  QueryAssets  true  "query options"
// @Success      200  {object}  ListInsightsJson
// @Failure      400  {object}	ResponseStatus
// @Failure      401  {object}	ResponseStatus
// @Router       /api/v1/assets [POST]
// @Security     BearerAuth
func (s *Server) listAssetsHandler(c echo.Context) error {
	user, err := getUserDomain(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
			"error":  err.Error(),
		})
	}
	query := QueryAssets{}
	err = c.Bind(&query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}

	ls, err := s.domain.ListAssets(c.Request().Context(), user, query.QueryAssets, query.Who)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
				"status": "Unauthorized",
				"error":  err.Error(),
			})
		}
		return c.JSON(http.StatusBadRequest, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, ls)
}

func (s *Server) favourAnAssetHandler(c echo.Context) error {
	user, err := getUserDomain(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
			"error":  err.Error(),
		})
	}
	idStr := c.Param("id")
	assetId, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseStatus{
			Status: FailureStatus,
			Error:  "asset ID not a number",
		})
	}
	at := c.Param("assetType")
	var assetType domain.AssetType
	switch at {
	case AssetTypeInsights:
		assetType = domain.InsightAssetType
	case AssetTypeCharts:
		assetType = domain.ChartAssetType
	case AssetTypeAudiences:
		assetType = domain.AudienceAssetType
	}

	isFavourite := true
	if c.Request().Method == "DELETE" {
		isFavourite = false
	}
	err = s.domain.FavouriteAsset(c.Request().Context(), user, uint(assetId), assetType, isFavourite)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status": SuccessStatus,
	})
}

// @Summary      List of favourite assets
// @Description  Get list of favourite assets of the user based on the asset type, the number of assets in the page and the last ID to start counting
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        query  body  domain.QueryAssets  true  "query options"
// @Success      200  {object}  ListChartsJson
// @Failure      400  {object}	ResponseStatus
// @Failure      401  {object}	ResponseStatus
// @Router       /api/v1/me/favourites [POST]
// @Security     BearerAuth
func (s *Server) listMyFavourites(c echo.Context) error {
	user, err := getUserDomain(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
			"error":  err.Error(),
		})
	}
	query := domain.QueryAssets{}
	err = c.Bind(&query)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}
	favQur := domain.QueryFavouriteAssets{
		FromUserID: user.ID,
		OnlyFav:    true,
	}
	ls, err := s.domain.ListAssets(c.Request().Context(), user, query, &favQur)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, ls)
}
