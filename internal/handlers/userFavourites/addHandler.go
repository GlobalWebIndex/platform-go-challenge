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

type AddUserFavouriteHandler struct {
	UserFavouriteService services.UserFavouriteServiceInterface
	AssetTypeValidator   helpers.AssetTypeValidatorInterface
	logger               logger.LoggerInterface
}

func NewAddUserFavouriteHandler(
	service services.UserFavouriteServiceInterface,
	AssetTypeValidator helpers.AssetTypeValidatorInterface,
	logger logger.LoggerInterface,
) *AddUserFavouriteHandler {
	return &AddUserFavouriteHandler{
		UserFavouriteService: service,
		AssetTypeValidator:   AssetTypeValidator,
		logger:               logger,
	}
}

// swagger:operation POST /user/{userId}/favourites addUserFavouriteAsset
//
// # It adds Asset with ID {assetId} to User's with ID {userId} favourites
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
//		description: Asset added to favourites
//		schema:
//			$ref: "#/definitions/AddUserFavouriteAssetResponse"
//	"400":
//		description: Bad request - request parameters are missing or invalid
//		schema:
//			$ref: "#/definitions/AddUserFavouriteAssetResponse"
//	"404":
//		description: Requested user or asset not found
//		schema:
//			$ref: "#/definitions/AddUserFavouriteAssetResponse"
//	"500":
//		description: Internal server error - check logs for details
//		schema:
//			$ref: "#/definitions/AddUserFavouriteAssetResponse"
func (handler *AddUserFavouriteHandler) AddUserFavouriteAssetController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	var err error
	response := &AddUserFavouriteAssetResponse{}
	userFavouriteRequest := &AddUserFavouriteAssetRequest{}

	userIdAsString := mux.Vars(r)["user_id"]
	if userIdAsString == "" {
		response.ErrorMessage = "missing user id"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	userId, err := strconv.Atoi(userIdAsString)
	if err != nil {
		handler.logger.Error("Error in adding user favourite assets",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})

		response.ErrorMessage = "malformed user id"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	err = json.NewDecoder(r.Body).Decode(userFavouriteRequest)
	if err != nil {
		handler.logger.Error("Error in adding user favourite assets",
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

	err = handler.UserFavouriteService.AddAsset(
		ctx,
		uint(userId),
		userFavouriteRequest.AssetId,
		userFavouriteRequest.AssetType,
	)
	if userNotFoundError, ok := err.(apierrors.UserNotFoundErrorWrapper); ok {
		handler.logger.Error("Error in adding user favourite assets",
			map[string]interface{}{
				"errorMessage": userNotFoundError.Unwrap(),
			})

		response.ErrorMessage = err.Error()
		handler.JsonResponse(w, userNotFoundError.ReturnedStatusCode, response)

		return
	}

	if err != nil {
		handler.logger.Error("Error in adding user favourite assets",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})

		response.ErrorMessage = "error in adding user favourite assets"
		handler.JsonResponse(w, http.StatusInternalServerError, response)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *AddUserFavouriteHandler) JsonResponse(
	w http.ResponseWriter,
	statusCode int,
	response *AddUserFavouriteAssetResponse,
) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.ErrorMessage = "error in adding user favourite asset - json response"

		handler.logger.Error("Error in adding user favourite asset - json response",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})
	}
}
