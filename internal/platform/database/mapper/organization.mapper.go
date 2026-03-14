package mapper

import (
	"github.com/vyolayer/vyolayer/internal/domain"
	"github.com/vyolayer/vyolayer/internal/platform/database/models"
	"github.com/vyolayer/vyolayer/internal/platform/database/types"
)

// ToDomainOrganization converts an Organization model to an Organization domain object.
func ToDomainOrganization(orgModel *models.Organization) *domain.Organization {
	if orgModel == nil {
		return nil
	}

	var deactivatedBy *types.UserID
	if orgModel.DeactivatedBy != nil {
		id, _ := types.ReconstructUserID(orgModel.DeactivatedBy.String())
		deactivatedBy = &id
	}

	return domain.ReconstructOrganization(
		orgModel.PublicID(),
		orgModel.Name,
		orgModel.Slug,
		orgModel.Description,
		orgModel.Owner.PublicID(),
		orgModel.IsActive,
		deactivatedBy,
		orgModel.DeactivatedAt,
		orgModel.MaxProjects,
		orgModel.MaxMembers,
	)
}

// ToDomainOrganizationWithMembers converts an Organization model with members to domain
func ToDomainOrganizationWithMembers(orgModel *models.Organization) *domain.Organization {
	org := ToDomainOrganization(orgModel)
	if org == nil {
		return nil
	}

	// Convert members
	members := make([]domain.OrganizationMember, 0, len(orgModel.Members))
	for _, memberModel := range orgModel.Members {
		if member := ToDomainOrganizationMember(&memberModel); member != nil {
			members = append(members, *member)
		}
	}

	// Load members into organization
	org.LoadMembers(members)

	return org
}
