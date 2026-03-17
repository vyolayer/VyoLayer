package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/account/domain"
	"gorm.io/gorm"
)

type userRepository struct {
	client *gorm.DB
}

func NewUserRepository(client *gorm.DB) UserRepository {
	return &userRepository{
		client: client,
	}
}

func (r *userRepository) Create(ctx context.Context, projectID uuid.UUID, user *domain.User) *RepoError {
	avatar := AvatarModel{
		UUID:          ModelID{ID: user.Avatar.ID},
		URL:           user.Avatar.URL,
		FallbackChar:  user.Avatar.FallbackChar,
		FallbackColor: user.Avatar.FallbackColor,
		TimeStamps: TimeStamps{
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: gorm.DeletedAt{Time: time.Time{}},
		},
	}

	um := UserModel{
		UUID:          ModelID{ID: user.ID},
		ProjectID:     projectID,
		Email:         user.Email,
		Username:      user.Username,
		Password:      user.HashedPassword,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		AvatarID:      user.Avatar.ID,
		EmailVerified: user.IsEmailVerified,
		Status:        user.Status.String(),
		TimeStamps: TimeStamps{
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: gorm.DeletedAt{Time: time.Time{}}, // null
		},
	}

	tx := r.client.Begin()
	if err := tx.Create(&avatar).Error; err != nil {
		return ConvertDBError(err, "Failed to create avatar")
	}

	if err := tx.Create(&um).Error; err != nil {
		return ConvertDBError(err, "Failed to create user")
	}

	if err := tx.Commit().Error; err != nil {
		return ConvertDBError(err, "Failed to commit transaction")
	}

	return nil
}

func (r *userRepository) FindByUsername(ctx context.Context, projectID uuid.UUID, username string) (*domain.User, *RepoError) {
	var um UserModel

	if err := r.client.
		Where("project_id = ? AND username = ?", projectID, username).
		First(&um).
		Error; err != nil {
		return nil, ConvertDBError(err, "Failed to find user")
	}

	return MapToDomainUser(&um), nil
}

func (r *userRepository) FindByID(ctx context.Context, projectID uuid.UUID, id uuid.UUID) (*domain.User, *RepoError) {
	var um UserModel
	if err := r.client.
		Preload("Avatar").
		Where("project_id = ? AND id = ?", projectID, id).
		First(&um).
		Error; err != nil {
		return nil, ConvertDBError(err, "Failed to find user")
	}

	return MapToDomainUser(&um), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, projectID uuid.UUID, email string) (*domain.User, *RepoError) {
	var um UserModel
	if err := r.client.
		Where("project_id = ? AND email = ?", projectID, email).
		First(&um).
		Error; err != nil {
		return nil, ConvertDBError(err, "Failed to find user")
	}

	return MapToDomainUser(&um), nil
}

func (r *userRepository) Update(ctx context.Context, projectID uuid.UUID, user *domain.User) *RepoError {
	updateData := &UserModel{
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		EmailVerified: user.IsEmailVerified,
		LastLoginAt:   user.LastLoginAt,
		Status:        user.Status.String(),
	}

	err := r.client.
		Model(&UserModel{}).
		Where("project_id = ? AND id = ?", projectID, user.ID).
		Select("Email", "FirstName", "LastName", "EmailVerified", "LastLoginAt", "UpdatedAt", "Status").
		Updates(updateData).Error
	if err != nil {
		return ConvertDBError(err, "Failed to update user")
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, projectID uuid.UUID, id uuid.UUID) *RepoError {
	if err := r.client.
		Where("project_id = ? AND id = ?", projectID, id).
		Delete(&UserModel{}).
		Error; err != nil {
		return ConvertDBError(err, "Failed to delete user")
	}

	return nil
}

func MapToDomainUser(u *UserModel) *domain.User {
	return &domain.User{
		ID:              u.UUID.ID,
		ProjectID:       u.ProjectID,
		Email:           u.Email,
		Username:        u.Username,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		HashedPassword:  u.Password,
		IsEmailVerified: u.EmailVerified,
		CreatedAt:       u.TimeStamps.CreatedAt,
		UpdatedAt:       u.TimeStamps.UpdatedAt,
		LastLoginAt:     u.LastLoginAt,
		Status:          domain.UserStatus(u.Status),
		Avatar: &domain.Avatar{
			ID:            u.Avatar.ID,
			URL:           u.Avatar.URL,
			FallbackChar:  u.Avatar.FallbackChar,
			FallbackColor: u.Avatar.FallbackColor,
		},
	}
}
