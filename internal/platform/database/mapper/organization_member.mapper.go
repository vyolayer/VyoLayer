package mapper

import (
	"worklayer/internal/domain"
	"worklayer/internal/platform/database/models"
	"worklayer/internal/platform/database/types"
)

// ToDomainOrganizationMember converts an OrganizationMember model to domain object
func ToDomainOrganizationMember(memberModel *models.OrganizationMember) *domain.OrganizationMember {
	if memberModel == nil {
		return nil
	}

	var invitedBy *types.OrganizationMemberID
	if memberModel.InvitedBy != nil {
		id, _ := types.ReconstructOrganizationMemberID(memberModel.InvitedBy.String())
		invitedBy = &id
	}

	var deactivatedBy *types.OrganizationMemberID
	if memberModel.RemovedBy != nil {
		id, _ := types.ReconstructOrganizationMemberID(memberModel.RemovedBy.String())
		deactivatedBy = &id
	}

	orgID, _ := types.ReconstructOrganizationID(memberModel.OrganizationID.String())
	userID, _ := types.ReconstructUserID(memberModel.UserID.String())

	// For joined at, use created time if not set
	joinedAt := memberModel.CreatedAt
	if memberModel.JoinedAt != nil {
		joinedAt = *memberModel.JoinedAt
	}

	return domain.ReconstructOrganizationMember(
		memberModel.PublicID(),
		orgID,
		userID,
		memberModel.User.Email,
		memberModel.User.FullName,
		memberModel.IsActive(),
		joinedAt,
		invitedBy,
		memberModel.InvitedAt,
		deactivatedBy,
		memberModel.RemovedAt,
	)
}
