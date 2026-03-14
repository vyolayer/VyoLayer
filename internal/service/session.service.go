package service

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/domain"
	"github.com/vyolayer/vyolayer/internal/platform/database/models"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
	"github.com/vyolayer/vyolayer/internal/repository"
	"github.com/vyolayer/vyolayer/internal/utils/hash"
	"github.com/vyolayer/vyolayer/pkg/errors"
)

type SessionService interface {
	SaveSession(ctx *fiber.Ctx, userId types.UserID, token string, expiry time.Duration) *errors.AppError
	DeleteSessionByToken(ctx *fiber.Ctx, token string) *errors.AppError
	RotateSession(ctx *fiber.Ctx, userId types.UserID, oldToken, newToken string, expiry time.Duration) (*domain.User, *errors.AppError)
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

func (ss *sessionService) SaveSession(ctx *fiber.Ctx, userId types.UserID, token string, expiry time.Duration) *errors.AppError {
	tokenHash := hash.HashToken(token)
	session := &models.UserSession{
		UserID:    userId.InternalID().ID(),
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(expiry),
		IpAddress: ctx.IP(),
		UserAgent: ctx.Get("User-Agent"),
	}

	if err := ss.session.Save(session); err != nil {
		return WrapRepositoryError(err, "save session")
	}
	return nil
}

func (ss *sessionService) DeleteSessionByToken(ctx *fiber.Ctx, token string) *errors.AppError {
	tokenHash := hash.HashToken(token)

	if err := ss.session.DeleteByTokenHash(tokenHash); err != nil {
		return WrapRepositoryError(err, "delete session")
	}
	return nil
}

func (ss *sessionService) RotateSession(ctx *fiber.Ctx, userId types.UserID, oldToken, newToken string, expiry time.Duration) (*domain.User, *errors.AppError) {
	oldTokenHash := hash.HashToken(oldToken)
	newTokenHash := hash.HashToken(newToken)
	newSession := &models.UserSession{
		UserID:    userId.InternalID().ID(),
		TokenHash: newTokenHash,
		ExpiresAt: time.Now().Add(expiry),
		IpAddress: ctx.IP(),
		UserAgent: ctx.Get("User-Agent"),
	}

	if err := ss.session.RotateByTokenHash(oldTokenHash, newSession); err != nil {
		log.Printf("SESSION SERVICE :: RotateSession : %v", err.Message)
		return nil, WrapRepositoryError(err, "rotate session")
	}

	user, err := ss.user.FindById(userId)
	if err != nil {
		return nil, WrapRepositoryError(err, "get user after session rotation")
	}
	return user, nil
}
