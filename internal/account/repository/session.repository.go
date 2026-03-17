package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/account/domain"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"gorm.io/gorm"
)

type sessionRepository struct {
	client *gorm.DB
}

func NewSessionRepository(client *gorm.DB) SessionRepository {
	return &sessionRepository{
		client: client,
	}
}

func (r *sessionRepository) Create(ctx context.Context, session *domain.Session) *RepoError {
	s := SessionModel{
		UUID:      ModelID{ID: session.ID},
		ProjectID: session.ProjectID,
		UserID:    session.UserID,
		TokenHash: session.TokenHash,
		ExpiresAt: session.ExpiresAt,
		IpAddress: session.IPAddress,
		UserAgent: session.UserAgent,
		Reason:    session.RevokedReason,
		TimeStamps: TimeStamps{
			CreatedAt: session.CreatedAt,
			UpdatedAt: session.CreatedAt,
		},
	}

	if err := r.client.Create(&s).Error; err != nil {
		return ConvertDBError(err, "Failed to create session")
	}
	return nil
}

func (r *sessionRepository) FindByID(ctx context.Context, projectID uuid.UUID, id uuid.UUID) (*domain.Session, *RepoError) {
	var s SessionModel
	if err := r.client.
		Where("project_id = ? AND id = ?", projectID, id).
		First(&s).
		Error; err != nil {
		return nil, ConvertDBError(err, "Failed to find session")
	}
	return MapToDomainSession(&s), nil
}

func (r *sessionRepository) FindByUserID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]*domain.Session, *RepoError) {
	var ss []SessionModel

	if err := r.client.
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Find(&ss).
		Error; err != nil {
		return nil, ConvertDBError(err, "Failed to find sessions")
	}

	var sessions []*domain.Session
	for _, s := range ss {
		sessions = append(sessions, MapToDomainSession(&s))
	}
	return sessions, nil
}

func (r *sessionRepository) FindByTokenHash(ctx context.Context, projectID uuid.UUID, tokenHash string) (*domain.Session, *RepoError) {
	var s SessionModel
	if err := r.client.
		Where("project_id = ? AND token_hash = ?", projectID, tokenHash).
		First(&s).
		Error; err != nil {
		return nil, ConvertDBError(err, "Failed to find session")
	}
	return MapToDomainSession(&s), nil
}

func (r *sessionRepository) Delete(ctx context.Context, projectID uuid.UUID, id uuid.UUID) *RepoError {
	if err := r.client.
		Where("project_id = ? AND id = ?", projectID, id).
		Delete(&SessionModel{}).
		Error; err != nil {
		return ConvertDBError(err, "Failed to delete session")
	}
	return nil
}

func (r *sessionRepository) DeleteExpired(ctx context.Context) *RepoError {
	return errors.NotImplemented("DeleteExpired not implemented")
}

func MapToDomainSession(s *SessionModel) *domain.Session {
	return &domain.Session{
		ID:            s.UUID.ID,
		ProjectID:     s.ProjectID,
		UserID:        s.UserID,
		TokenHash:     s.TokenHash,
		IPAddress:     s.IpAddress,
		UserAgent:     s.UserAgent,
		ExpiresAt:     s.ExpiresAt,
		RevokedAt:     &s.CreatedAt,
		RevokedReason: s.Reason,
		CreatedAt:     s.TimeStamps.CreatedAt,
	}
}
