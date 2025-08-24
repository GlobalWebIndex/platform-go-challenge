package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func DecodeJSON(r *http.Request, v any) error {
	var unmarshalErr *json.UnmarshalTypeError
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		if errors.As(err, &unmarshalErr) {
			return fmt.Errorf("invalid JSON format: %w", err)
		} else {
			return fmt.Errorf("invalid JSON format: %w", err)
		}
	}
	return nil
}
