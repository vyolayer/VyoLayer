package repository

import (
	"context"
	"vyolayer/internal/domain"
	"vyolayer/internal/platform/database/mapper"
	"vyolayer/internal/platform/database/types"
	"vyolayer/pkg/errors"

	"gorm.io/gorm"
)

type ApiKeyRepository interface {
	Create(ctx context.Context, apiKey *domain.ApiKey) (*domain.ApiKey, *errors.AppError)
	FindByID(ctx context.Context, apiKeyID types.ApiKeyID) (*domain.ApiKey, *errors.AppError)
	FindByKeyHash(ctx context.Context, keyHash string) (*domain.ApiKey, *errors.AppError)
	FindByProjectID(ctx context.Context, projectID types.ProjectID) ([]domain.ApiKey, *errors.AppError)
	Revoke(ctx context.Context, apiKeyID types.ApiKeyID, revokedBy types.UserID) *errors.AppError
	UpdateLastUsed(ctx context.Context, apiKeyID types.ApiKeyID) *errors.AppError
	CountByProjectID(ctx context.Context, projectID types.ProjectID) (int64, *errors.AppError)
}

type apiKeyRepository struct {
	db *gorm.DB
}

func NewApiKeyRepository(db *gorm.DB) ApiKeyRepository {
	return &apiKeyRepository{db: db}
}

func (r *apiKeyRepository) Create(
	ctx context.Context,
	apiKey *domain.ApiKey,
) (*domain.ApiKey, *errors.AppError) {
	model := TApiKey{
		OrganizationID: apiKey.OrganizationID.InternalID().ID(),
		ProjectID:      apiKey.ProjectID.InternalID().ID(),
		Name:           apiKey.Name,
		KeyPrefix:      apiKey.KeyPrefix,
		KeyHash:        apiKey.KeyHash,
		Mode:           apiKey.Mode,
		CreatedBy:      apiKey.CreatedBy.InternalID().ID(),
		RequestLimit:   apiKey.RequestLimit,
		RateLimit:      apiKey.RateLimit,
	}

	if apiKey.ExpiresAt != nil {
		model.ExpiresAt = apiKey.ExpiresAt
	}

	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return nil, ConvertDBError(err, "creating API key")
	}

	return mapper.ToDomainApiKey(&model), nil
}

func (r *apiKeyRepository) FindByID(
	ctx context.Context,
	apiKeyID types.ApiKeyID,
) (*domain.ApiKey, *errors.AppError) {
	var model TApiKey
	err := r.db.
		Where("id = ?", apiKeyID.InternalID().ID()).
		First(&model).Error
	if err != nil {
		return nil, ConvertDBError(err, "finding API key by ID")
	}
	return mapper.ToDomainApiKey(&model), nil
}

func (r *apiKeyRepository) FindByKeyHash(
	ctx context.Context,
	keyHash string,
) (*domain.ApiKey, *errors.AppError) {
	var model TApiKey
	err := r.db.
		Where("key_hash = ? AND revoked_at IS NULL", keyHash).
		First(&model).Error
	if err != nil {
		return nil, ConvertDBError(err, "finding API key by hash")
	}
	return mapper.ToDomainApiKey(&model), nil
}

func (r *apiKeyRepository) FindByProjectID(
	ctx context.Context,
	projectID types.ProjectID,
) ([]domain.ApiKey, *errors.AppError) {
	var models []TApiKey
	err := r.db.
		Where("project_id = ? AND deleted_at IS NULL", projectID.InternalID().ID()).
		Order("created_at DESC").
		Find(&models).Error
	if err != nil {
		return nil, ConvertDBError(err, "listing API keys by project ID")
	}

	result := make([]domain.ApiKey, 0, len(models))
	for _, m := range models {
		if ak := mapper.ToDomainApiKey(&m); ak != nil {
			result = append(result, *ak)
		}
	}
	return result, nil
}

func (r *apiKeyRepository) Revoke(
	ctx context.Context,
	apiKeyID types.ApiKeyID,
	revokedBy types.UserID,
) *errors.AppError {
	revokedByID := revokedBy.InternalID().ID()
	err := r.db.WithContext(ctx).
		Model(&TApiKey{}).
		Where("id = ?", apiKeyID.InternalID().ID()).
		Updates(map[string]interface{}{
			"revoked_at": gorm.Expr("NOW()"),
			"revoked_by": revokedByID,
		}).Error
	if err != nil {
		return ConvertDBError(err, "revoking API key")
	}
	return nil
}

func (r *apiKeyRepository) UpdateLastUsed(
	ctx context.Context,
	apiKeyID types.ApiKeyID,
) *errors.AppError {
	err := r.db.WithContext(ctx).
		Model(&TApiKey{}).
		Where("id = ?", apiKeyID.InternalID().ID()).
		Update("last_used_at", gorm.Expr("NOW()")).Error
	if err != nil {
		return ConvertDBError(err, "updating API key last used")
	}
	return nil
}

func (r *apiKeyRepository) CountByProjectID(
	ctx context.Context,
	projectID types.ProjectID,
) (int64, *errors.AppError) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&TApiKey{}).
		Where("project_id = ? AND revoked_at IS NULL AND deleted_at IS NULL",
			projectID.InternalID().ID()).
		Count(&count).Error
	if err != nil {
		return 0, ConvertDBError(err, "counting API keys by project ID")
	}
	return count, nil
}
