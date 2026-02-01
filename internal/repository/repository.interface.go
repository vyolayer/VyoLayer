package repository

import (
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"
)

var (
	user    = &models.User{}
	session = &models.UserSession{}
)

type UserRepository interface {
	CreateUser(user domain.User) RepositoryError
	FindByEmail(email string) (*domain.User, RepositoryError)
	FindById(id types.UserID) (*domain.User, RepositoryError)
}

type SessionRepository interface {
	Save(session *models.UserSession) RepositoryError
	FindByUserId(userId types.UserID) (*models.UserSession, RepositoryError)
	FindByTokenHash(hashedToken string) (*models.UserSession, RepositoryError)
	RotateByTokenHash(oldHashedToken string, newSession *models.UserSession) RepositoryError
	DeleteByTokenHash(tokenHash string) RepositoryError
}
