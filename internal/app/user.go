package app

import (
	"context"
	"log"

	"ownify_api/internal/dto"
	desc "ownify_api/pkg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)



func (m *MicroserviceServer) SignInWithPhone(ctx context.Context, req *desc.PhoneAuthRequest) (*desc.PhoneAuthResponse, error) {

	//Firebase token check
	firebaseToken, err := m.getUserIdFromToken(ctx)
	if err != nil {
		return nil, err
	}
	_, err = m.tokenManager.ValidateFirebase(firebaseToken)
	if err != nil {
		return nil, err
	}

	//Generate New UserId
	lastUserId, err := m.userService.GetLastUserId(req.WalletType)
	if err != nil {
		return nil, err
	}

	//Generate App Token for user
	newId := *lastUserId + 1
	token, err := m.authService.SignInWithPhone(firebaseToken, newId)
	if err != nil {
		return nil, err
	}

	//create user
	userId, err := m.userService.CreateUser(dto.BriefUser{
		ChainId:    int(req.ChainId),
		Wallet:     req.Wallet,
		WalletType: req.WalletType,
	})

	if err != nil {
		return nil, err
	}

	err = grpc.SendHeader(ctx, metadata.New(map[string]string{
		"Token": *token,
	}))
	if err != nil {
		return nil, err
	}
	return &desc.PhoneAuthResponse{
		UserId: *userId,
		Token:  *token,
	}, nil
}

func (m *MicroserviceServer) CreateUser(ctx context.Context, req *desc.SignUpRequest) (*desc.SignUpResponse, error) {
	// user := domain.Person{
	// 	FirstName:   req.GetFirstName(),
	// 	LastName:    req.GetLastName(),
	// 	Email:       req.GetEmail(),
	// 	Password:    req.GetPassword(),
	// 	PhoneNumber: req.GetPhoneNumber(),
	// 	Role:        "user",
	// }
	// id, err := m.authService.SignUp(user)
	// if err != nil {
	// 	return nil, err
	// }

	// return &desc.SignUpResponse{Id: *id}, nil
	return nil, nil
}

func (m *MicroserviceServer) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*desc.UpdateUserResponse, error) {
	_, err := m.getUserIdFromToken(ctx)
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

func (m *MicroserviceServer) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
	_, err := m.getUserIdFromToken(ctx)
	if err != nil {
		log.Println("user isn't authorized")
	}

	//_, err := m.userService.GetUser(req.GetId(), userID)
	if err != nil {
		return nil, err
	}
	return nil, nil

	// return &desc.GetUserResponse{
	// 	Id:          user.ID,
	// 	FirstName:   user.FirstName,
	// 	LastName:    user.LastName,
	// 	Email:       user.Email,
	// 	PhoneNumber: user.PhoneNumber,
	// 	Role:        string(user.Role),
	// 	Verified:    user.Verified,
	// 	Balance:     user.Balance}, nil
}

func (m *MicroserviceServer) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
	//userID, err := m.getUserIdFromToken(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// err = m.userService.DeleteUser(req.GetId(), 0)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

