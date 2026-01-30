package service

import (
	"time"
	"worklayer/internal/app/dto"
	"worklayer/internal/repository"
	"worklayer/internal/utils/response"
)

type UserService interface {
	GetUser(userId uint) (*dto.UserDTO, ServiceError)
}

type userService struct {
	user repository.UserRepository
}

func NewUserService(user repository.UserRepository) UserService {
	return &userService{user: user}
}

func (us *userService) GetUser(userId uint) (*dto.UserDTO, ServiceError) {
	user, err := us.user.FindById(userId)
	if err != nil {
		return nil, NewServiceError(response.InternalServerError("Failed to get user"))
	}

	return &dto.UserDTO{
		ID:              user.ID,
		Email:           user.Email,
		FullName:        user.FullName,
		IsActive:        user.IsActive,
		IsEmailVerified: user.IsEmailVerified,
		CreatedAt:       user.CreatedAt.Format(time.RFC3339),
	}, nil
}
