package repository

import (
	"context"
	"time"
	"vyolayer/internal/platform/database/types"
	"vyolayer/pkg/errors"

	"gorm.io/gorm"
)

type ApiKeyUsageLogRepository interface {
	Log(ctx context.Context, log *TApiKeyUsageLog) *errors.AppError
	CountByApiKeyIDSince(ctx context.Context, apiKeyID types.ApiKeyID, since time.Time) (int64, *errors.AppError)
}

type apiKeyUsageLogRepository struct {
	db *gorm.DB
}

func NewApiKeyUsageLogRepository(db *gorm.DB) ApiKeyUsageLogRepository {
	return &apiKeyUsageLogRepository{db: db}
}

func (r *apiKeyUsageLogRepository) Log(
	ctx context.Context,
	log *TApiKeyUsageLog,
) *errors.AppError {
	err := r.db.WithContext(ctx).Create(log).Error
	if err != nil {
		return ConvertDBError(err, "logging API key usage")
	}
	return nil
}

func (r *apiKeyUsageLogRepository) CountByApiKeyIDSince(
	ctx context.Context,
	apiKeyID types.ApiKeyID,
	since time.Time,
) (int64, *errors.AppError) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&TApiKeyUsageLog{}).
		Where("api_key_id = ? AND created_at >= ?", apiKeyID.InternalID().ID(), since).
		Count(&count).Error
	if err != nil {
		return 0, ConvertDBError(err, "counting API key usage")
	}
	return count, nil
}
