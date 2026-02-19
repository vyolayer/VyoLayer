package v1

import "github.com/gofiber/fiber/v2"

func (r *routes) registerOrganizationRoutes(router fiber.Router, d *dependencies) {
	org := router.Group("/organizations")
	org.Use(d.AuthMiddleware.JwtValidated())

	// User access
	org.Post("/", d.OrganizationCtrl.CreateOrganization)             // Create organization
	org.Get("/", d.OrganizationCtrl.ListOrganizations)               // List organizations
	org.Post("/onboarding", d.OrganizationCtrl.OnboardOrganization)  // Onboard organization
	org.Get("/slug/:slug", d.OrganizationCtrl.GetOrganizationBySlug) // Get organization by slug

	// Organization Invitation routes (user)
	org.Post("/invitations/pending", d.OrganizationMemInvCtrl.GetPendingInvitations) // List invitations (pending) should be accessible by user
	org.Get("/invitations/accept", d.OrganizationMemInvCtrl.AcceptInvitation)        // Accept invitation by user

	// Organization member level access
	org.Get("/:orgId", d.OrgMiddleware.CheckOrganizationMembership(), d.OrganizationCtrl.GetOrganizationByID)         // Get organization by ID
	org.Get("/:orgId/members/me", d.OrgMiddleware.CheckOrganizationMembership(), d.OrganizationMemCtrl.CurrentMember) // User data as organization member

	// Admin level access
	admin := org.Group("/:orgId")
	admin.Use(d.OrgMiddleware.CheckOrganizationMembership(), d.OrgMiddleware.IsAdmin())

	admin.Get("/members", d.OrganizationMemCtrl.GetAllMembersByOrgID)                  // List all members
	admin.Get("/members/:memberId", d.OrganizationMemCtrl.GetMemberByOrgIDAndMemberID) // Get member by ID

	admin.Post("/invitations", d.OrganizationMemInvCtrl.CreateInvitation) // Create invitation
	admin.Get("/invitations", d.OrganizationMemInvCtrl.ListInvitations)   // List invitations

	admin.Get("/rbac/permissions", d.OrganizationRBACCtrl.GetAllPermissions) // List all permissions
	admin.Get("/rbac/roles", d.OrganizationRBACCtrl.GetAllRoles)             // List all roles
}
