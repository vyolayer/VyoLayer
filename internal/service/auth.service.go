package service

import (
	"worklayer/internal/app/dto"
	"worklayer/internal/domain"
	"worklayer/internal/repository"
	"worklayer/internal/utils/hash"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	RegisterUser(ctx *fiber.Ctx, input dto.RegisterUserSchema) (*domain.User, ServiceError)
	LoginUser(ctx *fiber.Ctx, input dto.LoginUserSchema) (*domain.User, ServiceError)
}

type authService struct {
	user repository.UserRepository
}

func NewAuthService(user repository.UserRepository) AuthService {
	return &authService{user: user}
}

func (as *authService) RegisterUser(ctx *fiber.Ctx, input dto.RegisterUserSchema) (*domain.User, ServiceError) {
	user, err := domain.NewUser(input.Email, input.Password, input.FullName)
	if err != nil {
		return nil, NewServiceError(err.Code, err.Message)
	}

	// Check if user already exists
	existUser, _ := as.user.FindByEmail(user.Email)
	if existUser != nil {
		return nil, NewServiceError(409, "User already exists")
	}

	// Create user
	userResponse, domainErr := as.user.CreateUser(*user)
	if domainErr != nil {
		return nil, NewServiceError(domainErr.Code, domainErr.Message)
	}

	// repository user created successfully
	userResponse, repoErr := as.user.FindById(userResponse.ID)
	if repoErr != nil {
		return nil, NewServiceError(repoErr.Code, repoErr.Message)
	}

	return userResponse, nil
}

func (as *authService) LoginUser(ctx *fiber.Ctx, input dto.LoginUserSchema) (*domain.User, ServiceError) {
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
