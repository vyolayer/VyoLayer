package service

import (
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/repository"
)

type UserService interface {
	GetUser(userId types.UserID) (*domain.User, ServiceError)
}

type userService struct {
	user repository.UserRepository
}

func NewUserService(user repository.UserRepository) UserService {
	return &userService{user: user}
}

func (us *userService) GetUser(userId types.UserID) (*domain.User, ServiceError) {
	userModel, err := us.user.FindById(userId)
	if err != nil {
		return nil, ServiceError(err)
	}

	if userModel == nil {
		return nil, domain.NewDomainError(404, "User not found")
	}

	return userModel, nil
}
