package grpc

// Data Mappers domain -> gRPC

import (
	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

// mapOrganizationToProto maps a domain.Organization to a tenantV1.Organization
func mapOrganizationToProto(org *domain.Organization) *tenantV1.Organization {
	return &tenantV1.Organization{
		Id:           org.GetIDString(),
		Name:         org.GetName(),
		Description:  org.GetDescription(),
		Slug:         org.GetSlug(),
		OwnerId:      org.GetOwnerIDString(),
		IsActive:     org.GetIsActive(),
		MaxProjects:  org.GetMaxProjects(),
		MaxMembers:   org.GetMaxMembers(),
		ProjectCount: org.GetProjectCount(),
		MemberCount:  org.GetMemberCount(),
		CreatedAt:    org.GetCreatedAtString(),
		UpdatedAt:    org.GetUpdatedAtString(),
	}
}

// mapOrganizationWithMembersToProto maps a domain.OrganizationWithMembers to a tenantV1.OrganizationResponse
func mapOrganizationWithMembersToProto(orgWithMembers *domain.OrganizationWithMember) *tenantV1.Organization {
	return &tenantV1.Organization{
		Id:           orgWithMembers.GetIDString(),
		Name:         orgWithMembers.GetName(),
		Description:  orgWithMembers.GetDescription(),
		Slug:         orgWithMembers.GetSlug(),
		OwnerId:      orgWithMembers.GetOwnerIDString(),
		IsActive:     orgWithMembers.GetIsActive(),
		MaxProjects:  orgWithMembers.GetMaxProjects(),
		MaxMembers:   orgWithMembers.GetMaxMembers(),
		ProjectCount: orgWithMembers.GetProjectCount(),
		MemberCount:  orgWithMembers.GetMemberCount(),
		CreatedAt:    orgWithMembers.GetCreatedAtString(),
		UpdatedAt:    orgWithMembers.GetUpdatedAtString(),
	}
}

// mapOrganizationMemberToProto maps a domain.OrganizationMember to a tenantV1.OrganizationMember
func mapOrganizationMemberToProto(member *domain.OrganizationMember) *tenantV1.OrganizationMember {
	if member == nil {
		return nil
	}

	// Initialize the proto struct with the guaranteed/required fields
	protoMember := &tenantV1.OrganizationMember{
		Id:       member.GetIDString(),
		UserId:   member.GetUserIDString(),
		FullName: member.GetFullName(),
		Email:    member.GetEmail(),
		Status:   member.GetStatus(),
		JoinedAt: member.GetJoinedAtString(),
	}

	if member.InvitedAt != nil && !member.InvitedAt.IsZero() {
		protoMember.InvitedAt = strPtr(member.GetInvitedAtString())
	}

	if member.InvitedBy != nil && *member.InvitedBy != uuid.Nil {
		protoMember.InvitedBy = strPtr(member.GetInvitedByString())
	}

	if member.RemovedBy != nil && *member.RemovedBy != uuid.Nil {
		protoMember.DeactivatedBy = strPtr(member.GetRemovedByString())
	}
	if member.RemovedAt != nil && !member.RemovedAt.IsZero() {
		protoMember.DeactivatedAt = strPtr(member.GetRemovedAtString())
	}

	return protoMember
}

func mapOrganizationMemberWithRolesToProto(member *domain.OrganizationMemberWithRoles) *tenantV1.OrganizationMember {
	protoMember := mapOrganizationMemberToProto(&member.OrganizationMember)

	roles := make([]*tenantV1.OrganizationRole, len(member.Roles))
	for i, role := range member.Roles {
		roles[i] = &tenantV1.OrganizationRole{
			Id:           role.ID.String(),
			Name:         role.Name,
			Description:  role.Description,
			IsSystemRole: role.IsSystem,
			IsDefault:    role.IsDefault,
		}
	}
	protoMember.Roles = roles
	return protoMember
}

func mapOrganizationMemberWithRolesAndPermissionsToProto(member *domain.OrganizationMemberWithRolesAndPermissions) *tenantV1.OrganizationMemberWithRBACResponse {
	var protoMember *tenantV1.OrganizationMemberWithRBACResponse

	permissions := make([]*tenantV1.OrganizationPermission, len(member.Permissions))
	for i, permission := range member.Permissions {
		permissions[i] = &tenantV1.OrganizationPermission{
			Id:       permission.ID.String(),
			Resource: permission.Resource,
			Action:   permission.Action,
			IsSystem: permission.IsSystem,
			Code:     permission.Resource + "." + permission.Action,
			Group:    permission.Group,
		}
	}

	protoMember = &tenantV1.OrganizationMemberWithRBACResponse{
		Member:      mapOrganizationMemberToProto(&member.OrganizationMember),
		Roles:       make([]*tenantV1.OrganizationRole, len(member.Roles)),
		Permissions: make([]*tenantV1.OrganizationPermission, len(member.Permissions)),
	}

	for i, role := range member.Roles {
		protoMember.Roles[i] = &tenantV1.OrganizationRole{
			Id:           role.ID.String(),
			Name:         role.Name,
			Description:  role.Description,
			IsSystemRole: role.IsSystem,
			IsDefault:    role.IsDefault,
		}
	}
	protoMember.Permissions = permissions
	return protoMember
}

func mapOrganizationMemberInvitationToProto(invitation *domain.OrganizationMemberInvitation) *tenantV1.OrganizationMemberInvitation {
	return &tenantV1.OrganizationMemberInvitation{
		Id:             invitation.GetIDString(),
		OrganizationId: invitation.GetOrganizationIDString(),
		Email:          invitation.GetEmail(),
		RoleIds:        invitation.GetRoleIDStrings(),
		InvitedBy:      invitation.GetInvitedByString(),
		IsAccepted:     invitation.GetIsAccepted(),
		AcceptedAt:     strPtr(invitation.GetAcceptedAtString()),
		ExpiredAt:      invitation.GetExpiredAtString(),
		InvitedAt:      invitation.GetInvitedAtString(),
		IsPending:      invitation.IsPending(),
	}
}

func mapInvitationWithInviterToProto(inv *domain.OrganizationMemberInvitationWithInviter) *tenantV1.OrganizationMemberInvitationForOrg {
	return &tenantV1.OrganizationMemberInvitationForOrg{
		Invitation: mapOrganizationMemberInvitationToProto(&inv.OrganizationMemberInvitation),
		InvitedBy: &tenantV1.InvitedBy{
			MemberId: inv.Inviter.MemberID.String(),
			FullName: inv.Inviter.FullName,
			Email:    inv.Inviter.Email,
		},
	}
}

func mapOrganizationPermissionToProto(p *domain.OrganizationPermission) *tenantV1.OrganizationPermission {
	return &tenantV1.OrganizationPermission{
		Id:       p.GetIDString(),
		Resource: p.GetResource(),
		Action:   p.GetAction(),
		Code:     p.GetCode(),
		Group:    p.GetGroup(),
		IsSystem: p.GetIsSystem(),
	}
}

func mapOrganizationRoleToProto(r *domain.OrganizationRole) *tenantV1.OrganizationRole {
	return &tenantV1.OrganizationRole{
		Id:           r.GetIDString(),
		Name:         r.GetName(),
		Description:  r.GetDescription(),
		IsSystemRole: r.GetIsSystem(),
		IsDefault:    r.GetIsDefault(),
	}
}
