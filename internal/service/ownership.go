package service

import (
	"fmt"
	"ownify_api/internal/domain"
	"ownify_api/internal/dto"
	"ownify_api/internal/repository"
)

type OwnershipService interface {
	GetOwnerShip(pubKey string) (*dto.BriefBusiness, *dto.BriefUser, error)
}

type ownershipService struct {
	dbHandler repository.DBHandler
}

func NewOwnershipService(dbHandler repository.DBHandler) OwnershipService {
	return &ownershipService{dbHandler}
}

// GetOwnerShip implements OwnershipService
func (o *ownershipService) GetOwnerShip(pubKey string) (*dto.BriefBusiness, *dto.BriefUser, error) {
	bChan := make(chan domain.Result[dto.BriefBusiness])
	uChan := make(chan domain.Result[dto.BriefUser])
	defer close(bChan)
	defer close(uChan)
	go func() {
		business, err := o.dbHandler.NewBusinessQuery().GetBusinessByWalletAddress(pubKey)
		if err != nil {
			bChan <- domain.Result[dto.BriefBusiness]{
				Err: err,
			}
			return
		}
		bChan <- domain.Result[dto.BriefBusiness]{Val: *business}
	}()

	go func() {
		user, err := o.dbHandler.NewUserQuery().ValidUser(pubKey, "")
		if err != nil {
			uChan <- domain.Result[dto.BriefUser]{
				Err: err,
			}
			return
		}
		uChan <- domain.Result[dto.BriefUser]{Val: *user}
	}()

	b := <-bChan
	u := <-uChan
	if b.Err != nil && u.Err != nil {
		return nil, nil, fmt.Errorf("[Err] this user don't exist in ownify: %s", pubKey)
	}
	return &b.Val, &u.Val, nil
}
