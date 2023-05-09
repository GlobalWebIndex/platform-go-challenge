package userFavourite

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/loukaspe/platform-go-challenge/internal/core/services"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"github.com/loukaspe/platform-go-challenge/pkg/logger"
	"net/http"
	"strconv"
)

type GetUserFavouriteHandler struct {
	UserFavouriteService services.UserFavouriteServiceInterface
	logger               logger.LoggerInterface
}

func NewGetUserFavouriteHandler(
	service services.UserFavouriteServiceInterface,
	logger logger.LoggerInterface,
) *GetUserFavouriteHandler {
	return &GetUserFavouriteHandler{
		UserFavouriteService: service,
		logger:               logger,
	}
}

// swagger:operation GET /user/{userId}/favourites getUserFavouriteAsset
//
// # It returns favourite Assets of User's with ID {userId}
//
// ---
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes:
//	- http
//	- https
//
// responses:
//
//	"200":
//		description: Assets retrieved
//		schema:
//			$ref: "#/definitions/GetUserFavouriteAssetResponse"
//	"204":
//		description: Successful call but user has no favourites
//		schema:
//			$ref: "#/definitions/GetUserFavouriteAssetResponse"
//	"400":
//		description: Bad request - request parameters are missing or invalid
//		schema:
//			$ref: "#/definitions/GetUserFavouriteAssetResponse"
//	"404":
//		description: Requested user not found
//		schema:
//			$ref: "#/definitions/GetUserFavouriteAssetResponse"
//	"500":
//		description: Internal server error - check logs for details
//		schema:
//			$ref: "#/definitions/GetUserFavouriteAssetResponse"
func (handler *GetUserFavouriteHandler) GetUserFavouriteController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var err error
	response := &GetUserFavouriteAssetsResponse{}

	ctx := r.Context()

	userIdAsString := mux.Vars(r)["user_id"]
	if userIdAsString == "" {
		response.ErrorMessage = "missing user id"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	userId, err := strconv.Atoi(userIdAsString)
	if err != nil {
		handler.logger.Error("Error in getting user favourite assets",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})

		response.ErrorMessage = "malformed user id"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	userFavourites, err := handler.UserFavouriteService.GetAssets(
		ctx,
		uint(userId),
	)
	if userNotFoundError, ok := err.(apierrors.UserNotFoundErrorWrapper); ok {
		handler.logger.Error("Error in getting users favourite assets",
			map[string]interface{}{
				"errorMessage": userNotFoundError.Unwrap(),
			})

		response.ErrorMessage = err.Error()
		handler.JsonResponse(w, userNotFoundError.ReturnedStatusCode, response)

		return
	}

	if userHasNoFavouritesError, ok := err.(apierrors.NoFavouriteAssetsErrorWrapper); ok {
		handler.JsonResponse(w, userHasNoFavouritesError.ReturnedStatusCode, response)

		return
	}

	if err != nil {
		handler.logger.Error("Error in getting user's favourite assets",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})

		response.ErrorMessage = "error in getting user favourite assets"
		handler.JsonResponse(w, http.StatusInternalServerError, response)

		return
	}

	response.FavouritesAssets = userFavourites
	handler.JsonResponse(w, http.StatusOK, response)
}

func (handler *GetUserFavouriteHandler) JsonResponse(
	w http.ResponseWriter,
	statusCode int,
	response *GetUserFavouriteAssetsResponse,
) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.ErrorMessage = "error in getting user favourite asset - json response"

		handler.logger.Error("Error in getting user favourite asset - json response",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})
	}
}
