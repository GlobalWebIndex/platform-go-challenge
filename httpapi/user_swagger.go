package httpapi

import "github.com/labstack/echo/v4"

// adding fake controllers for creating swagger documentation

// @Summary      Get Insight
// @Description  Get an insight based on ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Insight ID"
// @Success      200  {object}  AssetInsightJson
// @Failure      401  {object}	ResponseStatus
// @Failure      404  {object}	ResponseStatus
// @Router       /api/v1/insights/{id} [GET]
// @Security     BearerAuth
func (s *Server) getInsightHandler(c echo.Context) error {
	return nil
}

// @Summary      Get Chart
// @Description  Get a chart based on ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Chart ID"
// @Success      200  {object}  AssetChartJson
// @Failure      401  {object}	ResponseStatus
// @Failure      404  {object}	ResponseStatus
// @Router       /api/v1/charts/{id} [GET]
// @Security     BearerAuth
func (s *Server) getChartHandler(c echo.Context) error {
	return nil
}

// @Summary      Get Audience
// @Description  Get an audience based on ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Audience ID"
// @Success      200  {object}  AssetAudienceJson
// @Failure      401  {object}	ResponseStatus
// @Failure      404  {object}	ResponseStatus
// @Router       /api/v1/audiences/{id} [GET]
// @Security     BearerAuth
func (s *Server) getAudienceHandler(c echo.Context) error {
	return nil
}

// @Summary      Favour an Insight
// @Description  Favour an insight based on ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Insight ID"
// @Success      200  {object}  ResponseStatus
// @Failure      401  {object}	ResponseStatus
// @Failure      404  {object}	ResponseStatus
// @Router       /api/v1/insights/{id}/favourite [PUT]
// @Security     BearerAuth
func (s *Server) favourInsightHandler(c echo.Context) error {
	return nil
}

// @Summary      Favour Chart
// @Description  Favour a chart based on ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Chart ID"
// @Success      200  {object}  ResponseStatus
// @Failure      401  {object}	ResponseStatus
// @Failure      404  {object}	ResponseStatus
// @Router       /api/v1/charts/{id}/favourite [PUT]
// @Security     BearerAuth
func (s *Server) favourChartHandler(c echo.Context) error {
	return nil
}

// @Summary      Favour Audience
// @Description  Favour an audience based on ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Audience ID"
// @Success      200  {object}  ResponseStatus
// @Failure      401  {object}	ResponseStatus
// @Failure      404  {object}	ResponseStatus
// @Router       /api/v1/audiences/{id}/favourite [PUT]
// @Security     BearerAuth
func (s *Server) favourAudienceHandler(c echo.Context) error {
	return nil
}
