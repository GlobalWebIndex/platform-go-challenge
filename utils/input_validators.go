package utils

import (
	favorite_model "gwi/models/favorite"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateFavorite validates the given Favorite struct.
// It returns an error if the validation fails, otherwise it returns nil.
func ValidateFavorite(fav *favorite_model.Favorite) error {
	return validate.Struct(fav)
}
