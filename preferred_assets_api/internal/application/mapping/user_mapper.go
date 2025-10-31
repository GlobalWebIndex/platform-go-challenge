package mapping

import (
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/application/dto"
	"github.com/MichailidouNatalia/GWI-Engineering-Challenge/preferred_assets_api/internal/domain"
	"github.com/google/uuid"
)

func UserReqToDomain(req dto.CreateUserRequest) (domain.User, error) {

	return domain.User{
		Id:       uuid.NewString(),
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}, nil
}

func UpdateReqToDomain(existingUser *domain.User, req dto.UpdateUserRequest) *domain.User {
	if req.Email != "" {
		existingUser.Email = req.Email
	}
	if req.Name != "" {
		existingUser.Name = req.Name
	}
	if req.Password != "" {
		existingUser.Password = req.Password
	}
	return existingUser
}

func DomainToUserRes(user domain.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}
}

func UserReqToResponseList(users []domain.User) []dto.UserResponse {
	responses := make([]dto.UserResponse, len(users))
	for i, user := range users {
		responses[i] = DomainToUserRes(user)
	}
	return responses
}
