package service

import (
	"log"
	"time"
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"
	"worklayer/internal/repository"
	"worklayer/internal/utils/hash"

	"github.com/gofiber/fiber/v2"
)

type SessionService interface {
	SaveSession(ctx *fiber.Ctx, userId types.UserID, token string, expiry time.Duration) ServiceError
	DeleteSessionByToken(ctx *fiber.Ctx, token string) ServiceError
	RotateSession(ctx *fiber.Ctx, userId types.UserID, oldToken, newToken string, expiry time.Duration) (*domain.User, ServiceError)
}

type sessionService struct {
	user    repository.UserRepository
	session repository.SessionRepository
}

func NewSessionService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository) SessionService {
	return &sessionService{
		user:    userRepo,
		session: sessionRepo,
	}
}

func (ss *sessionService) SaveSession(ctx *fiber.Ctx, userId types.UserID, token string, expiry time.Duration) ServiceError {
	tokenHash := hash.HashToken(token)
	session := &models.UserSession{
		UserID:    userId.InternalID(),
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(expiry),
		IpAddress: ctx.IP(),
		UserAgent: ctx.Get("User-Agent"),
	}

	if err := ss.session.Save(session); err != nil {
		return NewServiceError(err.Code, err.Message)
	}
	return nil
}

func (ss *sessionService) DeleteSessionByToken(ctx *fiber.Ctx, token string) ServiceError {
	tokenHash := hash.HashToken(token)

	if err := ss.session.DeleteByTokenHash(tokenHash); err != nil {
		return NewServiceError(err.Code, err.Message)
	}
	return nil
}

func (ss *sessionService) RotateSession(ctx *fiber.Ctx, userId types.UserID, oldToken, newToken string, expiry time.Duration) (*domain.User, ServiceError) {
	oldTokenHash := hash.HashToken(oldToken)
	newTokenHash := hash.HashToken(newToken)
	newSession := &models.UserSession{
		UserID:    userId.InternalID(),
		TokenHash: newTokenHash,
		ExpiresAt: time.Now().Add(expiry),
		IpAddress: ctx.IP(),
		UserAgent: ctx.Get("User-Agent"),
	}

	if err := ss.session.RotateByTokenHash(oldTokenHash, newSession); err != nil {
		log.Printf("SESSION SERVICE :: RotateSession : %v", err.Message)
		return nil, NewServiceError(err.Code, err.Message)
	}

	user, err := ss.user.FindById(userId)
	if err != nil {
		return nil, NewServiceError(err.Code, err.Message)
	}
	return user, nil
}
