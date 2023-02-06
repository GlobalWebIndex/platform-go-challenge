package app

import (
	"context"
	"fmt"

	"ownify_api/internal/utils"
	desc "ownify_api/pkg"
)

func (m *MicroserviceServer) GrantBusiness(ctx context.Context, req *desc.BusinessGrantRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	user_id, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}
	if utils.IsEmail(req.MyEmail) != nil || utils.IsEmail(req.BusinessEmail) != nil {
		return nil, fmt.Errorf("[ERR] include invalid email address")
	}

	isRegistered := m.authService.ValidAdmin(*user_id, req.MyEmail)
	if !isRegistered {
		return nil, fmt.Errorf("[ERR] You are not admin", req.MyEmail)
	}

	err = m.adminService.GrantBusiness(req.BusinessEmail, req.IsApproved)
	if err != nil {
		return nil, err
	}

	msg, err := utils.Encrypt("successfully granted", req.BusinessEmail)
	if err != nil {
		return nil, err
	}
	err = m.notifyService.SendMessage(msg)
	
	if err != nil {
		return nil, err
	}
	return &desc.NetWorkResponse{
		Msg:     "Successfully Granted!",
		Success: true,
	}, nil
}
