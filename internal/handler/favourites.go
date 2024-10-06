package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"svenvermeulen/platform-go-challenge/internal/auth"
	"svenvermeulen/platform-go-challenge/internal/repository/audience"
	"svenvermeulen/platform-go-challenge/internal/repository/chart"
	"svenvermeulen/platform-go-challenge/internal/repository/favourite"
	"svenvermeulen/platform-go-challenge/internal/repository/insight"
	"svenvermeulen/platform-go-challenge/pkg/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type FavouritesHandler struct {
	favouriteRepository *favourite.Repository
	audienceRepository  *audience.Repository
	chartRepository     *chart.Repository
	insightRepository   *insight.Repository
}

func NewFavouritesHandler(favouriteRepository *favourite.Repository,
	audienceRepository *audience.Repository,
	chartRepository *chart.Repository,
	insightRepository *insight.Repository) *FavouritesHandler {
	return &FavouritesHandler{
		favouriteRepository: favouriteRepository,
		audienceRepository:  audienceRepository,
		chartRepository:     chartRepository,
		insightRepository:   insightRepository,
	}
}

// GetUserFavourites godoc
//
// @Summary     Get user favourites
// @Description gets list of favourites for a specified user
// @Tags        favourites
// @Produce     json
// @Success     200	{object}	[]model.UserFavourite
// @Failure     400
// @Failure     401
// @Failure     404
// @Failure     500
// @Router      /favourites/:userid [get]
func (h *FavouritesHandler) GetFavourites(c *gin.Context) {
	userId, err := auth.GetUserIDFromToken(c)
	if err != nil {
		log.Errorf("Error obtaining userid from jwt token: %v\n", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	log.Info("Retrieving favourites for userid", userId)

	var offset int
	var pageSize int

	if offset, err = strconv.Atoi(c.Query("offset")); err != nil {
		offset = 0
	}
	if pageSize, err = strconv.Atoi(c.Query("pagesize")); err != nil {
		pageSize = 10
	}

	userFavourites := h.favouriteRepository.GetFavourites(userId, offset, pageSize)

	// Get favourite items for current user and extract these into 3 slices of IDs
	// IDs are then used to query the various repositories
	// Then I wait for all the queries to return and stitch the results back together
	audienceIDs, insightIDs, chartIDs := h.splitUserFavourites(userFavourites)

	var wg sync.WaitGroup
	wg.Add(3)

	// AUDIENCES
	var audiences map[uuid.UUID]*model.Audience

	go func() {
		defer wg.Done()
		audiences = h.audienceRepository.GetAudiences(audienceIDs)
	}()

	// CHARTS
	var charts map[uuid.UUID]*model.Chart

	go func() {
		defer wg.Done()
		charts = h.chartRepository.GetCharts(chartIDs)
	}()

	// INSIGHTS
	var insights map[uuid.UUID]*model.Insight

	go func() {
		defer wg.Done()
		insights = h.insightRepository.GetInsights(insightIDs)
	}()

	wg.Wait()

	// Stitch the responses from the various data sources back together in the
	// same order as the page of favourite ids was retrieved from the DB
	result := h.stitchResults(userFavourites, audiences, charts, insights)

	c.IndentedJSON(http.StatusOK, result)
}

func (*FavouritesHandler) stitchResults(userFavourites favourite.FavouriteEntries,
	audiences map[uuid.UUID]*model.Audience,
	charts map[uuid.UUID]*model.Chart,
	insights map[uuid.UUID]*model.Insight) []model.UserFavourite {
	result := make([]model.UserFavourite, 0, len(userFavourites))
	for _, f := range userFavourites {
		switch f.ResourceType {
		case "audience":
			{
				result = append(result, model.UserFavourite{Description: f.Description, Audience: audiences[f.FavouriteId]})
			}
		case "chart":
			{
				result = append(result, model.UserFavourite{Description: f.Description, Chart: charts[f.FavouriteId]})
			}
		case "insight":
			{
				result = append(result, model.UserFavourite{Description: f.Description, Insight: insights[f.FavouriteId]})
			}
		}
	}
	return result
}

// Splits the list of user favourite entries by resource type and returns 3 slices of ids
func (*FavouritesHandler) splitUserFavourites(userFavourites favourite.FavouriteEntries) ([]uuid.UUID, []uuid.UUID, []uuid.UUID) {
	audienceIDs := make([]uuid.UUID, 0, len(userFavourites))
	insightIDs := make([]uuid.UUID, 0, len(userFavourites))
	chartIDs := make([]uuid.UUID, 0, len(userFavourites))

	for _, f := range userFavourites {
		switch f.ResourceType {
		case "audience":
			{
				audienceIDs = append(audienceIDs, f.FavouriteId)
			}
		case "chart":
			{
				chartIDs = append(chartIDs, f.FavouriteId)
			}
		case "insight":
			{
				insightIDs = append(insightIDs, f.FavouriteId)
			}
		default:
			{
				log.Errorf("unknown userfavourite resource type %v", f.ResourceType)
			}
		}
	}
	return audienceIDs, insightIDs, chartIDs
}

// CreateUserFavourite godoc
//
// @Summary     Create user favourites
// @Description creates a favourite for a specified user
// @Tags        favourites
// @Param       data body model.UserFavouriteShort true "UserFavouriteShort to create"
// @Produce     json
// @Success     201	{object}	model.UserFavourite
// @Failure     400
// @Failure     401
// @Failure     404
// @Failure     500
// @Router      /favourites [post]
func (h *FavouritesHandler) CreateUserFavourite(c *gin.Context) {
	userId, err := auth.GetUserIDFromToken(c)
	if err != nil {
		log.Infof("Error obtaining userid from jwt token: %v\n", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var newFavourite model.UserFavouriteShort
	if err := c.BindJSON(&newFavourite); err != nil {
		log.Infof("Cannot bind request body to UserFavouriteShort object: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	switch newFavourite.ResourceType {
	case "audience", "chart", "insight":
		log.Infof("Adding favourite with id %v description %v for user %v", newFavourite.Id, newFavourite.Description, userId)
		h.favouriteRepository.AddFavourite(userId, newFavourite.Description, newFavourite.Id, newFavourite.ResourceType)
		c.Status(http.StatusCreated)
		body, err := json.Marshal(newFavourite)
		if err != nil {
			log.Errorf("error marshalling favourite to json: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		c.Writer.Write(body)
	default:
		log.Infof("Invalid resource type: %v", newFavourite.ResourceType)
		c.Status(http.StatusBadRequest)
	}
}

// UpdateUserFavourite godoc
//
// @Summary     Update user favourites
// @Description updates a favourite for a specified user
// @Tags        favourites
// @Param       data body model.UserFavouriteShort true "UserFavouriteShort to update"
// @Produce     json
// @Success     204
// @Failure     400
// @Failure     401
// @Failure     404
// @Failure     500
// @Router      /favourites/:favouriteid [patch]
func (h *FavouritesHandler) UpdateUserFavourite(c *gin.Context) {

	userId, err := auth.GetUserIDFromToken(c)
	if err != nil {
		log.Infof("Error obtaining userid from jwt token: %v\n", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	favouriteId, err := uuid.Parse(c.Param("favouriteid"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var updatedFavourite model.UserFavouriteShort
	if err := c.BindJSON(&updatedFavourite); err != nil {
		log.Infof("Cannot bind request body to UserFavouriteShort object: %v", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// ask repo to update this
	log.Infof("UPDATING favourite with uuid %v for user %v", favouriteId, userId)
	h.favouriteRepository.UpdateFavourite(userId,
		updatedFavourite.Description,
		favouriteId,
		updatedFavourite.ResourceType)

	// return 204 no content if ok
	// normally would return updated object
	c.Status(http.StatusNoContent)

}

// DeleteUserFavourite godoc
//
// @Summary     Delete user favourite
// @Description Deletes a favourite for a specified user
// @Tags        favourites
// @Param       favouriteid path uuid true "Favourite ID"
// @Produce     json
// @Success     204
// @Failure     400
// @Failure     401
// @Failure     404
// @Failure     500
// @Router      /favourites/:favouriteid [delete]
func (h *FavouritesHandler) DeleteUserFavourite(c *gin.Context) {

	userId, err := auth.GetUserIDFromToken(c)
	if err != nil {
		log.Infof("Error obtaining userid from jwt token: %v\n", err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	favouriteId, err := uuid.Parse(c.Param("favouriteid"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// ask repo to delete this
	log.Infof("DELETING favourite with uuid %v for user %v", favouriteId, userId)
	h.favouriteRepository.DeleteFavourite(userId, favouriteId)

	// return 204 no content if ok
	c.Status(http.StatusNoContent)

}
