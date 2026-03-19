package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/account/domain"
	servicemodelv1 "github.com/vyolayer/vyolayer/pkg/postgres/models/service/account/v1"
	sharedmodel "github.com/vyolayer/vyolayer/pkg/postgres/models/shared"
)

type (
	// Models
	ModelID    = sharedmodel.UUID
	TimeStamps = servicemodelv1.TimeStamps

	UserModel              = servicemodelv1.ServiceUser
	AvatarModel            = servicemodelv1.ServiceUserAvatar
	SessionModel           = servicemodelv1.ServiceUserSession
	VerificationTokenModel = servicemodelv1.ServiceUserVerificationToken
	UserAccountLockModel   = servicemodelv1.ServiceUserAccountLock
	UserLoginAttemptModel  = servicemodelv1.ServiceUserLoginAttempt
)

type UserRepository interface {
	Create(ctx context.Context, projectID uuid.UUID, user *domain.User) error
	FindByID(ctx context.Context, projectID uuid.UUID, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, projectID uuid.UUID, email string) (*domain.User, error)
	FindByUsername(ctx context.Context, projectID uuid.UUID, username string) (*domain.User, error)
	Update(ctx context.Context, projectID uuid.UUID, user *domain.User) error
	Delete(ctx context.Context, projectID uuid.UUID, id uuid.UUID) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	FindByID(ctx context.Context, projectID uuid.UUID, id uuid.UUID) (*domain.Session, error)
	FindByTokenHash(ctx context.Context, projectID uuid.UUID, tokenHash string) (*domain.Session, error)
	FindByUserID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]*domain.Session, error)
	Delete(ctx context.Context, projectID uuid.UUID, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	DeleteAll(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error
	RotateToken(ctx context.Context, session *domain.Session) error
}

type VerificationTokenRepository interface {
	Create(ctx context.Context, token *domain.VerificationToken) error
	FindByTokenHash(ctx context.Context, projectID uuid.UUID, tokenHash string) (*domain.VerificationToken, error)
	FindByUserIDAndType(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, tokenType string) ([]*domain.VerificationToken, error)
	Delete(ctx context.Context, projectID uuid.UUID, id uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	Update(ctx context.Context, projectID uuid.UUID, token *domain.VerificationToken) error
}
