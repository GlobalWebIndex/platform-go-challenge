package apierrors

type UserNotFoundErrorWrapper struct {
	ReturnedStatusCode int
	OriginalError      error
}

// Error the original error message remains as it is for logging reasons etc.
// and the wrapper error message is empty because we don't want the client to see anything
func (err UserNotFoundErrorWrapper) Error() string {
	return ""
}

func (err UserNotFoundErrorWrapper) Unwrap() error {
	return err.OriginalError
}

type AssetNotFoundErrorWrapper struct {
	ReturnedStatusCode int
	OriginalError      error
}

// Error the original error message remains as it is for logging reasons etc.
// and the wrapper error message is empty because we don't want the client to see anything
func (err AssetNotFoundErrorWrapper) Error() string {
	return ""
}

func (err AssetNotFoundErrorWrapper) Unwrap() error {
	return err.OriginalError
}

type UnknownAssetTypeErrorWrapper struct {
	ReturnedStatusCode int
	OriginalError      error
}

// Error the original error message remains as it is for logging reasons etc.
// and the wrapper error message is empty because we don't want the client to see anything
func (err UnknownAssetTypeErrorWrapper) Error() string {
	return ""
}

func (err UnknownAssetTypeErrorWrapper) Unwrap() error {
	return err.OriginalError
}

type NoFavouriteAssetsErrorWrapper struct {
	ReturnedStatusCode int
	OriginalError      error
}

// Error the original error message remains as it is for logging reasons etc.
// and the wrapper error message is empty because we don't want the client to see anything
func (err NoFavouriteAssetsErrorWrapper) Error() string {
	return ""
}

func (err NoFavouriteAssetsErrorWrapper) Unwrap() error {
	return err.OriginalError
}
