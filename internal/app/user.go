package app

import (
	"context"
	"fmt"
	"log"

	desc "ownify_api/pkg"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (m *MicroserviceServer[T]) Login(ctx context.Context, req *desc.SignInRequest) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	fmt.Println(md.Get("test"))

	token, err := m.authService.SignIn(req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	err = grpc.SendHeader(ctx, metadata.New(map[string]string{
		"Token": *token,
	}))
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (m *MicroserviceServer[T]) CreateUser(ctx context.Context, req *desc.SignUpRequest) (*desc.SignUpResponse, error) {
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

func (m *MicroserviceServer[T]) UpdateUser(ctx context.Context, req *desc.UpdateUserRequest) (*desc.UpdateUserResponse, error) {
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

func (m *MicroserviceServer[T]) GetUser(ctx context.Context, req *desc.GetUserRequest) (*desc.GetUserResponse, error) {
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

func (m *MicroserviceServer[T]) DeleteUser(ctx context.Context, req *desc.DeleteUserRequest) (*emptypb.Empty, error) {
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
