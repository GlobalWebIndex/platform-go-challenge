package handlers

import (
	"encoding/json"
	"github.com/loukaspe/platform-go-challenge/internal/core/domain"
	"github.com/loukaspe/platform-go-challenge/internal/core/services"
	"github.com/loukaspe/platform-go-challenge/pkg/logger"
	"net/http"
)

type JwtClaimsHandler struct {
	service *services.JwtService
	logger  logger.LoggerInterface
}
type JwtClaimsHandlerInterface interface {
	UrlAddController(w http.ResponseWriter, r *http.Request)
}

func NewJwtClaimsHandler(
	service *services.JwtService,
	logger logger.LoggerInterface,
) *JwtClaimsHandler {
	return &JwtClaimsHandler{
		service: service,
		logger:  logger,
	}
}

// request for generating jwt token
//
// swagger:parameters jwtToken
type JwtRequest struct {
	// in:body
	// Required: true
	Username string `json:"username"`
	// in:body
	// Required: true
	Password string `json:"password"`
}

// Response with jwtToken
// swagger:model JwtResponse
type JwtResponse struct {
	// jwt token
	//
	// Required: false
	Token string `json:"token"`
	// possible error message
	//
	// Required: false
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// swagger:operation POST /token jwtToken
//
// # Generates JWT token for authentication and authorization
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
//	responses:
//		"200":
//			description: OK
//			schema:
//				$ref: "#/definitions/JwtResponse"
//		"500":
//			description: Error
//			schema:
//				$ref: "#/definitions/JwtResponse"
func (handler *JwtClaimsHandler) JwtTokenController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request JwtRequest
	var response JwtResponse

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		handler.logger.Error("Error in creating jwt token - request decode",
			map[string]interface{}{
				"errorMessage": err.Error(),
			})

		response.ErrorMessage = "malformed auth request"

		handler.JsonResponse(w, http.StatusInternalServerError, &response)

		return
	}

	if request.Username == "" || request.Password == "" {
		response.ErrorMessage = "empty username or password"

		handler.JsonResponse(w, http.StatusBadRequest, &response)

		return
	}

	// This is where I would implement the login, but I did not have time

	domainUser := domain.User{
		Username: request.Username,
		Password: request.Password,
	}

	result, err := handler.service.CreateJwtTokenService(domainUser)
	if err != nil {
		response.ErrorMessage = "error during creation of the token"

		handler.JsonResponse(w, http.StatusInternalServerError, &response)

		return
	}

	response.Token = result
	handler.JsonResponse(w, http.StatusOK, &response)

	return
}

func (handler *JwtClaimsHandler) JsonResponse(
	w http.ResponseWriter,
	statusCode int,
	response *JwtResponse,
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
