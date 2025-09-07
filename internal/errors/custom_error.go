package errors

// CustomError is a structured error for API responses.
type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// Error implements the error interface.
func (e *CustomError) Error() string {
	return e.Message
}

// NewCustomError creates a new CustomError instance.
func NewCustomError(code int, message, detail string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}
