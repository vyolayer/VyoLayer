package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/vyolayer/vyolayer/internal/gateway/handlers/dto"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

func getOrgIDFromLocals(c *fiber.Ctx) string {
	orgID, _ := c.Locals("organization_id").(string)
	return orgID
}

func protoOrgToDTO(org *tenantV1.Organization) *dto.TOrganization {
	if org == nil {
		return nil
	}
	return &dto.TOrganization{
		ID:           org.GetId(),
		Name:         org.GetName(),
		Slug:         org.GetSlug(),
		Description:  org.GetDescription(),
		IsActive:     org.GetIsActive(),
		OwnerID:      org.GetOwnerId(),
		MaxMembers:   org.GetMaxMembers(),
		MaxProjects:  org.GetMaxProjects(),
		ProjectCount: org.GetProjectCount(),
		MemberCount:  org.GetMemberCount(),
		CreatedAt:    org.GetCreatedAt(),
		UpdatedAt:    org.GetUpdatedAt(),
	}
}

func protoMemberToDTO(m *tenantV1.OrganizationMember) *dto.TOrganizationMember {
	if m == nil {
		return nil
	}
	roles := make([]*dto.TOrganizationRole, len(m.GetRoles()))
	roleNames := make([]string, len(m.GetRoles()))
	for i, r := range m.GetRoles() {
		roles[i] = &dto.TOrganizationRole{
			ID:           r.GetId(),
			Name:         r.GetName(),
			Description:  r.GetDescription(),
			IsSystemRole: r.GetIsSystemRole(),
			IsDefault:    r.GetIsDefault(),
		}
		roleNames[i] = r.GetName()
	}
	return &dto.TOrganizationMember{
		ID:       m.GetId(),
		UserID:   m.GetUserId(),
		FullName: m.GetFullName(),
		Email:    m.GetEmail(),
		Status:   m.GetStatus(),
		JoinedAt: m.GetJoinedAt(),
		Roles:    roleNames,
	}
}

func protoInvitationToDTO(inv *tenantV1.OrganizationMemberInvitation) *dto.TOrganizationInvitation {
	if inv == nil {
		return nil
	}
	return &dto.TOrganizationInvitation{
		ID:             inv.GetId(),
		OrganizationID: inv.GetOrganizationId(),
		Email:          inv.GetEmail(),
		RoleIDs:        inv.GetRoleIds(),
		InvitedBy:      inv.GetInvitedBy(),
		InvitedAt:      inv.GetInvitedAt(),
		IsAccepted:     inv.GetIsAccepted(),
		AcceptedAt:     inv.GetAcceptedAt(),
		ExpiredAt:      inv.GetExpiredAt(),
		IsPending:      inv.GetIsPending(),
	}
}

func protoInvitationForOrgToDTO(inv *tenantV1.OrganizationMemberInvitationForOrg) *dto.TOrganizationInvitationForOrg {
	if inv == nil {
		return nil
	}

	invDto := protoInvitationToDTO(inv.GetInvitation())
	invByDto := &dto.TInvitedBy{
		MemberID: inv.GetInvitedBy().GetMemberId(),
		FullName: inv.GetInvitedBy().GetFullName(),
		Email:    inv.GetInvitedBy().GetEmail(),
	}

	return &dto.TOrganizationInvitationForOrg{
		Invitation: invDto,
		InvitedBy:  invByDto,
	}
}
