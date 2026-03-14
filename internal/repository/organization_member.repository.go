package repository

import (
	"context"
	"time"
	"vyolayer/internal/domain"
	"vyolayer/internal/platform/database/mapper"
	"vyolayer/pkg/errors"

	"gorm.io/gorm"
)

type OrganizationMemberRepository interface {
	GetByOrgID(ctx context.Context, userID UserID, orgID OrgID) ([]domain.OrganizationMember, *errors.AppError)

	Create(
		ctx context.Context,
		userID UserID,
		orgID OrgID,
		roleIDs []string,
	) (*domain.OrganizationMember, *errors.AppError)

	GetByOrgIDAndMemberID(
		ctx context.Context,
		orgID OrgID,
		memberID OrgMemberID,
	) (*domain.OrganizationMember, *errors.AppError)

	GetCurrentMember(
		ctx context.Context,
		orgID OrgID,
		userID UserID,
	) (*domain.OrganizationMemberWithRBAC, *errors.AppError)

	Delete(
		ctx context.Context,
		memberID OrgMemberID,
	) *errors.AppError

	CountOwners(
		ctx context.Context,
		orgID OrgID,
	) (int64, *errors.AppError)

	AssignRole(
		ctx context.Context,
		memberID OrgMemberID,
		orgID OrgID,
		roleID OrgRoleID,
		grantedByMemberID OrgMemberID,
	) *errors.AppError

	RevokeAllRoles(
		ctx context.Context,
		memberID OrgMemberID,
		orgID OrgID,
	) *errors.AppError
}

type organizationMemberRepository struct {
	db *gorm.DB
}

func NewOrganizationMemberRepository(db *gorm.DB) OrganizationMemberRepository {
	return &organizationMemberRepository{db: db}
}

func (orm *organizationMemberRepository) GetByOrgID(
	ctx context.Context,
	userID UserID,
	orgID OrgID,
) ([]domain.OrganizationMember, *errors.AppError) {
	// First, verify that the requesting user is a member of the organization
	var userMembership TOrganizationMember
	checkErr := orm.db.
		Model(&TOrganizationMember{}).
		Where("organization_id = ? AND user_id = ? AND deleted_at IS NULL", orgID.InternalID(), userID.InternalID()).
		First(&userMembership).Error

	if checkErr != nil {
		if checkErr == gorm.ErrRecordNotFound {
			return nil, errors.Forbidden("You are not a member of this organization")
		}
		return nil, ConvertDBError(checkErr, "checking user membership")
	}

	// Fetch all members of the organization
	var org TOrganization
	result := orm.db.Model(&TOrganization{}).
		Where("id = ?", orgID.InternalID().String()).
		Preload("Members").
		Preload("Members.User").
		Find(&org)

	if result.Error != nil {
		return nil, ConvertDBError(result.Error, "getting organization members")
	}

	if len(org.Members) == 0 {
		return []domain.OrganizationMember{}, nil
	}

	// Convert to domain members
	domainMembers := make([]domain.OrganizationMember, len(org.Members))
	for i, member := range org.Members {
		domainMembers[i] = *mapper.ToDomainOrganizationMember(&member)
	}

	return domainMembers, nil
}

func (orm *organizationMemberRepository) Create(
	ctx context.Context,
	userID UserID,
	orgID OrgID,
	roleIDs []string,
) (*domain.OrganizationMember, *errors.AppError) {
	// First, verify that the requesting user is a member of the organization
	var userMembership TOrganizationMember
	checkErr := orm.db.WithContext(ctx).
		Model(&TOrganizationMember{}).
		Where("organization_id = ? AND user_id = ? AND deleted_at IS NULL", orgID.InternalID(), userID.InternalID()).
		First(&userMembership).Error

	if checkErr == nil {
		return nil, errors.Conflict("User is already a member of this organization")
	}

	if checkErr != gorm.ErrRecordNotFound {
		return nil, ConvertDBError(checkErr, "checking user membership")
	}

	// Create the organization member
	member := &TOrganizationMember{
		OrganizationID: orgID.InternalID().ID(),
		UserID:         userID.InternalID().ID(),
	}

	result := orm.db.WithContext(ctx).Create(member)
	if result.Error != nil {
		return nil, ConvertDBError(result.Error, "creating organization member")
	}

	if err := orm.db.Preload("User").First(member, member.ID).Error; err != nil {
		return nil, ConvertDBError(err, "loading created organization member")
	}

	return mapper.ToDomainOrganizationMember(member), nil
}

