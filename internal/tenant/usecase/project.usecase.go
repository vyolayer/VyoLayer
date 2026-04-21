package usecase

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"github.com/vyolayer/vyolayer/pkg/utils"
)

type projectUC struct {
	logger          *logger.AppLogger
	projectRepo     domain.ProjectRepository
	tenantInfraRepo domain.TenantInfraRepository
}

// NewProjectUseCase creates a new instance of the ProjectUseCase
func NewProjectUseCase(
	logger *logger.AppLogger,
	projectRepo domain.ProjectRepository,
	tenantInfraRepo domain.TenantInfraRepository,
) domain.ProjectUseCase {
	return &projectUC{
		logger:          logger,
		projectRepo:     projectRepo,
		tenantInfraRepo: tenantInfraRepo,
	}
}

func (u *projectUC) Create(ctx context.Context, orgID, createdBy uuid.UUID, name, description string) (*domain.Project, error) {
	// tenant infra exist if not then create
	tenantInfra, _ := u.tenantInfraRepo.GetByOrgID(ctx, orgID)
	if tenantInfra == nil || !tenantInfra.CompareStatus(domain.TenantInfraStatusReady) {
		return nil, status.Error(codes.FailedPrecondition, "tenant infra not found or not ready")
	}

	slug := utils.ToSlug(name).
		Slugify().
		AddSuffix(orgID.String()[:8]).
		String()

	// Construct the Domain entity
	project := domain.NewProject(orgID, createdBy, name, slug, description)

	// Save to Control Plane database
	if err := u.projectRepo.Create(ctx, nil, project); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create project: %s", err.Error())
	}

	return project, nil
}

func (u *projectUC) GetByID(ctx context.Context, projectID uuid.UUID) (*domain.Project, error) {
	project, err := u.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch project: %s", err.Error())
	}
	if project == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	return project, nil
}

func (u *projectUC) Get(ctx context.Context, orgID, projectID uuid.UUID) (*domain.Project, error) {
	project, err := u.projectRepo.GetByOrgID(ctx, orgID, projectID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch project: %s", err.Error())
	}
	if project == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	return project, nil
}

func (u *projectUC) List(ctx context.Context, orgID uuid.UUID, limit, offset int32) ([]*domain.Project, int32, error) {
	// Default pagination safety
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return u.projectRepo.List(ctx, orgID, limit, offset)
}

func (u *projectUC) Update(ctx context.Context, orgID, projectID uuid.UUID, name, description *string) (*domain.Project, error) {
	// Fetch existing project first to ensure it exists and we own it
	project, err := u.projectRepo.GetByOrgID(ctx, orgID, projectID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch project for update: %s", err.Error())
	}
	if project == nil {
		return nil, status.Error(codes.NotFound, "project not found")
	}

	isModified := false
	if name != nil && *name != project.GetName() {
		project.Name = *name
		isModified = true
	}
	if description != nil && *description != project.GetDescription() {
		project.Description = *description
		isModified = true
	}

	// Save if changes occurred
	if isModified {
		if err := u.projectRepo.Update(ctx, project); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to update project: %s", err.Error())
		}
	}

	return project, nil
}

func (u *projectUC) Delete(ctx context.Context, orgID, projectID uuid.UUID, confirmName string) error {
	// Fetch the project to validate the confirmation name
	project, err := u.projectRepo.GetByOrgID(ctx, orgID, projectID)
	if err != nil {
		return status.Errorf(codes.Internal, "failed to fetch project for deletion: %s", err.Error())
	}
	if project == nil {
		return status.Error(codes.NotFound, "project not found")
	}

	// Business Rule: The user must explicitly type the exact name to trigger a destructive action
	if project.GetName() != confirmName {
		u.logger.Warn("project deletion failed due to name mismatch", map[string]any{
			"projectID": projectID,
			"expected":  project.GetName(),
			"provided":  confirmName,
		})
		return status.Error(codes.InvalidArgument, "confirmation name does not match project name")
	}

	// Hard delete from the Control Plane DB
	return u.projectRepo.Delete(ctx, orgID, projectID)
}
