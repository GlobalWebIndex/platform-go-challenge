package app

import (
	"context"
	"gwi_api/internal/dto"

	desc "gwi_api/pkg"
)

func (m *MicroserviceServer) SignUp(ctx context.Context, req *desc.SignUpRequest) (*desc.NetWorkResponse, error) {

	user := dto.UserDto{
		Gender:   dto.Gender(req.Gender),
		Country:  req.Country,
		Age:      int(req.Age),
		Email:    req.Email,
		Password: req.Password,
	}

	userId, token, err := m.authService.SignUp(user)
	if err != nil {
		return nil, err
	}

	type SignUpResponse struct {
		UserId int64
		Token  string
	}

	data := SignUpResponse{
		UserId: *userId,
		Token:  *token,
	}
	return BuildRes[SignUpResponse](data, "successfully sign up", true)
}

func (m *MicroserviceServer) SignIn(ctx context.Context, req *desc.SignInRequest) (*desc.NetWorkResponse, error) {

	token, err := m.authService.SignIn(req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	type SignInResponse struct {
		Token string
	}
	data := SignInResponse{
		Token: *token,
	}
	return BuildRes[SignInResponse](data, "successfully login", true)
}