func (orm *organizationMemberRepository) GetByOrgIDAndMemberID(
	ctx context.Context,
	orgID OrgID,
	memberID OrgMemberID,
) (*domain.OrganizationMember, *errors.AppError) {
	var member TOrganizationMember
	result := orm.db.WithContext(ctx).
		Model(&TOrganizationMember{}).
		Preload("User").
		Where("organization_id = ? AND id = ? AND deleted_at IS NULL", orgID.InternalID(), memberID.InternalID()).
		First(&member)

	if result.Error != nil {
		return nil, ConvertDBError(result.Error, "getting organization member")
	}

	return mapper.ToDomainOrganizationMember(&member), nil
}

func (orm *organizationMemberRepository) GetCurrentMember(
	ctx context.Context,
	orgID OrgID,
	userID UserID,
) (*domain.OrganizationMemberWithRBAC, *errors.AppError) {
	var member TOrganizationMember
	result := orm.db.WithContext(ctx).
		Model(&TOrganizationMember{}).
		Preload("User").
		Preload("Roles").
		Preload("Roles.Role").
		Preload("Roles.Role.Permissions").
		Where("organization_id = ? AND user_id = ? AND deleted_at IS NULL", orgID.InternalID(), userID.InternalID()).
		First(&member)

	if result.Error != nil {
		return nil, ConvertDBError(result.Error, "getting organization member")
	}

	return mapper.ToDomainOrganizationMemberWithRBAC(&member), nil
}

func (orm *organizationMemberRepository) Delete(
	ctx context.Context,
	memberID OrgMemberID,
) *errors.AppError {
	now := time.Now()
	err := orm.db.WithContext(ctx).
		Model(&TOrganizationMember{}).
		Where("id = ? AND deleted_at IS NULL", memberID.InternalID()).
		Updates(map[string]interface{}{
			"removed_at": now,
			"deleted_at": now,
		}).Error

	if err != nil {
		return ConvertDBError(err, "soft deleting organization member")
	}

	return nil
}

func (orm *organizationMemberRepository) CountOwners(
	ctx context.Context,
	orgID OrgID,
) (int64, *errors.AppError) {
	var count int64
	err := orm.db.WithContext(ctx).
		Model(&TMemberOrganizationRole{}).
		Joins("JOIN organization_roles ON organization_roles.id = member_organization_roles.role_id").
		Joins("JOIN organization_members ON organization_members.id = member_organization_roles.member_id").
		Where("member_organization_roles.organization_id = ? AND organization_roles.name = ? AND organization_members.deleted_at IS NULL",
			orgID.InternalID().ID(), "Owner").
		Count(&count).Error

	if err != nil {
		return 0, ConvertDBError(err, "counting owners")
	}

	return count, nil
}

func (orm *organizationMemberRepository) AssignRole(
	ctx context.Context,
	memberID OrgMemberID,
	orgID OrgID,
	roleID OrgRoleID,
	grantedByMemberID OrgMemberID,
) *errors.AppError {
	now := time.Now()
	roleAssignment := &TMemberOrganizationRole{
		MemberID:       memberID.InternalID().ID(),
		OrganizationID: orgID.InternalID().ID(),
		RoleID:         roleID.InternalID().ID(),
		GrantedBy:      grantedByMemberID.InternalID().ID(),
		GrantedAt:      now,
	}

	if err := orm.db.WithContext(ctx).Create(roleAssignment).Error; err != nil {
		return ConvertDBError(err, "assigning role to member")
	}

	return nil
}

func (orm *organizationMemberRepository) RevokeAllRoles(
	ctx context.Context,
	memberID OrgMemberID,
	orgID OrgID,
) *errors.AppError {
	err := orm.db.WithContext(ctx).
		Where("member_id = ? AND organization_id = ?", memberID.InternalID().ID(), orgID.InternalID().ID()).
		Delete(&TMemberOrganizationRole{}).Error

	if err != nil {
		return ConvertDBError(err, "revoking all roles from member")
	}

	return nil
}
