package userFavourite

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/loukaspe/platform-go-challenge/internal/core/services"
	apierrors "github.com/loukaspe/platform-go-challenge/pkg/errors"
	"github.com/loukaspe/platform-go-challenge/pkg/helpers"
	"github.com/loukaspe/platform-go-challenge/pkg/logger"
	"net/http"
	"strconv"
)

type DeleteUserFavouriteHandler struct {
	UserFavouriteService services.UserFavouriteServiceInterface
	AssetTypeValidator   helpers.AssetTypeValidatorInterface
	logger               logger.LoggerInterface
}

func NewDeleteUserFavouriteHandler(
	service services.UserFavouriteServiceInterface,
	AssetTypeValidator helpers.AssetTypeValidatorInterface,
	logger logger.LoggerInterface,
) *DeleteUserFavouriteHandler {
	return &DeleteUserFavouriteHandler{
		UserFavouriteService: service,
		AssetTypeValidator:   AssetTypeValidator,
		logger:               logger,
	}
}

// swagger:operation DELETE /user/{userId}/favourites deleteUserFavouriteAsset
//
// # It deletes favourite Asset with ID {assetId} of User's with ID {userId}
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
//		description: Asset deleted
//		schema:
//			$ref: "#/definitions/DeleteUserFavouriteAssetResponse"
//	"400":
//		description: Bad request - request parameters are missing or invalid
//		schema:
//			$ref: "#/definitions/DeleteUserFavouriteAssetResponse"
//	"404":
//		description: Requested user or asset not found
//		schema:
//			$ref: "#/definitions/DeleteUserFavouriteAssetResponse"
//	"500":
//		description: Internal server error - check logs for details
//		schema:
//			$ref: "#/definitions/DeleteUserFavouriteAssetResponse"
func (handler *DeleteUserFavouriteHandler) DeleteUserFavouriteController(w http.ResponseWriter, r *http.Request) {
	var err error
	response := &DeleteUserFavouriteAssetResponse{}
	userFavouriteRequest := &DeleteUserFavouriteAssetRequest{}

	ctx := r.Context()

	userIdAsString := mux.Vars(r)["user_id"]
	if userIdAsString == "" {
		response.ErrorMessage = "missing user id"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	userId, err := strconv.Atoi(userIdAsString)
	if err != nil {
		handler.logger.Error("Error in deleting user favourite assets",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})

		response.ErrorMessage = "malformed user id"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	err = json.NewDecoder(r.Body).Decode(userFavouriteRequest)
	if err != nil {
		handler.logger.Error("Error in deleting user favourite assets",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})

		response.ErrorMessage = "malformed user favourite assets request"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	err = handler.AssetTypeValidator.ValidateAssetType(userFavouriteRequest.AssetType)
	if err != nil {
		response.ErrorMessage = err.Error()
		handler.JsonResponse(w, http.StatusBadRequest, response)
		return
	}

	err = handler.UserFavouriteService.RemoveAsset(
		ctx,
		uint(userId),
		userFavouriteRequest.AssetId,
		userFavouriteRequest.AssetType,
	)
	if userNotFoundError, ok := err.(apierrors.UserNotFoundErrorWrapper); ok {
		handler.logger.Error("Error in deleting user favourite assets",
			map[string]interface{}{
				"errorMessage": userNotFoundError.Unwrap(),
			})

		response.ErrorMessage = err.Error()
		handler.JsonResponse(w, userNotFoundError.ReturnedStatusCode, response)

		return
	}

	if err != nil {
		handler.logger.Error("Error in deleting user favourite assets",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})

		response.ErrorMessage = "error in deleting user favourite assets"
		handler.JsonResponse(w, http.StatusInternalServerError, response)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *DeleteUserFavouriteHandler) JsonResponse(
	w http.ResponseWriter,
	statusCode int,
	response *DeleteUserFavouriteAssetResponse,
) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.ErrorMessage = "error in deleting user favourite asset - json response"

		handler.logger.Error("Error in deleting user favourite asset - json response",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})
	}
}
