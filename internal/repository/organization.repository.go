package repository

import (
	"context"
	"time"
	"vyolayer/internal/domain"
	"vyolayer/internal/platform/database/mapper"
	"vyolayer/pkg/errors"

	"gorm.io/gorm"
)

type OrganizationRepository interface {
	ListByUserID(ctx context.Context, userID UserID) ([]domain.Organization, *errors.AppError)
	Create(ctx context.Context, org *domain.Organization) (*domain.Organization, *errors.AppError)
	FindByID(ctx context.Context, orgID OrgID) (*domain.Organization, *errors.AppError)
	FindBySlug(ctx context.Context, slug string) (*domain.Organization, *errors.AppError)
	Update(ctx context.Context, org *domain.Organization) *errors.AppError
	Delete(ctx context.Context, orgID OrgID) *errors.AppError
	SlugExists(ctx context.Context, slug string, excludeOrgID OrgID) (bool, *errors.AppError)
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (or *organizationRepository) ListByUserID(
	ctx context.Context,
	userID UserID,
) ([]domain.Organization, *errors.AppError) {
	var orgs []TOrganization
	err := or.db.
		Joins("JOIN organization_members om ON organizations.id = om.organization_id").
		Preload("Owner").
		Preload("Members").
		Where("om.user_id = ?", userID.InternalID().ID()).
		Find(&orgs).Error

	if err != nil {
		return nil, ConvertDBError(err, "listing organizations")
	}

	var domainOrgs []domain.Organization
	for _, org := range orgs {
		domainOrgs = append(domainOrgs, *mapper.ToDomainOrganizationWithMembers(&org))
	}

	return domainOrgs, nil
}

func (or *organizationRepository) Create(
	ctx context.Context,
	org *domain.Organization,
) (*domain.Organization, *errors.AppError) {
	tx := or.db.Begin()
	defer func() { tx.Rollback() }()

	if tx.Error != nil {
		return nil, ConvertDBError(tx.Error, "beginning transaction")
	}

	orgId := org.ID.InternalID().ID()
	ownerUserId := org.OwnerID.InternalID().ID()
	now := time.Now()

	// Create organization
	err := gorm.G[TOrganization](tx).Create(ctx, &TOrganization{
		BaseModel:   TBaseModel{ID: orgId},
		Name:        org.Name,
		Slug:        org.Slug,
		Description: org.Description,
		OwnerID:     ownerUserId,
	})

	if err != nil {
		return nil, ConvertDBError(err, "creating organization")
	}

	// Create the organization owner member
	memberErr := gorm.G[TOrganizationMember](tx).
		Create(ctx, &TOrganizationMember{
			OrganizationID: orgId,
			UserID:         ownerUserId,
			JoinedAt:       &now,
			InvitedBy:      nil,
			InvitedAt:      nil,
		})
	if memberErr != nil {
		return nil, ConvertDBError(memberErr, "creating organization owner member")
	}

	// Get the created member to assign role
	createdMember, memberErr := gorm.G[TOrganizationMember](tx).
		Where("organization_id = ? AND user_id = ?", orgId, ownerUserId).
		First(ctx)

	if memberErr != nil {
		return nil, ConvertDBError(memberErr, "getting created member")
	}

	// Find the "Owner" role
	ownerRole, roleErr := gorm.G[TOrganizationRole](tx).
		Where("name = ? AND is_system = ?", "Owner", true).
		First(ctx)

	if roleErr != nil {
		return nil, ConvertDBError(roleErr, "finding owner role")
	}

	// Assign owner role to the member
	roleAssignErr := gorm.G[TMemberOrganizationRole](tx).
		Create(ctx, &TMemberOrganizationRole{
			MemberID:       createdMember.ID,
			OrganizationID: orgId,
			RoleID:         ownerRole.ID,
			GrantedBy:      createdMember.ID, // Owner grants role to themselves
			GrantedAt:      now,
		})

	if roleAssignErr != nil {
		return nil, ConvertDBError(roleAssignErr, "assigning owner role")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, ConvertDBError(err, "committing transaction")
	}

	// Get the created organization
	createdOrg := TOrganization{}
	createdOrgErr := or.db.
		Where("id = ?", orgId).
		Preload("Members").
		Preload("Members.User").
		First(&createdOrg).Error

	if createdOrgErr != nil {
		return nil, ConvertDBError(createdOrgErr, "getting created organization")
	}

	return mapper.ToDomainOrganizationWithMembers(&createdOrg), nil
}

func (or *organizationRepository) FindByID(
	ctx context.Context,
	orgID OrgID,
) (*domain.Organization, *errors.AppError) {
	var org TOrganization
	err := or.db.
		Where("id = ?", orgID.InternalID().ID()).
		Preload("Owner").
		Preload("Members").
		Preload("Members.User").
		First(&org).Error

	if err != nil {
		return nil, ConvertDBError(err, "finding organization by ID")
	}

	return mapper.ToDomainOrganizationWithMembers(&org), nil
}

func (or *organizationRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (*domain.Organization, *errors.AppError) {
	var org TOrganization
	err := or.db.
		Where("slug = ?", slug).
		Preload("Owner").
		Preload("Members").
		Preload("Members.User").
		First(&org).Error

	if err != nil {
		return nil, ConvertDBError(err, "finding organization by slug")
	}

	return mapper.ToDomainOrganizationWithMembers(&org), nil
}

func (or *organizationRepository) Update(
	ctx context.Context,
	org *domain.Organization,
) *errors.AppError {
	orgID := org.ID.InternalID().ID()

	updates := map[string]interface{}{
		"name":        org.Name,
		"slug":        org.Slug,
		"description": org.Description,
		"is_active":   org.IsActive,
	}

	if org.DeactivatedBy != nil {
		deactivatedByID := (*org.DeactivatedBy).InternalID().ID()
		updates["deactivated_by"] = deactivatedByID
	} else {
		updates["deactivated_by"] = nil
	}

	if org.DeactivatedAt != nil {
		updates["deactivated_at"] = org.DeactivatedAt
	} else {
		updates["deactivated_at"] = nil
	}

	err := or.db.WithContext(ctx).
		Model(&TOrganization{}).
		Where("id = ?", orgID).
		Updates(updates).Error

	if err != nil {
		return ConvertDBError(err, "updating organization")
	}

	return nil
}

func (or *organizationRepository) Delete(
	ctx context.Context,
	orgID OrgID,
) *errors.AppError {
	err := or.db.WithContext(ctx).
		Unscoped().
		Where("id = ?", orgID.InternalID().ID()).
		Delete(&TOrganization{}).Error

	if err != nil {
		return ConvertDBError(err, "deleting organization")
	}

	return nil
}

func (or *organizationRepository) SlugExists(
	ctx context.Context,
	slug string,
	excludeOrgID OrgID,
) (bool, *errors.AppError) {
	var count int64
	err := or.db.WithContext(ctx).
		Model(&TOrganization{}).
		Where("slug = ? AND id != ?", slug, excludeOrgID.InternalID().ID()).
		Count(&count).Error

	if err != nil {
		return false, ConvertDBError(err, "checking slug existence")
	}

	return count > 0, nil
}
