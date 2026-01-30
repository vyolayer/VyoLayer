package repository

import (
	"worklayer/internal/platform/database/models"
)

var (
	user    = &models.User{}
	session = &models.UserSession{}
)

type UserRepository interface {
	CreateUser(email string, hashedPassword string, fullName string) error
	FindByEmail(email string) (*models.User, error)
	FindById(id uint) (*models.User, error)
}

type SessionRepository interface {
	Save(session *models.UserSession) error
	FindByUserId(userId uint) (*models.UserSession, error)
	FindByTokenHash(hashedToken string) (*models.UserSession, error)
	RotateByTokenHash(oldHashedToken string, newSession *models.UserSession) error
	DeleteByTokenHash(tokenHash string) error
}
