package repository

import (
	"log"

	"github.com/vyolayer/vyolayer/internal/domain"
	"github.com/vyolayer/vyolayer/internal/platform/database/mapper"
	"github.com/vyolayer/vyolayer/internal/platform/database/models"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(user domain.User) (*domain.User, *errors.AppError) {
	err := ur.db.Create(&models.User{
		Email:           user.Email,
		PasswordHash:    user.HashedPassword,
		FullName:        user.FullName,
		IsActive:        user.IsActive,
		IsEmailVerified: user.IsEmailVerified,
	}).Error

	if err != nil {
		return nil, ConvertDBError(err, "creating user")
	}

	return &user, nil
}

func (ur *userRepository) FindByEmail(email string) (*domain.User, *errors.AppError) {
	var user models.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		return nil, ConvertDBError(err, "finding user by email")
	}

	log.Printf("USER REPOSITORY :: User: %v", user)
	return mapper.ToDomainUser(&user), nil
}

func (ur *userRepository) FindById(id types.UserID) (*domain.User, *errors.AppError) {
	var user models.User
	err := ur.db.Where("id = ?", id.InternalID().String()).First(&user).Error
	if err != nil {
		log.Printf("Error finding user by id: %v", err)
		return nil, ConvertDBError(err, "finding user by ID")
	}

	log.Printf("USER REPOSITORY :: User: %v", user)
	return mapper.ToDomainUser(&user), nil
}
