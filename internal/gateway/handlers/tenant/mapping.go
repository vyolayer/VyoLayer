package tenant

import (
	tenantdto "github.com/vyolayer/vyolayer/internal/shared/dto/tenant"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
)

// Proto -> DTO (Organization)
func protoOrgToDTO(org *tenantV1.Organization) *tenantdto.Organization {
	if org == nil {
		return nil
	}

	return &tenantdto.Organization{
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

// Proto -> DTO (OrganizationResponse -> OrganizationDetailResponse)
func protoOrgResponseToDTO(resp *tenantV1.OrganizationResponse) *tenantdto.OrganizationDetailResponse {
	if resp == nil {
		return nil
	}

	members := make([]*tenantdto.OrganizationMember, len(resp.GetMembers()))
	for i, m := range resp.GetMembers() {
		members[i] = protoMemberToDTO(m)
	}

	return &tenantdto.OrganizationDetailResponse{
		Organization: protoOrgToDTO(resp.GetOrganization()),
		Members:      members,
	}
}

// Proto -> DTO (Organization member)
func protoMemberToDTO(m *tenantV1.OrganizationMember) *tenantdto.OrganizationMember {
	if m == nil {
		return nil
	}

	// roles := make([]*tenantdto.OrganizationRole, len(m.GetRoles()))
	roleNames := make([]string, len(m.GetRoles()))
	for i, r := range m.GetRoles() {
		// roles[i] = protoOrgRoleToDTO(r)
		roleNames[i] = r.GetName()
	}

	return &tenantdto.OrganizationMember{
		ID:            m.GetId(),
		UserID:        m.GetUserId(),
		FullName:      m.GetFullName(),
		Email:         m.GetEmail(),
		Status:        m.GetStatus(),
		JoinedAt:      m.GetJoinedAt(),
		InvitedAt:     m.GetInvitedAt(),
		InvitedBy:     m.GetInvitedBy(),
		DeactivatedBy: m.GetDeactivatedBy(),
		DeactivatedAt: m.GetDeactivatedAt(),
		Roles:         roleNames,
	}
}

// Proto -> DTO (Organization role)
func protoOrgRoleToDTO(r *tenantV1.OrganizationRole) *tenantdto.OrganizationRole {
	if r == nil {
		return nil
	}

	return &tenantdto.OrganizationRole{
		ID:           r.GetId(),
		Name:         r.GetName(),
		Description:  r.GetDescription(),
		IsSystemRole: r.GetIsSystemRole(),
		IsDefault:    r.GetIsDefault(),
	}
}

// Proto -> DTO (Organization permission)
func protoPermToDTO(p *tenantV1.OrganizationPermission) *tenantdto.OrganizationPerm {
	if p == nil {
		return nil
	}

	return &tenantdto.OrganizationPerm{
		ID:           p.GetId(),
		Resource:     p.GetResource(),
		Action:       p.GetAction(),
		Code:         p.GetCode(),
		Group:        p.GetGroup(),
		IsSystemPerm: p.GetIsSystem(),
	}
}

// Proto -> DTO (Organization invitation)
func protoInvitationToDTO(inv *tenantV1.OrganizationMemberInvitation) *tenantdto.OrganizationInvitation {
	if inv == nil {
		return nil
	}

	return &tenantdto.OrganizationInvitation{
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

// Proto -> DTO (Organization invitation for org)
func protoInvitationForOrgToDTO(inv *tenantV1.OrganizationMemberInvitationForOrg) *tenantdto.OrganizationInvitationForOrg {
	if inv == nil {
		return nil
	}

	invDto := protoInvitationToDTO(inv.GetInvitation())
	invByDto := &tenantdto.InvitedBy{
		MemberID: inv.GetInvitedBy().GetMemberId(),
		FullName: inv.GetInvitedBy().GetFullName(),
		Email:    inv.GetInvitedBy().GetEmail(),
	}

	return &tenantdto.OrganizationInvitationForOrg{
		Invitation: invDto,
		InvitedBy:  invByDto,
	}
}

// Proto -> DTO (Project)
func protoProjectToDTO(p *tenantV1.Project) *tenantdto.Project {

	if p == nil {
		return nil
	}

	return &tenantdto.Project{
		ID:             p.GetId(),
		OrganizationID: p.GetOrganizationId(),
		Name:           p.GetName(),
		Slug:           p.GetSlug(),
		Description:    p.GetDescription(),
		IsActive:       p.GetIsActive(),
		CreatedBy:      p.GetCreatedBy(),
		MaxAPIKeys:     p.GetMaxApiKeys(),
		MaxMembers:     p.GetMaxMembers(),
		MemberCount:    p.GetMemberCount(),
		CreatedAt:      p.GetCreatedAt(),
	}
}

// Proto -> DTO (ProjectResponse)
func protoProjectResponseToDTO(resp *tenantV1.ProjectResponse) *tenantdto.ProjectResponse {
	if resp == nil {
		return nil
	}

	members := make([]*tenantdto.ProjectMember, len(resp.GetMembers()))
	for i, m := range resp.GetMembers() {
		members[i] = protoProjectMemberToDTO(m)
	}

	return &tenantdto.ProjectResponse{
		Project: protoProjectToDTO(resp.GetProject()),
		Members: members,
	}
}

// Proto -> DTO (ProjectMember)
func protoProjectMemberToDTO(m *tenantV1.ProjectMember) *tenantdto.ProjectMember {
	if m == nil {
		return nil
	}

	return &tenantdto.ProjectMember{
		ID:        m.GetId(),
		UserID:    m.GetUserId(),
		Email:     m.GetEmail(),
		FullName:  m.GetFullName(),
		Role:      m.GetRole(),
		IsActive:  m.GetIsActive(),
		JoinedAt:  m.GetJoinedAt(),
		RemovedAt: m.RemovedAt, // *string – optional in proto3
	}
}
