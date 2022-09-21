package httpapi

import "github.com/labstack/echo/v4"

// adding fake controllers for creating swagger documentation

// @Summary      Add Insight
// @Description  Add new asset for insights
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        insight  body  domain.Insight  true  "insight"
// @Success      200  {object}  AssetInsightJson
// @Failure      401  {object}	ResponseStatus
// @Router       /api/v1/admin/insights [POST]
// @Security     BearerAuth
func (s *Server) addInsightHandler(c echo.Context) error {
	return nil
}

// @Summary      Add Chart
// @Description  Add new asset for charts
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        chart  body  domain.Chart  true  "chart"
// @Success      200  {object}  AssetChartJson
// @Failure      401  {object}	ResponseStatus
// @Router       /api/v1/admin/charts [POST]
// @Security     BearerAuth
func (s *Server) addChartHandler(c echo.Context) error {
	return nil
}

// @Summary      Add Audience
// @Description  Add new asset for audiences
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        audience  body  domain.Audience  true  "audience"
// @Success      200  {object}  AssetAudienceJson
// @Failure      401  {object}	ResponseStatus
// @Router       /api/v1/admin/audiences [POST]
// @Security     BearerAuth
func (s *Server) addAudienceHandler(c echo.Context) error {
	return nil
}

// @Summary      Update Insight
// @Description  Update an existing asset from insights
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Insight ID"
// @Param        insight  body  domain.Insight  true  "insight"
// @Success      200  {object}  AssetInsightJson
// @Failure      401  {object}	ResponseStatus
// @Router       /api/v1/admin/insights/{id} [PUT]
// @Security     BearerAuth
func (s *Server) updateInsightHandler(c echo.Context) error {
	return nil
}

// @Summary      Update Chart
// @Description  Update an existing asset from charts
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Chart ID"
// @Param        chart  body  domain.Chart  true  "chart"
// @Success      200  {object}  AssetChartJson
// @Failure      401  {object}	ResponseStatus
// @Router       /api/v1/admin/charts/{id} [PUT]
// @Security     BearerAuth
func (s *Server) updateChartHandler(c echo.Context) error {
	return nil
}

// @Summary      Update Audience
// @Description  Update an existing asset from audiences
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Audience ID"
// @Param        audience  body  domain.Audience  true  "audience"
// @Success      200  {object}  AssetAudienceJson
// @Failure      401  {object}	ResponseStatus
// @Router       /api/v1/admin/audiences/{id} [PUT]
// @Security     BearerAuth
func (s *Server) updateAudienceHandler(c echo.Context) error {
	return nil
}

// @Summary      Delete Insight
// @Description  Delete an asset from insights
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Insight ID"
// @Success      200  {object}  ResponseStatus
// @Failure      401  {object}	ResponseStatus
// @Router       /api/v1/admin/insights/{id} [DELETE]
// @Security     BearerAuth
func (s *Server) deleteInsightHandler(c echo.Context) error {
	return nil
}

// @Summary      Delete Chart
// @Description  Delete an asset from charts
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Chart ID"
// @Param        chart  body  domain.Chart  true  "chart"
// @Success      200  {object}  ResponseStatus
// @Failure      401  {object}	ResponseStatus
// @Router       /api/v1/admin/charts/{id} [DELETE]
// @Security     BearerAuth
func (s *Server) deleteChartHandler(c echo.Context) error {
	return nil
}

// @Summary      Delete Audience
// @Description  Delete an asset from audiences
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Audience ID"
// @Param        audience  body  domain.Audience  true  "audience"
// @Success      200  {object}  ResponseStatus
// @Failure      401  {object}	ResponseStatus
// @Router       /api/v1/admin/audiences/{id} [DELETE]
// @Security     BearerAuth
func (s *Server) deleteAudienceHandler(c echo.Context) error {
	return nil
}
