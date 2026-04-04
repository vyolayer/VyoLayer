package tenantrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	"github.com/vyolayer/vyolayer/pkg/logger"
	"gorm.io/gorm"
)

type organizationMemberRepo struct {
	gormRepo
}

func NewOrganizationMemberRepo(db *gorm.DB, logger *logger.AppLogger) OrganizationMemberRepository {
	return &organizationMemberRepo{
		gormRepo: gormRepo{
			db:     db,
			logger: logger,
		},
	}
}

// --- Write Implementation ---

func (r *organizationMemberRepo) AddMember(ctx context.Context, tx *gorm.DB, member *domain.OrganizationMember) error {
	db := r.db
	if tx != nil {
		db = tx
	}

	model := toMemberModel(member)
	if err := db.WithContext(ctx).Create(model).Error; err != nil {
		r.logger.ErrorWithErr("AddMember", err)
		return ConvertDBError(err, "Failed to add member")
	}

	// Sync timestamps back to domain model
	member.SetCreatedAt(model.CreatedAt)
	member.SetUpdatedAt(model.UpdatedAt)
	r.logger.Debug("AddMember success", "")
	return nil
}

func (r *organizationMemberRepo) UpdateMember(ctx context.Context, member *domain.OrganizationMember) error {
	model := toMemberModel(member)

	return r.db.WithContext(ctx).
		Model(&OrganizationMember{}).
		Where("id = ?", member.ID).
		Omit("CreatedAt", "OrganizationID", "UserID"). // Prevent immutable fields from changing
		Updates(model).Error
}

func (r *organizationMemberRepo) RemoveMember(ctx context.Context, member *domain.OrganizationMember) error {
	now := time.Now()
	member.RemovedAt = &now
	member.IsActive = false

	// In BaaS, we usually update the status and removed_at rather than hard deleting the member
	// to preserve audit trails of who did what, when.
	// NOTE: OrganizationMember.IsActive is a method on the GORM model (not a field),
	// so we pass the raw bool value from the domain object directly.
	return r.db.WithContext(ctx).
		Model(&OrganizationMember{}).
		Where("id = ?", member.ID).
		Updates(map[string]any{
			"joined_at":  nil, // clear joined_at so IsActive() → false
			"removed_at": member.RemovedAt,
			"removed_by": member.RemovedBy,
			"updated_at": now,
			"deleted_at": now, // Trigger GORM soft delete if configured
		}).Error
}

// --- Read Implementation ---

func (r *organizationMemberRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.OrganizationMemberWithRolesAndPermissions, error) {
	var model OrganizationMember

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Roles.Role").             // Preload the Role via the MemberOrganizationRole junction
		Preload("Roles.Role.Permissions"). // Preload the Permissions attached to the Role
		Where("id = ?", id).
		First(&model).Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, nil
		// }
		return nil, ConvertDBError(err, "Failed to get member")
	}

	r.logger.Debug("Get Member by ID", model)
	member := toMemberWithRolesAndPermissionsDomain(&model)
	r.logger.Debug("Get Member by ID", member)
	return member, nil
}

func (r *organizationMemberRepo) GetByUserIdAndOrgId(ctx context.Context, userID, orgID uuid.UUID) (*domain.OrganizationMemberWithRolesAndPermissions, error) {
	var model OrganizationMember

	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Roles.Role").
		Preload("Roles.Role.Permissions").
		Where("user_id = ? AND organization_id = ?", userID, orgID).
		First(&model).Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, nil
		// }
		return nil, ConvertDBError(err, "Failed to get member")
	}

	return toMemberWithRolesAndPermissionsDomain(&model), nil
}

func (r *organizationMemberRepo) List(ctx context.Context, organizationID uuid.UUID) ([]*domain.OrganizationMemberWithRoles, error) {
	var (
		models []OrganizationMember
		result []*domain.OrganizationMemberWithRoles
	)

	// Only load the Roles, we usually don't need all Permissions for a list view
	err := r.db.WithContext(ctx).
		Preload("Roles.Role").
		Preload("User"). // Optional: if you want to extract FullName and Email from the IAM User model
		Where("organization_id = ?", organizationID).
		Find(&models).Error

	r.logger.Info("List", map[string]any{
		"organizationID": organizationID,
		"models":         models,
	})
	if err != nil {
		return nil, ConvertDBError(err, "Failed to list members")
	}

	result = make([]*domain.OrganizationMemberWithRoles, 0, len(models))
	for _, m := range models {
		r.logger.Debug("Member", m)
		if mapped := toMemberWithRolesDomain(&m); mapped != nil {
			r.logger.Debug("Mapped", mapped)
			result = append(result, mapped)
		}
	}

	r.logger.Debug("List result", map[string]any{
		"organizationID": organizationID,
		"result":         result,
	})

	return result, nil
}

func (r *organizationMemberRepo) GetByOrgIdAndEmail(ctx context.Context, orgID uuid.UUID, email string) (*domain.OrganizationMember, error) {
	var model OrganizationMember

	err := r.db.WithContext(ctx).
		Joins("JOIN iam.users u ON u.id = tenant.organization_members.user_id").
		Where("organization_id = ? AND u.email = ?", orgID, email).
		First(&model).Error

	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, nil
		// }
		return nil, ConvertDBError(err, "Failed to get member")
	}

	return toMemberDomain(&model), nil
}
