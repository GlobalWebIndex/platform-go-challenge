package app

import (
	"context"
	"fmt"
	"log"

	"ownify_api/internal/dto"
	desc "ownify_api/pkg"
)

func (m *MicroserviceServer) CreateUser(ctx context.Context, req *desc.CreateUserRequest) (*desc.NetWorkResponse, error) {

	// validate token.
	user_id, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}

	isRegistered := m.authService.ValidUser(req.WalletAddress, req.IdFingerprint)
	if isRegistered {
		return nil, fmt.Errorf("[ERR] Already registered! with %s", req.WalletAddress)
	}

	user := dto.BriefUser{
		PubAddr:       req.WalletAddress,
		UserId:        *user_id,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		BirthDay:      req.BirthDay,
		Gender:        req.Gender,
		Nationality:   req.Nationality,
		IdFingerprint: req.IdFingerprint,
	}
	if err = user.Valid(); err != nil {
		return nil, err
	}
	err = m.userService.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &desc.NetWorkResponse{
		Msg:     "Successfully verified",
		Success: true,
	}, nil
}

func (m *MicroserviceServer) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*desc.UpdateUserResponse, error) {
	_, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}

	// updatedUser, err := m.userService.UpdateUser(dto.Person{
	// 	ID:          userID,
	// 	Email:       req.GetEmail(),
	// 	FirstName:   req.GetFirstName(),
	// 	LastName:    req.GetLastName(),
	// 	PhoneNumber: req.GetPhoneNumber()})
	if err != nil {
		return nil, err
	}
	return nil, nil
	// return &desc.UpdateUserResponse{Id: updatedUser.ID, FirstName: updatedUser.FirstName,
	// 	LastName: updatedUser.LastName, Email: updatedUser.Email, PhoneNumber: updatedUser.PhoneNumber}, nil
}

func (m *MicroserviceServer) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.NetWorkResponse, error) {
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		return nil, err
	}

	user, err := m.userService.GetUser(*uid, req.WalletAddress)
	if err != nil {
		return nil, fmt.Errorf("[ERR] you can't access other's user information")
	}

	return BuildRes(user, "Here is your business info", true)
}

func (m *MicroserviceServer) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*desc.NetWorkResponse, error) {
	uid, err := m.TokenInterceptor(ctx)
	if err != nil {
		log.Println("user isn't authorized")
	}

	_, err = m.userService.GetUser(*uid, req.WalletAddress)
	if err != nil {
		return nil, fmt.Errorf("[ERR] you can't access other's user information")
	}

	err = m.userService.DeleteUser(req.WalletAddress)
	if err != nil {
		return nil, err
	}
	return &desc.NetWorkResponse{
		Msg:     "Successfully deleted",
		Success: true,
	}, nil
}
