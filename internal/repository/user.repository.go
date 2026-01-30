package repository

import (
	"log"
	"time"
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/models"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) CreateUser(email string, hashedPassword string, fullName string) error {
	return ur.db.Create(&models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		FullName:     fullName,
		LastLoginAt:  &time.Time{},
	}).Error
}

func (ur *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Printf("Error finding user by email: %v", err)
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) FindById(id domain.UserID) (*models.User, error) {
	var user models.User
	err := ur.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		log.Printf("Error finding user by id: %v", err)
		return nil, err
	}
	return &user, nil
}
