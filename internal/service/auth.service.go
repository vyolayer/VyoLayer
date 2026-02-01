package service

import (
	"worklayer/internal/app/dto"
	"worklayer/internal/domain"
	"worklayer/internal/repository"
	"worklayer/internal/utils/hash"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	RegisterUser(ctx *fiber.Ctx, input dto.RegisterUserDTO) ServiceError
	LoginUser(ctx *fiber.Ctx, input dto.LoginUserDTO) (*domain.User, ServiceError)
}

type authService struct {
	user repository.UserRepository
}

func NewAuthService(user repository.UserRepository) AuthService {
	return &authService{user: user}
}

func (as *authService) RegisterUser(ctx *fiber.Ctx, input dto.RegisterUserDTO) ServiceError {
	user, err := domain.NewUser(input.Email, input.Password, input.FullName)
	if err != nil {
		return NewServiceError(err.Code, err.Message)
	}

	// Check if user already exists
	existUser, _ := as.user.FindByEmail(user.Email)
	if existUser != nil {
		return NewServiceError(409, "User already exists")
	}

	// Create user
	if err := as.user.CreateUser(*user); err != nil {
		return NewServiceError(err.Code, err.Message)
	}

	return nil
}

func (as *authService) LoginUser(ctx *fiber.Ctx, input dto.LoginUserDTO) (*domain.User, ServiceError) {
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return nil, NewServiceError(err.Code, err.Message)
	}

	// Check if user already exists
	user, repoErr := as.user.FindByEmail(email.String())
	if repoErr != nil {
		return nil, NewServiceError(repoErr.Code, repoErr.Message)
	}
	if user == nil {
		return nil, NewServiceError(404, "User not found")
	}

	// Compare password
	if !hash.CheckPasswordHash(input.Password, user.HashedPassword) {
		return nil, NewServiceError(401, "Invalid password")
	}

	return user, nil
}
