package domain

import "errors"

var (
	ErrWrongAssetInput = errors.New("wrong input for asset")
	ErrWrongQueryInput = errors.New("wrong input for query")
	ErrWrongUserInput  = errors.New("wrong input for user")
	ErrWrongLoginInput = errors.New("wrong input for login")

	ErrInternalDBFailure = errors.New("internal failure with the DB")
)
