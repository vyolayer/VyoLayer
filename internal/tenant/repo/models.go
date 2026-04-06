package tenantrepo

import (
	modelv1 "github.com/vyolayer/vyolayer/internal/tenant/models/v1"
)

type (
	// Organization
	Organization = modelv1.Organization

	// Organization RBAC
	OrganizationRole           = modelv1.OrganizationRole
	OrganizationPermission     = modelv1.OrganizationPermission
	OrganizationRolePermission = modelv1.OrganizationRolePermission

	// Organization Membership
	OrganizationMember           = modelv1.OrganizationMember
	OrganizationMemberInvitation = modelv1.OrganizationMemberInvitation
	MemberOrganizationRole       = modelv1.MemberOrganizationRole

	// Project
	Project = modelv1.Project

	// Project Membership
	ProjectMember = modelv1.ProjectMember

	// API Key
	APIKey = modelv1.ApiKey

	// Tenant Infra
	TenantInfra = modelv1.TenantInfra
)

const (
	APIKeyModeDev  = "dev"
	APIKeyModeLive = "live"
)
