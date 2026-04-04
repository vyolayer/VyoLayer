package tenantrepo

import (
	"time"

	"github.com/google/uuid"
	"github.com/vyolayer/vyolayer/internal/tenant/domain"
	tenantmodelv1 "github.com/vyolayer/vyolayer/internal/tenant/models/v1"
)

// --- Data Mappers (Domain -> GORM Models) ---

func toBaseModel(id uuid.UUID, createdAt, updatedAt time.Time) tenantmodelv1.BaseModel {
	return tenantmodelv1.BaseModel{
		ID: id,
		TimeStamps: tenantmodelv1.TimeStamps{
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		},
	}
}

func toOrgModel(d *domain.Organization) *tenantmodelv1.Organization {

	if d == nil {
		return nil
	}

	return &tenantmodelv1.Organization{
		BaseModel:     toBaseModel(d.ID, d.CreatedAt, d.UpdatedAt),
		Name:          d.Name,
		Slug:          d.Slug,
		Description:   d.Description,
		OwnerID:       d.OwnerID,
		IsActive:      d.IsActive,
		MaxProjects:   d.MaxProjects,
		MaxMembers:    d.MaxMembers,
		ProjectCount:  d.ProjectCount,
		MemberCount:   d.MemberCount,
		DeactivatedBy: d.DeactivatedBy,
		DeactivatedAt: d.DeactivatedAt,
		ArchivedAt:    d.ArchivedAt,
	}
}

func toMemberModel(d *domain.OrganizationMember) *tenantmodelv1.OrganizationMember {
	if d == nil {
		return nil
	}

	model := &tenantmodelv1.OrganizationMember{
		BaseModel:      toBaseModel(d.ID, d.CreatedAt, d.UpdatedAt),
		OrganizationID: d.OrganizationID,
		UserID:         d.UserID,
		InvitedBy:      d.InvitedBy,
		RemovedBy:      d.RemovedBy,
		RemovedAt:      d.RemovedAt,
	}

	if !d.JoinedAt.IsZero() {
		model.JoinedAt = d.JoinedAt
	}
	if d.InvitedAt != nil && !d.InvitedAt.IsZero() {
		model.InvitedAt = d.InvitedAt
	}

	return model
}

func toMemberRoleModel(d *domain.MemberOrganizationRole) *tenantmodelv1.MemberOrganizationRole {
	if d == nil {
		return nil
	}

	model := &tenantmodelv1.MemberOrganizationRole{
		BaseModel:      toBaseModel(d.ID, d.CreatedAt, d.UpdatedAt),
		MemberID:       d.MemberID,
		OrganizationID: d.OrganizationID,
		RoleID:         d.RoleID,
	}

	if d.GrantedBy != nil {
		model.GrantedBy = d.GrantedBy
	}
	if d.RevokedBy != nil {
		model.RevokedBy = d.RevokedBy
	}
	if d.RevokedAt != nil {
		model.RevokedAt = d.RevokedAt
	}
	if !d.GrantedAt.IsZero() {
		model.GrantedAt = d.GrantedAt
	}

	return model
}

func toOrganizationMemberInvitationModel(d *domain.OrganizationMemberInvitation) *OrganizationMemberInvitation {
	if d == nil {
		return nil
	}

	model := &OrganizationMemberInvitation{
		BaseModel:      toBaseModel(d.ID, d.CreatedAt, d.UpdatedAt),
		OrganizationID: d.OrganizationID,
		InvitedBy:      d.InvitedBy,
		Email:          d.Email,
		Token:          d.Token,
		IsAccepted:     d.IsAccepted,
		DeletedBy:      d.DeletedBy,
		InvitedAt:      d.InvitedAt,
		AcceptedAt:     d.AcceptedAt,
		ExpiredAt:      d.ExpiredAt,
	}

	return model
}

