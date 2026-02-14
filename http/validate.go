package http

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationError returns a single user-facing message from validator errors.
// If err is nil or not a validator.ValidationErrors, it returns "".
func ValidationError(err error) string {
	if err == nil {
		return ""
	}
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return err.Error()
	}
	msgs := make([]string, 0, len(errs))
	for _, e := range errs {
		if e.Field() != "" {
			msgs = append(msgs, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		} else {
			msgs = append(msgs, e.Tag())
		}
	}
	return strings.Join(msgs, "; ")
}

// validateVar validates a single value with the given tag and returns a message on failure.
// Note: "required" rejects empty string; for optional string fields use other tags or skip validation when the value is empty.
func validateVar(field interface{}, tag string) string {
	if err := validate.Var(field, tag); err != nil {
		return ValidationError(err)
	}
	return ""
}
