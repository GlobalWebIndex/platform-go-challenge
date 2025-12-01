package storage

import (
	"github.com/dintzeler/platform-go-challenge/models"
	"encoding/json"
	"os"
	"github.com/dintzeler/platform-go-challenge/customerrors"
)

type DataStore struct {
	Charts []models.Chart `json:"charts"`
	Insights []models.Insight `json:"insights"`
	Audiences []models.Audience `json:"audiences"`
	Users []models.User `json:"users"`
	Favorites []models.Favorite `json:"favorites"`
}

func LoadData(fileName string) (*DataStore, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var dataStore DataStore
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dataStore)
	if err != nil {
		return nil, err
	}
	return &dataStore, nil
}

func SaveData(fileName string, dataStore *DataStore) error {
    file, err := os.Create(fileName)
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ") // Optional: for pretty formatting
    err = encoder.Encode(dataStore)
    if err != nil {
        return &customerrors.ValidationError{
			Message:   "Error saving data",
			ErrorCode: "DATA_SAVE_ERROR",
		}
    }
    return nil
}