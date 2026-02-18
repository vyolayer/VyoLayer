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

// ToDomainOrganizationMemberWithRBAC converts an OrganizationMember model to domain object
func ToDomainOrganizationMemberWithRBAC(memberModel *models.OrganizationMember) *domain.OrganizationMemberWithRBAC {
	if memberModel == nil {
		return nil
	}

	roles := make(map[string]domain.OrganizationRole)
	permissions := make(map[string]domain.OrganizationPermission)

	for _, role := range memberModel.Roles {
		for _, permission := range role.Role.Permissions {
			domainPermission := domain.OrganizationPermission{
				ID:       permission.PublicID().String(),
				Resource: permission.Resource,
				Action:   permission.Action,
				Group:    permission.Group,
				IsSystem: permission.IsSystem,
			}
			permissions[permission.PublicID().String()] = domainPermission
		}

		domainRole := domain.OrganizationRole{
			ID:           role.Role.PublicID().String(),
			Name:         role.Role.Name,
			Description:  role.Role.Description,
			IsSystemRole: role.Role.IsSystem,
			IsDefault:    role.Role.IsDefault,
		}
		roles[role.Role.PublicID().String()] = domainRole
	}

	return &domain.OrganizationMemberWithRBAC{
		OrganizationMember: *ToDomainOrganizationMember(memberModel),
		Roles:              roles,
		Permissions:        permissions,
	}
}
