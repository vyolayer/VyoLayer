package repository

import (
	"context"

	"github.com/vyolayer/vyolayer/internal/domain"
	"github.com/vyolayer/vyolayer/internal/platform/database/mapper"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
	"github.com/vyolayer/vyolayer/pkg/errors"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Create(
		ctx context.Context,
		project *domain.Project,
	) (*domain.Project, *errors.AppError)

	FindByID(
		ctx context.Context,
		projectID types.ProjectID,
	) (*domain.Project, *errors.AppError)

	FindBySlug(
		ctx context.Context,
		orgID types.OrganizationID,
		slug string,
	) (*domain.Project, *errors.AppError)

	FindByOrganizationID(
		ctx context.Context,
		orgID types.OrganizationID,
	) ([]domain.Project, *errors.AppError)

	Update(
		ctx context.Context,
		project *domain.Project,
	) *errors.AppError

	Delete(
		ctx context.Context,
		projectID types.ProjectID,
	) *errors.AppError

	CountByOrganizationID(
		ctx context.Context,
		orgID types.OrganizationID,
	) (int64, *errors.AppError)

	SlugExists(
		ctx context.Context,
		orgID types.OrganizationID,
		slug string,
		excludeProjectID types.ProjectID,
	) (bool, *errors.AppError)
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) Create(
	ctx context.Context,
	project *domain.Project,
) (*domain.Project, *errors.AppError) {
	tx := r.db.Begin()
	defer func() { tx.Rollback() }()

	if tx.Error != nil {
		return nil, ConvertDBError(tx.Error, "beginning transaction")
	}

	projectID := project.ID.InternalID().ID()
	orgID := project.OrganizationID.InternalID().ID()
	createdByID := project.CreatedBy.InternalID().ID()

	// Create the project
	err := gorm.G[TProject](tx).Create(ctx, &TProject{
		BaseModel:      TBaseModel{ID: projectID},
		OrganizationID: orgID,
		Name:           project.Name,
		Slug:           project.Slug,
		Description:    project.Description,
		IsActive:       true,
		CreatedBy:      createdByID,
		MaxApiKeys:     project.MaxApiKeys,
		MaxMembers:     project.MemberInfo.MaxNoOfMembers,
	})
	if err != nil {
		return nil, ConvertDBError(err, "creating project")
	}

	// Create the project creator as admin member
	members, _ := project.GetMembers()
	if len(members) > 0 {
		m := members[0]
		err = gorm.G[TProjectMember](tx).Create(ctx, &TProjectMember{
			ProjectID: projectID,
			UserID:    m.UserID.InternalID().ID(),
			Role:      m.Role,
			AddedBy:   createdByID,
		})
		if err != nil {
			return nil, ConvertDBError(err, "creating project admin member")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, ConvertDBError(err, "committing transaction")
	}

	// Reload
	var created TProject
	loadErr := r.db.
		Where("id = ?", projectID).
		Preload("Members").
		Preload("Members.User").
		First(&created).Error
	if loadErr != nil {
		return nil, ConvertDBError(loadErr, "loading created project")
	}

	return mapper.ToDomainProjectWithMembers(&created), nil
}

func (r *projectRepository) FindByID(
	ctx context.Context,
	projectID types.ProjectID,
) (*domain.Project, *errors.AppError) {
	var project TProject
	err := r.db.
		Where("id = ?", projectID.InternalID().ID()).
		Preload("Members").
		Preload("Members.User").
		First(&project).Error
	if err != nil {
		return nil, ConvertDBError(err, "finding project by ID")
	}
	return mapper.ToDomainProjectWithMembers(&project), nil
}

func (r *projectRepository) FindBySlug(
	ctx context.Context,
	orgID types.OrganizationID,
	slug string,
) (*domain.Project, *errors.AppError) {
	var project TProject
	err := r.db.
		Where("organization_id = ? AND slug = ?", orgID.InternalID().ID(), slug).
		Preload("Members").
		Preload("Members.User").
		First(&project).Error
	if err != nil {
		return nil, ConvertDBError(err, "finding project by slug")
	}
	return mapper.ToDomainProjectWithMembers(&project), nil
}

func (r *projectRepository) FindByOrganizationID(
	ctx context.Context,
	orgID types.OrganizationID,
) ([]domain.Project, *errors.AppError) {
	var projects []TProject
	err := r.db.
		Where("organization_id = ?", orgID.InternalID().ID()).
		Preload("Members").
		Preload("Members.User").
		Find(&projects).Error
	if err != nil {
		return nil, ConvertDBError(err, "listing projects by org ID")
	}

	result := make([]domain.Project, 0, len(projects))
	for _, p := range projects {
		if dp := mapper.ToDomainProjectWithMembers(&p); dp != nil {
			result = append(result, *dp)
		}
	}
	return result, nil
}

func (r *projectRepository) Update(
	ctx context.Context,
	project *domain.Project,
) *errors.AppError {
	updates := map[string]interface{}{
		"name":        project.Name,
		"slug":        project.Slug,
		"description": project.Description,
		"is_active":   project.IsActive,
	}

	err := r.db.WithContext(ctx).
		Model(&TProject{}).
		Where("id = ?", project.ID.InternalID().ID()).
		Updates(updates).Error
	if err != nil {
		return ConvertDBError(err, "updating project")
	}
	return nil
}

func (r *projectRepository) Delete(
	ctx context.Context,
	projectID types.ProjectID,
) *errors.AppError {
	err := r.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", projectID.InternalID().ID()).
		Delete(&TProject{}).Error
	if err != nil {
		return ConvertDBError(err, "deleting project")
	}
	return nil
}

func (r *projectRepository) CountByOrganizationID(
	ctx context.Context,
	orgID types.OrganizationID,
) (int64, *errors.AppError) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&TProject{}).
		Where("organization_id = ? AND deleted_at IS NULL", orgID.InternalID().ID()).
		Count(&count).Error
	if err != nil {
		return 0, ConvertDBError(err, "counting projects by org ID")
	}
	return count, nil
}

func (r *projectRepository) SlugExists(
	ctx context.Context,
	orgID types.OrganizationID,
	slug string,
	excludeProjectID types.ProjectID,
) (bool, *errors.AppError) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&TProject{}).
		Where("organization_id = ? AND slug = ? AND id != ?",
			orgID.InternalID().ID(), slug, excludeProjectID.InternalID().ID()).
		Count(&count).Error
	if err != nil {
		return false, ConvertDBError(err, "checking project slug existence")
	}
	return count > 0, nil
}
