package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/account/domain"
	"gorm.io/gorm"
)

type verificationTokenRepository struct {
	client *gorm.DB
}

func NewVerificationTokenRepository(client *gorm.DB) VerificationTokenRepository {
	return &verificationTokenRepository{
		client: client,
	}
}

func (r *verificationTokenRepository) Create(ctx context.Context, token *domain.VerificationToken) error {
	tokenModel := VerificationTokenModel{
		UUID:      ModelID{ID: token.ID},
		ProjectID: token.ProjectID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		Type:      token.Type,
		CreatedAt: token.CreatedAt,
		ExpiresAt: token.ExpiresAt,
	}
	if err := r.client.Create(&tokenModel).Error; err != nil {
		return ConvertDBError(err, "Failed to create verification token")
	}
	return nil
}

func (r *verificationTokenRepository) FindByTokenHash(
	ctx context.Context,
	projectID uuid.UUID,
	tokenHash string,
) (*domain.VerificationToken, error) {
	var tokenModel VerificationTokenModel
	err := r.client.
		Where("token_hash = ? AND project_id = ?", tokenHash, projectID).
		First(&tokenModel).
		Error
	if err != nil {
		return nil, ConvertDBError(err, "Failed to find verification token")
	}

	return &domain.VerificationToken{
		ID:        tokenModel.UUID.ID,
		ProjectID: tokenModel.ProjectID,
		UserID:    tokenModel.UserID,
		TokenHash: tokenModel.TokenHash,
		Type:      tokenModel.Type,
		CreatedAt: tokenModel.CreatedAt,
		ExpiresAt: tokenModel.ExpiresAt,
		UsedAt:    tokenModel.UsedAt,
	}, nil
}

func (r *verificationTokenRepository) FindByUserIDAndType(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, tokenType string) ([]*domain.VerificationToken, error) {
	var tokens []*VerificationTokenModel
	err := r.client.
		Where("user_id = ? AND project_id = ? AND token_type = ?", userID, projectID, tokenType).
		Find(&tokens).
		Error
	if err != nil {
		return nil, ConvertDBError(err, "Failed to find verification tokens")
	}

	var verificationTokens []*domain.VerificationToken
	for _, token := range tokens {
		verificationTokens = append(verificationTokens, &domain.VerificationToken{
			ID:        token.UUID.ID,
			ProjectID: token.ProjectID,
			UserID:    token.UserID,
			TokenHash: token.TokenHash,
			Type:      token.Type,
			CreatedAt: token.CreatedAt,
			ExpiresAt: token.ExpiresAt,
			UsedAt:    &token.ExpiresAt,
		})
	}

	return verificationTokens, nil
}

func (r *verificationTokenRepository) Update(ctx context.Context, projectID uuid.UUID, token *domain.VerificationToken) error {
	tokenModel := VerificationTokenModel{
		UUID:      ModelID{ID: token.ID},
		ProjectID: token.ProjectID,
		UserID:    token.UserID,
		TokenHash: token.TokenHash,
		Type:      token.Type,
		CreatedAt: token.CreatedAt,
		ExpiresAt: token.ExpiresAt,
		UsedAt:    token.UsedAt,
	}

	if err := r.client.Save(&tokenModel).Error; err != nil {
		return ConvertDBError(err, "Failed to update verification token")
	}
	return nil
}

func (r *verificationTokenRepository) Delete(ctx context.Context, projectID uuid.UUID, id uuid.UUID) error {
	if err := r.client.
		Where("id = ? AND project_id = ?", id, projectID).
		Delete(&VerificationTokenModel{}).
		Error; err != nil {
		return ConvertDBError(err, "Failed to delete verification token")
	}

	return nil
}

func (r *verificationTokenRepository) DeleteExpired(ctx context.Context) error {
	if err := r.client.
		Where("expires_at < ?", time.Now()).
		Delete(&VerificationTokenModel{}).
		Error; err != nil {
		return ConvertDBError(err, "Failed to delete expired verification tokens")
	}

	return nil
}
