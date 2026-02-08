package service

import (
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/repository"
	"worklayer/pkg/errors"
)

type UserService interface {
	GetUser(userId types.UserID) (*domain.User, *errors.AppError)
}

type userService struct {
	user repository.UserRepository
}

func NewUserService(user repository.UserRepository) UserService {
	return &userService{user: user}
}

func (us *userService) GetUser(userId types.UserID) (*domain.User, *errors.AppError) {
	userModel, err := us.user.FindById(userId)
	if err != nil {
		return nil, WrapRepositoryError(err, "get user")
	}

	if userModel == nil {
		return nil, domain.UserNotFoundError(userId.InternalID().String())
	}

	return userModel, nil
}
