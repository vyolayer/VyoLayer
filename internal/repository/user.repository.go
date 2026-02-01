package repository

import (
	"log"
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/mapper"
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"

	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists  RepositoryError = NewRepositoryError(409, "User already exists")
	ErrUserNotFound       RepositoryError = NewRepositoryError(404, "User not found")
	ErrFailedToCreateUser RepositoryError = NewRepositoryError(500, "Failed to create user")
	ErrFailedToFindUser   RepositoryError = NewRepositoryError(500, "Failed to find user")
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(user domain.User) RepositoryError {
	err := ur.db.Create(&models.User{
		Email:           user.Email,
		PasswordHash:    user.HashedPassword,
		FullName:        user.FullName,
		IsActive:        user.IsActive,
		IsEmailVerified: user.IsEmailVerified,
	}).Error
	if err != nil {
		// check if the error is due to duplicate key
		if err == gorm.ErrDuplicatedKey {
			return ErrUserAlreadyExists
		}
		return NewRepositoryError(500, err.Error())
	}
	return nil
}

func (ur *userRepository) FindByEmail(email string) (*domain.User, RepositoryError) {
	var user models.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, NewRepositoryError(500, err.Error())
	}

	log.Printf("USER REPOSITORY :: User: %v", user)
	return mapper.ToDomainUser(&user), nil
}

func (ur *userRepository) FindById(id types.UserID) (*domain.User, RepositoryError) {
	var user models.User
	err := ur.db.Where("id = ?", id.InternalID().String()).First(&user).Error
	if err != nil {
		log.Printf("Error finding user by id: %v", err)
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, NewRepositoryError(500, err.Error())
	}

	log.Printf("USER REPOSITORY :: User: %v", user)
	return mapper.ToDomainUser(&user), nil
}
