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

type UpdateUserFavouriteHandler struct {
	UserFavouriteService services.UserFavouriteServiceInterface
	AssetTypeValidator   helpers.AssetTypeValidatorInterface
	logger               logger.LoggerInterface
}

func NewUpdateUserFavouriteHandler(
	service services.UserFavouriteServiceInterface,
	AssetTypeValidator helpers.AssetTypeValidatorInterface,
	logger logger.LoggerInterface,
) *UpdateUserFavouriteHandler {
	return &UpdateUserFavouriteHandler{
		UserFavouriteService: service,
		AssetTypeValidator:   AssetTypeValidator,
		logger:               logger,
	}
}

// swagger:operation PATCH /user/{userId}/favourites/{assetId} updateUserFavouriteAssetDescription
//
// # It updates Asset's with ID {assetId} description of User's with ID {userId} favourites
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
//		description: Asset's description updated
//		schema:
//			$ref: "#/definitions/UpdateUserFavouriteAssetDescriptionResponse"
//	"400":
//		description: Bad request - request parameters are missing or invalid
//		schema:
//			$ref: "#/definitions/UpdateUserFavouriteAssetDescriptionResponse"
//	"404":
//		description: Requested user or asset not found
//		schema:
//			$ref: "#/definitions/UpdateUserFavouriteAssetDescriptionResponse"
//	"500":
//		description: Internal server error - check logs for details
//		schema:
//			$ref: "#/definitions/UpdateUserFavouriteAssetDescriptionResponse"
func (handler *UpdateUserFavouriteHandler) UpdateUserFavouriteController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var err error
	response := &UpdateUserFavouriteAssetDescriptionResponse{}
	userFavouriteRequest := &UpdateUserFavouriteAssetDescriptionRequest{}

	ctx := r.Context()

	userIdAsString := mux.Vars(r)["user_id"]
	if userIdAsString == "" {
		response.ErrorMessage = "missing user id"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	userId, err := strconv.Atoi(userIdAsString)
	if err != nil {
		handler.logger.Error("Error in updating user favourite asset description",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})

		response.ErrorMessage = "malformed user id"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	assetIdAsString := mux.Vars(r)["asset_id"]
	if assetIdAsString == "" {
		response.ErrorMessage = "missing asset id"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	assetId, err := strconv.Atoi(assetIdAsString)
	if err != nil {
		handler.logger.Error("Error in updating user favourite assets description",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})

		response.ErrorMessage = "malformed asset id"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	err = json.NewDecoder(r.Body).Decode(userFavouriteRequest)
	if err != nil {
		handler.logger.Error("Error in updating user favourite asset description",
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

	if userFavouriteRequest.Description == "" {
		response.ErrorMessage = "empty description"

		handler.JsonResponse(w, http.StatusBadRequest, response)

		return
	}

	err = handler.UserFavouriteService.EditAssetDescription(
		ctx,
		uint(userId),
		uint(assetId),
		userFavouriteRequest.AssetType,
		userFavouriteRequest.Description,
	)
	if userNotFoundError, ok := err.(apierrors.UserNotFoundErrorWrapper); ok {
		handler.logger.Error("Error in updating user favourite asset description",
			map[string]interface{}{
				"errorMessage": userNotFoundError.Unwrap(),
			})

		response.ErrorMessage = err.Error()
		handler.JsonResponse(w, userNotFoundError.ReturnedStatusCode, response)

		return
	}

	if err != nil {
		handler.logger.Error("Error in updating user favourite assets description",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})

		response.ErrorMessage = "error in updating user favourite assets description"
		handler.JsonResponse(w, http.StatusInternalServerError, response)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *UpdateUserFavouriteHandler) JsonResponse(
	w http.ResponseWriter,
	statusCode int,
	response *UpdateUserFavouriteAssetDescriptionResponse,
) {
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.ErrorMessage = "error in updating user favourite asset description - json response"

		handler.logger.Error("Error in updating user favourite asset description - json response",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})
	}
}
