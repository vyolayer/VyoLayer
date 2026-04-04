package seed

import tenantmodelv1 "github.com/vyolayer/vyolayer/internal/tenant/models/v1"

const (
	// Resources
	ResourceOrganization = "organization"
	ResourceMember       = "member"
	ResourceRole         = "role"
	ResourceProject      = "project"
	ResourceAudit        = "audit"

	// Actions
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionManage = "manage" // Special action for "assigning" things
	ActionInvite = "invite"
	ActionRemove = "remove"
	ActionList   = "list"
	ActionView   = "view"

	// Role Names
	RoleOwner  = "owner"
	RoleAdmin  = "admin"
	RoleMember = "member"
	RoleViewer = "viewer"
)

// 1. Define Global Permissions
var OrgPermissions = []tenantmodelv1.OrganizationPermission{
	// Organization
	{
		Resource:    ResourceOrganization,
		Action:      ActionRead,
		Description: "View organization details",
		Group:       "settings",
		Code:        "organization.read",
		IsSystem:    true,
	},
	{
		Resource:    ResourceOrganization,
		Action:      ActionUpdate,
		Description: "Update settings",
		Group:       "settings",
		Code:        "organization.update",
		IsSystem:    true,
	},
	{
		Resource:    ResourceOrganization,
		Action:      ActionDelete,
		Description: "Delete organization",
		Group:       "settings",
		Code:        "organization.delete",
		IsSystem:    true,
	},
	{
		Resource:    ResourceOrganization,
		Action:      "transfer",
		Description: "Transfer organization ownership",
		Group:       "settings",
		Code:        "organization.transfer",
		IsSystem:    true,
	},
	{
		Resource:    "billing",
		Action:      ActionManage,
		Description: "Manage billing and subscriptions",
		Group:       "billing",
		Code:        "billing.manage",
		IsSystem:    true,
	},

	// Members Permissions
	{
		Resource:    ResourceMember,
		Action:      ActionInvite,
		Description: "Invite members to organization",
		Group:       "team",
		Code:        "member.invite",
		IsSystem:    true,
	},
	{
		Resource:    ResourceMember,
		Action:      ActionRemove,
		Description: "Remove members from organization",
		Group:       "team",
		Code:        "member.remove",
		IsSystem:    true,
	},
	{
		Resource:    ResourceMember,
		Action:      ActionList,
		Description: "List organization members",
		Group:       "team",
		Code:        "member.list",
		IsSystem:    true,
	},
	{
		Resource:    ResourceMember,
		Action:      ActionView,
		Description: "View member details",
		Group:       "team",
		Code:        "member.view",
		IsSystem:    true,
	},
	{
		Resource:    ResourceMember,
		Action:      ActionUpdate,
		Description: "Update member roles and status",
		Group:       "team",
		Code:        "member.update",
		IsSystem:    true,
	},

	// Roles (RBAC)
	{
		Resource:    ResourceRole,
		Action:      ActionCreate,
		Description: "Create new roles",
		Group:       "settings",
		Code:        "role.create",
		IsSystem:    true,
	},
	{
		Resource:    ResourceRole,
		Action:      ActionUpdate,
		Description: "Update roles",
		Group:       "settings",
		Code:        "role.update",
		IsSystem:    true,
	},
	{
		Resource:    ResourceRole,
		Action:      ActionDelete,
		Description: "Delete roles",
		Group:       "settings",
		Code:        "role.delete",
		IsSystem:    true,
	},
	{
		Resource:    ResourceRole,
		Action:      ActionView,
		Description: "View roles",
		Group:       "settings",
		Code:        "role.view",
		IsSystem:    true,
	},
	{
		Resource:    ResourceRole,
		Action:      ActionManage,
		Description: "Manage roles",
		Group:       "settings",
		Code:        "role.manage",
		IsSystem:    true,
	},

	// Projects (Business Logic)
	{
		Resource:    ResourceProject,
		Action:      ActionCreate,
		Description: "Create projects",
		Group:       "projects",
		Code:        "project.create",
		IsSystem:    true,
	},
	{
		Resource:    ResourceProject,
		Action:      ActionRead,
		Description: "View projects",
		Group:       "projects",
		Code:        "project.read",
		IsSystem:    true,
	},
	{
		Resource:    ResourceProject,
		Action:      ActionUpdate,
		Description: "Edit projects",
		Group:       "projects",
		Code:        "project.update",
		IsSystem:    true,
	},
	{
		Resource:    ResourceProject,
		Action:      ActionDelete,
		Description: "Delete projects",
		Group:       "projects",
		Code:        "project.delete",
		IsSystem:    true,
	},

	// Audit
	{
		Resource:    ResourceAudit,
		Action:      ActionRead,
		Description: "View audit logs",
		Group:       "audit",
		Code:        "audit.read",
		IsSystem:    true,
	},
}

// 2. Define Global Roles (OrgID is NIL)
var OrgRoles = []tenantmodelv1.OrganizationRole{
	{
		Name:           RoleOwner,
		Description:    "Owner of the organization",
		IsSystemRole:   true,
		IsDefault:      false,
		HierarchyLevel: 100,
	}, // Owner is assigned on creation, not default join
	{
		Name:           RoleAdmin,
		Description:    "Full administrative access",
		IsSystemRole:   true,
		IsDefault:      false,
		HierarchyLevel: 80,
	},
	{
		Name:           RoleMember,
		Description:    "Standard member",
		IsSystemRole:   true,
		IsDefault:      true,
		HierarchyLevel: 50,
	}, // Default for new joins
	{
		Name:           RoleViewer,
		Description:    "Read-only access",
		IsSystemRole:   true,
		IsDefault:      false,
		HierarchyLevel: 10,
	},
}
