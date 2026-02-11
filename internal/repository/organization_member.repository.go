package repository

import (
	"context"
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/mapper"
	"worklayer/pkg/errors"

	"gorm.io/gorm"
)

type OrganizationMemberRepository interface {
	GetByOrgID(ctx context.Context, userID UserID, orgID OrgID) ([]domain.OrganizationMember, *errors.AppError)
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
		Where("organization_id = ? AND user_id = ?", orgID.InternalID().String(), userID.InternalID().String()).
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
