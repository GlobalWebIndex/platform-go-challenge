package customerrors

type ValidationError struct {
	Message string `json:"message"`
	ErrorCode string `json:"error_code"`
}

func (e *ValidationError) Error() string {
	return e.Message
}