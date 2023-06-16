// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: proto/core/favourite/v1/favourite.proto

package favourite_pb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on FavouriteAsset with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *FavouriteAsset) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on FavouriteAsset with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in FavouriteAssetMultiError,
// or nil if none found.
func (m *FavouriteAsset) ValidateAll() error {
	return m.validate(true)
}

func (m *FavouriteAsset) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetId()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, FavouriteAssetValidationError{
					field:  "Id",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, FavouriteAssetValidationError{
					field:  "Id",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetId()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FavouriteAssetValidationError{
				field:  "Id",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetIdUser()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, FavouriteAssetValidationError{
					field:  "IdUser",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, FavouriteAssetValidationError{
					field:  "IdUser",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetIdUser()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FavouriteAssetValidationError{
				field:  "IdUser",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetIdAsset()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, FavouriteAssetValidationError{
					field:  "IdAsset",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, FavouriteAssetValidationError{
					field:  "IdAsset",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetIdAsset()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return FavouriteAssetValidationError{
				field:  "IdAsset",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if m.Md != nil {

		if all {
			switch v := interface{}(m.GetMd()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, FavouriteAssetValidationError{
						field:  "Md",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, FavouriteAssetValidationError{
						field:  "Md",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetMd()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return FavouriteAssetValidationError{
					field:  "Md",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if m.Asset != nil {

		if all {
			switch v := interface{}(m.GetAsset()).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, FavouriteAssetValidationError{
						field:  "Asset",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, FavouriteAssetValidationError{
						field:  "Asset",
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(m.GetAsset()).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return FavouriteAssetValidationError{
					field:  "Asset",
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return FavouriteAssetMultiError(errors)
	}

	return nil
}

// FavouriteAssetMultiError is an error wrapping multiple validation errors
// returned by FavouriteAsset.ValidateAll() if the designated constraints
// aren't met.
type FavouriteAssetMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m FavouriteAssetMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m FavouriteAssetMultiError) AllErrors() []error { return m }

// FavouriteAssetValidationError is the validation error returned by
// FavouriteAsset.Validate if the designated constraints aren't met.
type FavouriteAssetValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e FavouriteAssetValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e FavouriteAssetValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e FavouriteAssetValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e FavouriteAssetValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e FavouriteAssetValidationError) ErrorName() string { return "FavouriteAssetValidationError" }

// Error satisfies the builtin error interface
func (e FavouriteAssetValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sFavouriteAsset.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = FavouriteAssetValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = FavouriteAssetValidationError{}
