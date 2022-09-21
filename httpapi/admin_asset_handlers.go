package httpapi

import (
	"net/http"
	"platform-go-challenge/domain"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) addAssetHandler(c echo.Context) error {
	user, err := getUserDomain(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
			"error":  err.Error(),
		})
	}
	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
		})
	}

	at := c.Param("assetType")
	var assetData interface{}
	switch at {
	case AssetTypeInsights:
		in := domain.Insight{}
		err := c.Bind(&in)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseStatus{
				Status: FailureStatus,
				Error:  err.Error(),
			})
		}
		assetData = &in
	case AssetTypeCharts:
		in := domain.Chart{}
		err := c.Bind(&in)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseStatus{
				Status: FailureStatus,
				Error:  err.Error(),
			})
		}
		assetData = &in
	case AssetTypeAudiences:
		in := domain.Audience{}
		err := c.Bind(&in)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseStatus{
				Status: FailureStatus,
				Error:  err.Error(),
			})
		}
		assetData = &in
	}
	asset := domain.InputAsset{
		Data: assetData,
	}
	newAsset, err := s.domain.AddAsset(c.Request().Context(), user, asset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}

	return c.JSON(http.StatusOK, newAsset)
}

func (s *Server) deleteAssetHandler(c echo.Context) error {
	user, err := getUserDomain(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
			"error":  err.Error(),
		})
	}
	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
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

	err = s.domain.DeleteAsset(c.Request().Context(), user, uint(assetId), assetType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status": SuccessStatus,
	})
}

func (s *Server) updateAssetHandler(c echo.Context) error {
	user, err := getUserDomain(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
			"error":  err.Error(),
		})
	}
	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
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
	var assetData interface{}
	switch at {
	case AssetTypeInsights:
		in := domain.Insight{}
		err := c.Bind(&in)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseStatus{
				Status: FailureStatus,
				Error:  err.Error(),
			})
		}
		assetData = &in
	case AssetTypeCharts:
		in := domain.Chart{}
		err := c.Bind(&in)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseStatus{
				Status: FailureStatus,
				Error:  err.Error(),
			})
		}
		assetData = &in
	case AssetTypeAudiences:
		in := domain.Audience{}
		err := c.Bind(&in)
		if err != nil {
			return c.JSON(http.StatusBadRequest, ResponseStatus{
				Status: FailureStatus,
				Error:  err.Error(),
			})
		}
		assetData = &in
	}
	asset := domain.InputAsset{
		Data: assetData,
	}

	newAsset, err := s.domain.UpdateAsset(c.Request().Context(), user, uint(assetId), asset)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseStatus{
			Status: FailureStatus,
			Error:  err.Error(),
		})
	}
	return c.JSON(http.StatusOK, newAsset)
}

func (s *Server) listUserFavouriteAssetsHandler(c echo.Context) error {
	user, err := getUserDomain(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
			"error":  err.Error(),
		})
	}
	if !user.IsAdmin {
		return echo.NewHTTPError(http.StatusUnauthorized, echo.Map{
			"status": "Unauthorized",
		})
	}
	return nil
}
