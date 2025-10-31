package middleware

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorsResponse struct {
	Errors []ValidationError `json:"errors"`
}