// toInvitationDomain converts a GORM OrganizationMemberInvitation model to a domain OrganizationMemberInvitation
func toInvitationDomain(m *OrganizationMemberInvitation) *domain.OrganizationMemberInvitation {
	if m == nil {
		return nil
	}
	return &domain.OrganizationMemberInvitation{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		InvitedBy:      m.InvitedBy,
		Email:          m.Email,
		Token:          m.Token,
		IsAccepted:     m.IsAccepted,
		DeletedBy:      m.DeletedBy,
		InvitedAt:      m.InvitedAt,
		AcceptedAt:     m.AcceptedAt,
		ExpiredAt:      m.ExpiredAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func toInvitationDomainWithInviter(m *OrganizationMemberInvitation) *domain.OrganizationMemberInvitationWithInviter {
	if m == nil {
		return nil
	}
	return &domain.OrganizationMemberInvitationWithInviter{
		OrganizationMemberInvitation: *toInvitationDomain(m),
		Inviter: domain.Inviter{
			MemberID: m.Inviter.ID,
			FullName: m.Inviter.User.FullName,
			Email:    m.Inviter.User.Email,
		},
	}
}

// --- Data Mappers (GORM Models -> Domain) ---

// toOrgDomain converts a GORM Organization model to a domain Organization
func toOrgDomain(m *tenantmodelv1.Organization) *domain.Organization {
	if m == nil {
		return nil
	}
	return &domain.Organization{
		ID:            m.ID,
		Name:          m.Name,
		Slug:          m.Slug,
		Description:   m.Description,
		IsActive:      m.IsActive,
		OwnerID:       m.OwnerID,
		MaxMembers:    m.MaxMembers,
		MaxProjects:   m.MaxProjects,
		ProjectCount:  m.ProjectCount,
		MemberCount:   m.MemberCount,
		DeactivatedBy: m.DeactivatedBy,
		DeactivatedAt: m.DeactivatedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		ArchivedAt:    m.ArchivedAt,
	}
}

// toOrgWithMemberDomain converts a GORM Organization model with preloaded members to a domain OrganizationWithMember
func toOrgWithMemberDomain(m *tenantmodelv1.Organization) *domain.OrganizationWithMember {
	if m == nil {
		return nil
	}

	org := toOrgDomain(m)

	// Map members
	var members []domain.OrganizationMember
	for _, dbMember := range m.Members {
		var joinedAt time.Time
		if dbMember.JoinedAt != nil {
			joinedAt = *dbMember.JoinedAt
		}

		var invitedAt time.Time
		if dbMember.InvitedAt != nil {
			invitedAt = *dbMember.InvitedAt
		}

		members = append(members, domain.OrganizationMember{
			ID:             dbMember.ID,
			OrganizationID: dbMember.OrganizationID,
			UserID:         dbMember.UserID,
			Email:          dbMember.User.Email,
			FullName:       dbMember.User.FullName,
			IsActive:       dbMember.IsActive(),
			RemovedAt:      dbMember.RemovedAt,
			RemovedBy:      dbMember.RemovedBy,
			JoinedAt:       &joinedAt,
			InvitedBy:      dbMember.InvitedBy,
			InvitedAt:      &invitedAt,
			CreatedAt:      dbMember.CreatedAt,
			UpdatedAt:      dbMember.UpdatedAt,
		})
	}

	return &domain.OrganizationWithMember{
		Organization: *org,
		Members:      members,
	}
}

// toMemberDomain converts a GORM OrganizationMember model to a domain OrganizationMember
func toMemberDomain(m *tenantmodelv1.OrganizationMember) *domain.OrganizationMember {
	if m == nil {
		return nil
	}

	domainMember := &domain.OrganizationMember{
		ID:             m.ID,
		OrganizationID: m.OrganizationID,
		UserID:         m.UserID,
		IsActive:       m.IsActive(),
		InvitedBy:      m.InvitedBy,
		RemovedBy:      m.RemovedBy,
		RemovedAt:      m.RemovedAt,
		InvitedAt:      m.InvitedAt, // Safe: both are *time.Time
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}

	// Safe pointer mapping
	if m.JoinedAt != nil {
		domainMember.JoinedAt = m.JoinedAt
	}

	// SAFE PRELOAD CHECK: Only map User fields if the User was actually joined/preloaded
	if m.User.ID != uuid.Nil {
		domainMember.Email = m.User.Email
		domainMember.FullName = m.User.FullName
	}

	return domainMember
}

// toMemberWithRolesDomain converts a GORM OrganizationMember model with preloaded roles
func toMemberWithRolesDomain(m *tenantmodelv1.OrganizationMember) *domain.OrganizationMemberWithRoles {
	if m == nil {
		return nil
	}

	domainMember := toMemberDomain(m)
	var roles []domain.OrganizationRole

	for _, junction := range m.Roles {
		roleDomain := toRoleDomain(&junction.Role)

		if roleDomain != nil {
			roles = append(roles, *roleDomain)
		}
	}

	return &domain.OrganizationMemberWithRoles{
		OrganizationMember: *domainMember,
		Roles:              roles,
	}
}

// toMemberWithRolesAndPermissionsDomain converts a GORM OrganizationMember model with preloaded roles and permissions to a domain OrganizationMemberWithRolesAndPermissions
func toMemberWithRolesAndPermissionsDomain(m *tenantmodelv1.OrganizationMember) *domain.OrganizationMemberWithRolesAndPermissions {
	if m == nil {
		return nil
	}

	domainMember := toMemberDomain(m)
	var roles []domain.OrganizationRole
	var permissions []domain.OrganizationPermission

	// Use a map to deduplicate permissions if a user has multiple roles with overlapping permissions
	permMap := make(map[uuid.UUID]bool)

	for _, junction := range m.Roles {
		// 1. Map the Role
		roles = append(roles, *toRoleDomain(&junction.Role))

		// 2. Map and deduplicate the Permissions attached to this Role
		for _, dbPerm := range junction.Role.Permissions {
			if !permMap[dbPerm.ID] {
				permMap[dbPerm.ID] = true
				permissions = append(permissions, *toPermissionDomain(&dbPerm))
			}
		}
	}

	return &domain.OrganizationMemberWithRolesAndPermissions{
		OrganizationMember: *domainMember,
		Roles:              roles,
		Permissions:        permissions,
	}
}

// --- Helper Mappers for Roles and Permissions ---

func toRoleDomain(r *tenantmodelv1.OrganizationRole) *domain.OrganizationRole {
	if r == nil {
		return nil
	}
	return &domain.OrganizationRole{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		IsSystem:    r.IsSystemRole,
		IsDefault:   r.IsDefault,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func toPermissionDomain(p *tenantmodelv1.OrganizationPermission) *domain.OrganizationPermission {
	if p == nil {
		return nil
	}
	return &domain.OrganizationPermission{
		ID:          p.ID,
		Resource:    p.Resource,
		Action:      p.Action,
		Code:        p.Code,
		Group:       p.Group,
		Description: p.Description,
		IsSystem:    p.IsSystem,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

// toMemberRoleDomain converts a GORM MemberOrganizationRole model to a domain MemberOrganizationRole
func toMemberRoleDomain(m *tenantmodelv1.MemberOrganizationRole) *domain.MemberOrganizationRole {
	if m == nil {
		return nil
	}
	d := &domain.MemberOrganizationRole{
		ID:             m.ID,
		MemberID:       m.MemberID,
		OrganizationID: m.OrganizationID,
		RoleID:         m.RoleID,
		GrantedAt:      m.GrantedAt,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}

	if m.GrantedBy != nil {
		d.GrantedBy = m.GrantedBy
	}
	if m.RevokedBy != nil {
		d.RevokedBy = m.RevokedBy
	}
	if m.RevokedAt != nil {
		d.RevokedAt = m.RevokedAt
	}

	return d
}
