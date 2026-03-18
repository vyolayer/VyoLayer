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
	Create(ctx context.Context, projectID uuid.UUID, user *domain.User) *RepoError
	FindByID(ctx context.Context, projectID uuid.UUID, id uuid.UUID) (*domain.User, *RepoError)
	FindByEmail(ctx context.Context, projectID uuid.UUID, email string) (*domain.User, *RepoError)
	FindByUsername(ctx context.Context, projectID uuid.UUID, username string) (*domain.User, *RepoError)
	Update(ctx context.Context, projectID uuid.UUID, user *domain.User) *RepoError
	Delete(ctx context.Context, projectID uuid.UUID, id uuid.UUID) *RepoError
}

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) *RepoError
	FindByID(ctx context.Context, projectID uuid.UUID, id uuid.UUID) (*domain.Session, *RepoError)
	FindByTokenHash(ctx context.Context, projectID uuid.UUID, tokenHash string) (*domain.Session, *RepoError)
	FindByUserID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]*domain.Session, *RepoError)
	Delete(ctx context.Context, projectID uuid.UUID, id uuid.UUID) *RepoError
	DeleteExpired(ctx context.Context) *RepoError
	DeleteAll(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) *RepoError
	RotateToken(ctx context.Context, session *domain.Session) *RepoError
}

type VerificationTokenRepository interface {
	Create(ctx context.Context, token *domain.VerificationToken) *RepoError
	FindByTokenHash(ctx context.Context, projectID uuid.UUID, tokenHash string) (*domain.VerificationToken, *RepoError)
	FindByUserIDAndType(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, tokenType string) ([]*domain.VerificationToken, *RepoError)
	Delete(ctx context.Context, projectID uuid.UUID, id uuid.UUID) *RepoError
	DeleteExpired(ctx context.Context) *RepoError
	Update(ctx context.Context, projectID uuid.UUID, token *domain.VerificationToken) *RepoError
}
