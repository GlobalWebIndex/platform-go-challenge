package dto

import (
	"errors"
	"fmt"
)

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

type UserDto struct {
	ID       string
	Gender   Gender
	Email    string
	Password string
	Country  string
	Age      int
}

func (u UserDto) String() string {
	return fmt.Sprintf("ID: %s, Gender: %s, Email: %s, Country: %s, Age: %d", u.ID, u.Gender, u.Email, u.Country, u.Age)
}

// List of valid countries, you might want to update this with actual values
var validCountries = []string{
	"United States",
	"Canada",
	// Add all countries you consider valid...
}

func CreateUserDto(id string, gender Gender, email string, password string, country string, age int) (*UserDto, error) {
	// Validate the gender
	if gender != Male && gender != Female {
		return nil, errors.New("invalid gender value")
	}

	// Validate the age
	if age < 16 || age > 100 {
		return nil, errors.New("invalid age value")
	}

	// Validate the country
	if !isValidCountry(country) {
		return nil, errors.New("invalid country value")
	}

	// Create the UserDto
	return &UserDto{
		ID:       id,
		Gender:   gender,
		Email:    email,
		Password: password,
		Country:  country,
		Age:      age,
	}, nil
}

func isValidCountry(country string) bool {
	for _, validCountry := range validCountries {
		if country == validCountry {
			return true
		}
	}

	return false
}
