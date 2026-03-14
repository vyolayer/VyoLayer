package repository

import (
	"context"
	"vyolayer/internal/domain"
	"vyolayer/internal/platform/database/mapper"
	"vyolayer/internal/platform/database/types"
	"vyolayer/pkg/errors"

	"gorm.io/gorm"
)

type ProjectMemberRepository interface {
	AddMember(ctx context.Context, member *domain.ProjectMember) (*domain.ProjectMember, *errors.AppError)
	FindByID(ctx context.Context, memberID types.ProjectMemberID) (*domain.ProjectMember, *errors.AppError)
	FindByProjectID(ctx context.Context, projectID types.ProjectID) ([]domain.ProjectMember, *errors.AppError)
	FindByUserAndProject(ctx context.Context, userID types.UserID, projectID types.ProjectID) (*domain.ProjectMember, *errors.AppError)
	UpdateRole(ctx context.Context, memberID types.ProjectMemberID, role string) *errors.AppError
	Remove(ctx context.Context, memberID types.ProjectMemberID, removedBy types.UserID) *errors.AppError
	CountByProjectID(ctx context.Context, projectID types.ProjectID) (int64, *errors.AppError)
}

type projectMemberRepository struct {
	db *gorm.DB
}

func NewProjectMemberRepository(db *gorm.DB) ProjectMemberRepository {
	return &projectMemberRepository{db: db}
}

func (r *projectMemberRepository) AddMember(
	ctx context.Context,
	member *domain.ProjectMember,
) (*domain.ProjectMember, *errors.AppError) {
	model := TProjectMember{
		ProjectID: member.ProjectID.InternalID().ID(),
		UserID:    member.UserID.InternalID().ID(),
		Role:      member.Role,
		AddedBy:   member.AddedBy.InternalID().ID(),
	}

	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return nil, ConvertDBError(err, "adding project member")
	}

	// Reload with user info
	var created TProjectMember
	loadErr := r.db.
		Where("id = ?", model.ID).
		Preload("User").
		First(&created).Error
	if loadErr != nil {
		return nil, ConvertDBError(loadErr, "loading created project member")
	}

	return mapper.ToDomainProjectMember(&created), nil
}

func (r *projectMemberRepository) FindByID(
	ctx context.Context,
	memberID types.ProjectMemberID,
) (*domain.ProjectMember, *errors.AppError) {
	var member TProjectMember
	err := r.db.
		Where("id = ?", memberID.InternalID().ID()).
		Preload("User").
		First(&member).Error
	if err != nil {
		return nil, ConvertDBError(err, "finding project member by ID")
	}
	return mapper.ToDomainProjectMember(&member), nil
}

func (r *projectMemberRepository) FindByProjectID(
	ctx context.Context,
	projectID types.ProjectID,
) ([]domain.ProjectMember, *errors.AppError) {
	var members []TProjectMember
	err := r.db.
		Where("project_id = ? AND removed_at IS NULL", projectID.InternalID().ID()).
		Preload("User").
		Find(&members).Error
	if err != nil {
		return nil, ConvertDBError(err, "listing project members")
	}

	result := make([]domain.ProjectMember, 0, len(members))
	for _, m := range members {
		if pm := mapper.ToDomainProjectMember(&m); pm != nil {
			result = append(result, *pm)
		}
	}
	return result, nil
}

func (r *projectMemberRepository) FindByUserAndProject(
	ctx context.Context,
	userID types.UserID,
	projectID types.ProjectID,
) (*domain.ProjectMember, *errors.AppError) {
	var member TProjectMember
	err := r.db.
		Where("user_id = ? AND project_id = ? AND removed_at IS NULL",
			userID.InternalID().ID(), projectID.InternalID().ID()).
		Preload("User").
		First(&member).Error
	if err != nil {
		return nil, ConvertDBError(err, "finding project member by user and project")
	}
	return mapper.ToDomainProjectMember(&member), nil
}

func (r *projectMemberRepository) UpdateRole(
	ctx context.Context,
	memberID types.ProjectMemberID,
	role string,
) *errors.AppError {
	err := r.db.WithContext(ctx).
		Model(&TProjectMember{}).
		Where("id = ?", memberID.InternalID().ID()).
		Update("role", role).Error
	if err != nil {
		return ConvertDBError(err, "updating project member role")
	}
	return nil
}

func (r *projectMemberRepository) Remove(
	ctx context.Context,
	memberID types.ProjectMemberID,
	removedBy types.UserID,
) *errors.AppError {
	removedByID := removedBy.InternalID().ID()
	err := r.db.WithContext(ctx).
		Model(&TProjectMember{}).
		Where("id = ?", memberID.InternalID().ID()).
		Updates(map[string]interface{}{
			"removed_at": gorm.Expr("NOW()"),
			"removed_by": removedByID,
		}).Error
	if err != nil {
		return ConvertDBError(err, "removing project member")
	}
	return nil
}

func (r *projectMemberRepository) CountByProjectID(
	ctx context.Context,
	projectID types.ProjectID,
) (int64, *errors.AppError) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&TProjectMember{}).
		Where("project_id = ? AND removed_at IS NULL AND deleted_at IS NULL",
			projectID.InternalID().ID()).
		Count(&count).Error
	if err != nil {
		return 0, ConvertDBError(err, "counting project members")
	}
	return count, nil
}
