package seed

import "vyolayer/internal/platform/database/models"

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
	RoleOwner  = "Owner"
	RoleAdmin  = "Admin"
	RoleMember = "Member"
	RoleViewer = "Viewer"
)

// 1. Define Global Permissions
var OrgPermissions = []models.OrganizationPermission{
	// Organization
	{
		Resource:    ResourceOrganization,
		Action:      ActionRead,
		Description: "View organization details",
		Group:       "Settings",
		IsSystem:    true,
	},
	{
		Resource:    ResourceOrganization,
		Action:      ActionUpdate,
		Description: "Update settings",
		Group:       "Settings",
		IsSystem:    true,
	},
	{
		Resource:    ResourceOrganization,
		Action:      ActionDelete,
		Description: "Delete organization",
		Group:       "Settings",
		IsSystem:    true,
	},

	// Members Permissions
	{
		Resource:    ResourceMember,
		Action:      ActionInvite,
		Description: "Invite members to organization",
		Group:       "Team",
		IsSystem:    true,
	},
	{
		Resource:    ResourceMember,
		Action:      ActionRemove,
		Description: "Remove members from organization",
		Group:       "Team",
		IsSystem:    true,
	},
	{
		Resource:    ResourceMember,
		Action:      ActionList,
		Description: "List organization members",
		Group:       "Team",
		IsSystem:    true,
	},
	{
		Resource:    ResourceMember,
		Action:      ActionView,
		Description: "View member details",
		Group:       "Team",
		IsSystem:    true,
	},

	// Roles (RBAC)
	{
		Resource:    ResourceRole,
		Action:      ActionCreate,
		Description: "Create new roles",
		Group:       "Settings",
		IsSystem:    true,
	},
	{
		Resource:    ResourceRole,
		Action:      ActionUpdate,
		Description: "Update roles",
		Group:       "Settings",
		IsSystem:    true,
	},
	{
		Resource:    ResourceRole,
		Action:      ActionDelete,
		Description: "Delete roles",
		Group:       "Settings",
		IsSystem:    true,
	},
	{
		Resource:    ResourceRole,
		Action:      ActionView,
		Description: "View roles",
		Group:       "Settings",
		IsSystem:    true,
	},
	{
		Resource:    ResourceRole,
		Action:      ActionManage,
		Description: "Manage roles",
		Group:       "Settings",
		IsSystem:    true,
	},

	// Projects (Business Logic)
	{
		Resource:    ResourceProject,
		Action:      ActionCreate,
		Description: "Create projects",
		Group:       "Projects",
		IsSystem:    true,
	},
	{
		Resource:    ResourceProject,
		Action:      ActionRead,
		Description: "View projects",
		Group:       "Projects",
		IsSystem:    true,
	},
	{
		Resource:    ResourceProject,
		Action:      ActionUpdate,
		Description: "Edit projects",
		Group:       "Projects",
		IsSystem:    true,
	},
	{
		Resource:    ResourceProject,
		Action:      ActionDelete,
		Description: "Delete projects",
		Group:       "Projects",
		IsSystem:    true,
	},

	// Audit
	{
		Resource:    ResourceAudit,
		Action:      ActionRead,
		Description: "View audit logs",
		Group:       "Audit",
		IsSystem:    true,
	},
}

// 2. Define Global Roles (OrgID is NIL)
var OrgRoles = []models.OrganizationRole{
	{
		Name:        RoleOwner,
		Description: "Owner of the organization",
		IsSystem:    true,
		IsDefault:   false,
	}, // Owner is assigned on creation, not default join
	{
		Name:        RoleAdmin,
		Description: "Full administrative access",
		IsSystem:    true,
		IsDefault:   false,
	},
	{
		Name:        RoleMember,
		Description: "Standard member",
		IsSystem:    true,
		IsDefault:   true,
	}, // Default for new joins
	{
		Name:        RoleViewer,
		Description: "Read-only access",
		IsSystem:    true,
		IsDefault:   false,
	},
}
