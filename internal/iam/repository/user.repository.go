package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/iam/domain"
	m "github.com/vyolayer/vyolayer/internal/iam/models/v1"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"gorm.io/gorm"
)

// IAMUserRepository defines the persistence contract for IAM users.
type IAMUserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.IAMUser, error)
	FindByEmail(ctx context.Context, email string) (*domain.IAMUser, error)
	Create(ctx context.Context, user *domain.IAMUser) error
	Update(ctx context.Context, user *domain.IAMUser) error
}

type iAmUserRepository struct {
	log    *logger.AppLogger
	client *gorm.DB
}

func NewUserRepository(client *gorm.DB, log *logger.AppLogger) IAMUserRepository {
	return &iAmUserRepository{
		client: client,
		log:    log,
	}
}

func (r *iAmUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.IAMUser, error) {
	var u m.User
	if err := r.client.WithContext(ctx).
		Where("id = ?", id).
		First(&u).Error; err != nil {
		return nil, ConvertDBError(err, "Failed to find user")
	}

	r.log.Debug("Found user", u)

	return domain.ReconstructIAMUser(
		u.ID,
		u.Email, u.PasswordHash, u.FullName,
		u.IsEmailVerified, u.Status,
		u.CreatedAt, u.UpdatedAt,
	), nil
}

func (r *iAmUserRepository) FindByEmail(ctx context.Context, email string) (*domain.IAMUser, error) {
	var u m.User
	if err := r.client.WithContext(ctx).
		Where("email = ?", email).
		First(&u).Error; err != nil {
		return nil, ConvertDBError(err, "Failed to find user")
	}

	return domain.ReconstructIAMUser(
		u.ID,
		u.Email, u.PasswordHash, u.FullName,
		u.IsEmailVerified, u.Status,
		u.CreatedAt, u.UpdatedAt,
	), nil
}

func (r *iAmUserRepository) Create(ctx context.Context, user *domain.IAMUser) error {
	tx := r.client.WithContext(ctx).Begin()
	defer func() { tx.Rollback() }()

	now := time.Now()

	a := m.Avatar{
		URL:           user.Avatar.URL,
		FallbackChar:  user.Avatar.FallbackChar,
		FallbackColor: user.Avatar.FallbackColor,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := tx.Create(&a).Error; err != nil {
		return ConvertDBError(err, "Failed to create avatar")
	}

	u := m.User{
		ID:              user.ID,
		Email:           user.Email.String(),
		PasswordHash:    user.Password.String(),
		FullName:        user.FullName,
		IsEmailVerified: user.IsEmailVerified,
		Status:          user.Status.String(),
		AvatarID:        a.ID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := tx.Create(&u).Error; err != nil {
		return ConvertDBError(err, "Failed to create user")
	}

	if err := tx.Commit().Error; err != nil {
		return ConvertDBError(err, "Failed to commit transaction")
	}

	user.ID = u.ID
	return nil
}

// Update persists mutated fields (FullName, PasswordHash, Status) to the database.
func (r *iAmUserRepository) Update(ctx context.Context, user *domain.IAMUser) error {
	updates := map[string]any{
		"full_name":        user.FullName,
		"password_hash":    user.Password.String(),
		"is_email_verified": user.IsEmailVerified,
		"status":           user.Status.String(),
		"updated_at":       time.Now(),
	}

	result := r.client.WithContext(ctx).
		Model(&m.User{}).
		Where("id = ?", user.ID).
		Updates(updates)

	if result.Error != nil {
		return ConvertDBError(result.Error, "Failed to update user")
	}

	return nil
}
