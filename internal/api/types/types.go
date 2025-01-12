package types

const (
	// Generic error codes
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrTimedOut         = "TIMED_OUT"

	// Business logic error codes
	ErrCodeAssetNotFound    = "ASSET_NOT_FOUND"
	ErrCodeAlreadyFavorited = "ALREADY_FAVORITED"
	ErrCodeNotFavorited     = "NOT_FAVORITED"
	ErrCodeInvalidAssetType = "INVALID_ASSET_TYPE"

	// Authentication error codes
	ErrCodeInvalidToken = "INVALID_TOKEN"

	// Ratelimit error codes
	ErrCodeIPRateLimit = "IP_RATE_LIMIT"
)

// APIResponse is our standardized response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *MetaData   `json:"meta,omitempty"`
}

// APIError represents an error in our API
type APIError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	Field      string `json:"field,omitempty"`
	HTTPStatus int    `json:"-"`
}

// MetaData contains additional response information
type MetaData struct {
	Page      int    `json:"page,omitempty"`
	PerPage   int    `json:"per_page,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// Response constructors
func NewSuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

func NewErrorResponse(code, message string) APIResponse {
	return APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	}
}

// WithMeta adds metadata to a response
func (r APIResponse) WithMeta(meta *MetaData) APIResponse {
	r.Meta = meta
	return r
}
