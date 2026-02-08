package service

import (
	"worklayer/internal/app/dto"
	"worklayer/internal/domain"
	"worklayer/internal/repository"
	"worklayer/internal/utils/hash"
	"worklayer/pkg/errors"

	"github.com/gofiber/fiber/v2"
)

type AuthService interface {
	RegisterUser(ctx *fiber.Ctx, input dto.RegisterUserSchema) (*domain.User, *errors.AppError)
	LoginUser(ctx *fiber.Ctx, input dto.LoginUserSchema) (*domain.User, *errors.AppError)
}

type authService struct {
	user repository.UserRepository
}

func NewAuthService(user repository.UserRepository) AuthService {
	return &authService{user: user}
}

func (as *authService) RegisterUser(ctx *fiber.Ctx, input dto.RegisterUserSchema) (*domain.User, *errors.AppError) {
	user, err := domain.NewUser(input.Email, input.Password, input.FullName)
	if err != nil {
		return nil, err
	}

	// Check if user already exists
	existUser, _ := as.user.FindByEmail(user.Email)
	if existUser != nil {
		return nil, domain.UserAlreadyExistsError(input.Email)
	}

	// Create user
	userResponse, repoErr := as.user.CreateUser(*user)
	if repoErr != nil {
		return nil, WrapRepositoryError(repoErr, "register user")
	}

	return userResponse, nil
}

func (as *authService) LoginUser(ctx *fiber.Ctx, input dto.LoginUserSchema) (*domain.User, *errors.AppError) {
	email, err := domain.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	// Find user by email
	user, repoErr := as.user.FindByEmail(email.String())
	if repoErr != nil {
		// Don't leak information about whether user exists
		if errors.Is(repoErr, errors.ErrDBRecordNotFound) {
			return nil, domain.InvalidCredentialsError()
		}
		return nil, WrapRepositoryError(repoErr, "login user")
	}

	if user == nil {
		return nil, domain.InvalidCredentialsError()
	}

	// Verify user is active and verified
	if !user.IsActive {
		return nil, errors.NewWithMessage(errors.ErrUserInactive, "Account is inactive")
	}

	// Compare password
	if !hash.CheckPasswordHash(input.Password, user.HashedPassword) {
		return nil, domain.InvalidCredentialsError()
	}

	return user, nil
}
