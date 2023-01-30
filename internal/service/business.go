package service

import (
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"
)

type BusinessService interface {
	CreateBusiness(business *dto.BriefBusiness) error

	GetBusiness(email string) (*dto.BriefBusiness, error)
	DeleteBusiness(email string, userId string) error
}

type businessService struct {
	dbHandler repository.DBHandler
}

func NewBusinessService(dbHandler repository.DBHandler) BusinessService {
	return &businessService{dbHandler}
}

func (b *businessService) CreateBusiness(
	business *dto.BriefBusiness) error {
	err := b.dbHandler.NewBusinessQuery().CreateBusiness(business)
	if err != nil {
		return err
	}
	return nil
}

func (b *businessService) GetBusiness(email string) (*dto.BriefBusiness, error) {
	business, err := b.dbHandler.NewBusinessQuery().GetBusiness(email)
	if err != nil {
		return nil, err
	}
	return business, nil
}

func (b *businessService) DeleteBusiness(email string, userId string) error {
	err := b.dbHandler.NewBusinessQuery().DeleteBusiness(email, userId)
	if err != nil {
		return err
	}
	return nil
}
