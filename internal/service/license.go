package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"ownify_api/internal/repository"
	"ownify_api/internal/utils"
)

type LicenseService interface {
	GenerateAPIKey(email string, userId string) (*string, error)
	GetAPIKey(email string, userId string) ([]string, error)
}

type licenseService struct {
	dao repository.DBHandler
}

func (s *licenseService) GenerateAPIKey(email string, userId string) (*string, error) {
	randomPart, err := generateRandomString(12) // Generate a 10-character random string
	if err != nil {
		return nil, errors.New("could not generate API key: " + err.Error())
	}

	encryptedEmail, err := utils.Encrypt(email, randomPart)
	if err != nil {
		return nil, err
	}

	encryptedUserId, err := utils.Encrypt(userId, randomPart)
	if err != nil {
		return nil, err
	}
	// Concatenate email, userId, and the random part to create the API key
	apiKey := fmt.Sprintf("%s-%s-%s-%s", "ownify", encryptedEmail, encryptedUserId, randomPart)
	err = s.dao.NewLicenseQuery().SaveAPIKey(email, userId, apiKey)

	if err != nil {
		return nil, errors.New("could not save API key: " + err.Error())
	}

	return &apiKey, nil
}

func (s *licenseService) GetAPIKey(email string, userId string) ([]string, error) {
	apiKeys, err := s.dao.NewLicenseQuery().GetAPIKey(email, userId)
	if err != nil {
		return nil, errors.New("You did not create API Key: " + err.Error())
	}
	return apiKeys, nil
}

func NewLicenseService(dao repository.DBHandler) LicenseService {
	return &licenseService{dao: dao}
}

// generateRandomString generates a random string of the given length.
func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
