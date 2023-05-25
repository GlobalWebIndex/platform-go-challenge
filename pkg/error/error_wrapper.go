package app_error

const (
	MessageUnknownError = "An unknown error has occured"
)

type AppError struct {
	Err                 error
	UserFriendlyMessage string
}

func (e AppError) Error() string {
	return e.UserFriendlyMessage
}
