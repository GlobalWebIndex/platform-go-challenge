package persistence

import (
	"context"
	"errors"
	"log"
	"platform-go-challenge/consts"
	"platform-go-challenge/domain"

	"github.com/arangodb/go-driver"
)

type UserRepository struct {
	collection driver.Collection
}

func NewUserRepository(collection driver.Collection) *UserRepository {
	return &UserRepository{collection: collection}
}

func (ur UserRepository) CreateUsers() ([]domain.User, error) {
	numOfDocs, err := ur.collection.Count(context.Background())
	if err != nil {
		log.Println("Failed to get document count for user collection")
	}
	if numOfDocs > 0 {
		return []domain.User{}, errors.New("there are already documents in user collections")
	}

	users := consts.Users

	arangoUsers := []domain.User{}
	for _, user := range users {
		meta, err := ur.collection.CreateDocument(context.Background(), &user)
		if err != nil {
			log.Printf("Couldn't create new document")
			return []domain.User{}, err
		}
		user.Key = meta.Key
		arangoUsers = append(arangoUsers, user)
	}
	return arangoUsers, nil
}

func (ur UserRepository) GetUser(userID string) (string, error) {
	user, err := ur.collection.ReadDocument(context.Background(), userID, nil)
	if err != nil {
		return "", err
	}
	return string(user.ID), nil
}
