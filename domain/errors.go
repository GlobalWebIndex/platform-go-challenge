package domain

import "errors"

var (
	ErrWrongAssetInput = errors.New("wrong input for asset")
	ErrWrongQueryInput = errors.New("wrong input for query")
	ErrWrongUserInput  = errors.New("wrong input for user")
	ErrWrongLoginInput = errors.New("wrong input for login")
	ErrUserNotFound    = errors.New("user not found")

	ErrUnauthorized = errors.New("unauthorized")

	ErrInternalDBFailure = errors.New("internal failure with the DB")
)
