package repository

import (
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"
	"worklayer/pkg/errors"
)

var (
	user    = &models.User{}
	session = &models.UserSession{}
)

type UserRepository interface {
	CreateUser(user domain.User) (*domain.User, *errors.AppError)
	FindByEmail(email string) (*domain.User, *errors.AppError)
	FindById(id types.UserID) (*domain.User, *errors.AppError)
}

type SessionRepository interface {
	Save(session *models.UserSession) *errors.AppError
	FindByUserId(userId types.UserID) (*models.UserSession, *errors.AppError)
	FindByTokenHash(hashedToken string) (*models.UserSession, *errors.AppError)
	RotateByTokenHash(oldHashedToken string, newSession *models.UserSession) *errors.AppError
	DeleteByTokenHash(tokenHash string) *errors.AppError
}
