package service

import (
	"time"
	"worklayer/internal/app/dto"
	"worklayer/internal/repository"
	"worklayer/internal/utils/email"
	"worklayer/internal/utils/hash"
	"worklayer/internal/utils/response"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	RegisterUser(ctx *fiber.Ctx, input dto.RegisterUserDTO) ServiceError
	LoginUser(ctx *fiber.Ctx, input dto.LoginUserDTO) (*dto.UserDTO, ServiceError)
}

type authService struct {
	user repository.UserRepository
}

func NewAuthService(user repository.UserRepository) AuthService {
	return &authService{user: user}
}

func (as *authService) RegisterUser(ctx *fiber.Ctx, input dto.RegisterUserDTO) ServiceError {
	email := email.NewEmail(input.Email)
	hashedPassword, err := hash.HashPassword(input.Password)
	if err != nil {
		return NewServiceError(response.InternalServerError("Failed to hash password"))
	}

	// Check if user already exists
	existUser, _ := as.user.FindByEmail(email.Value())
	if existUser != nil {
		return NewServiceError(response.ConflictError("User already exists"))
	}

	// Create user
	if err := as.user.CreateUser(email.Value(), hashedPassword, input.FullName); err != nil {
		return NewServiceError(response.InternalServerError("Failed to create user"))
	}

	return nil
}

func (as *authService) LoginUser(ctx *fiber.Ctx, input dto.LoginUserDTO) (*dto.UserDTO, ServiceError) {
	email := email.NewEmail(input.Email)

	// Check if user already exists
	existUser, _ := as.user.FindByEmail(email.Value())
	if existUser == nil {
		return nil, NewServiceError(response.NotFoundError("User not found"))
	}

	// Compare password
	if !hash.CheckPasswordHash(input.Password, existUser.PasswordHash) {
		return nil, NewServiceError(response.UnauthorizedError("Invalid password"))
	}

	return &dto.UserDTO{
		ID:              existUser.ID,
		Email:           existUser.Email,
		IsActive:        existUser.IsActive,
		IsEmailVerified: existUser.IsEmailVerified,
		FullName:        existUser.FullName,
		CreatedAt:       existUser.CreatedAt.Format(time.RFC3339),
	}, nil
}
