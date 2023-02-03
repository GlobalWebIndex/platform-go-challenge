package app

import (
	"context"
	"fmt"

	"ownify_api/internal/dto"
	"ownify_api/internal/utils"
	desc "ownify_api/pkg"
)

func (m *MicroserviceServer) CreateBusiness(ctx context.Context, req *desc.CreateBusinessRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return &desc.NetWorkResponse{}, err
	}

	isRegistered := m.authService.CheckEmail(req.Email)
	if isRegistered {
		return nil, fmt.Errorf("Already registered!")
	}

	business := dto.BriefBusiness{
		UserId:      *uid,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		Pin:         utils.Hash(req.Password),
		Business:    req.Business,
		PhoneNumber: req.PhoneNumber,
		Location:    req.Location,
	}
	err = m.businessService.CreateBusiness(&business)
	if err != nil {
		return nil, err
	}
	return &desc.NetWorkResponse{Success: true, Msg: "Successfully created!"}, nil
}

func (m *MicroserviceServer) DeleteBusiness(ctx context.Context, req *desc.DeleteBusinessRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, fmt.Errorf("[ERR] no permission to user: %s", req.Email)
	}

	if !m.authService.ValidBusiness(*uid, req.Email) {
		return nil, fmt.Errorf("[ERR] no permission to user: %s", req.Email)
	}

	err = m.businessService.DeleteBusiness(req.Email, *uid)
	if err != nil {
		return nil, err
	}
	return &desc.NetWorkResponse{Success: true, Msg: "Successfully deleted."}, nil
}

func (m *MicroserviceServer) GetBusiness(ctx context.Context, req *desc.GetBusinessRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return &desc.NetWorkResponse{Success: false, Msg: "Access denied."}, err
	}

	err = utils.IsEmail(req.Email)
	if err != nil {
		return nil, err
	}

	data, err := m.businessService.GetBusiness(req.Email)
	if err != nil {
		return nil, err
	}

	return BuildRes(data, "Here is your business info", true)
}
