package v1

import "github.com/gofiber/fiber/v2"

func (r *routes) registerOrganizationRoutes(router fiber.Router, d *dependencies) {
	org := router.Group("/organizations")
	org.Use(d.AuthMiddleware.JwtValidated())

	// ─── User-level access (no org membership required) ─────────────────────────
	org.Post("/", d.OrganizationCtrl.CreateOrganization)             // Create organization
	org.Get("/", d.OrganizationCtrl.ListOrganizations)               // List organizations
	org.Post("/onboarding", d.OrganizationCtrl.OnboardOrganization)  // Onboard organization
	org.Get("/slug/:slug", d.OrganizationCtrl.GetOrganizationBySlug) // Get organization by slug

	// Organization Invitation routes (user-level)
	org.Post("/invitations/pending", d.OrganizationMemInvCtrl.GetPendingInvitations) // List pending invitations for user
	org.Get("/invitations/accept", d.OrganizationMemInvCtrl.AcceptInvitation)        // Accept invitation by token

	// ─── Member-level access ────────────────────────────────────────────────────
	memberAccess := org.Group("/:orgId")
	memberAccess.Use(d.OrgMiddleware.CheckOrganizationMembership())

	memberAccess.Get("/", d.OrganizationCtrl.GetOrganizationByID)                // Get organization by ID
	memberAccess.Get("/members/me", d.OrganizationMemCtrl.CurrentMember)         // Current user's membership
	memberAccess.Post("/members/leave", d.OrganizationMemCtrl.LeaveOrganization) // Leave organization

	// ─── Admin-level access ─────────────────────────────────────────────────────
	admin := org.Group("/:orgId")
	admin.Use(d.OrgMiddleware.CheckOrganizationMembership(), d.OrgMiddleware.IsAdmin())

	// Organization management
	admin.Patch("/", d.OrganizationCtrl.UpdateOrganization)        // Update organization
	admin.Post("/archive", d.OrganizationCtrl.ArchiveOrganization) // Archive organization
	admin.Post("/restore", d.OrganizationCtrl.RestoreOrganization) // Restore organization

	// Member management
	admin.Get("/members", d.OrganizationMemCtrl.GetAllMembersByOrgID)                  // List all members
	admin.Get("/members/:memberId", d.OrganizationMemCtrl.GetMemberByOrgIDAndMemberID) // Get member by ID
	admin.Delete("/members/:memberId", d.OrganizationMemCtrl.RemoveMember)             // Remove member
	admin.Patch("/members/:memberId/role", d.OrganizationMemCtrl.ChangeRole)           // Change member role

	// Invitation management
	admin.Post("/invitations", d.OrganizationMemInvCtrl.CreateInvitation)                      // Create invitation
	admin.Get("/invitations", d.OrganizationMemInvCtrl.ListInvitations)                        // List invitations
	admin.Post("/invitations/:invitationId/resend", d.OrganizationMemInvCtrl.ResendInvitation) // Resend invitation
	admin.Delete("/invitations/:invitationId", d.OrganizationMemInvCtrl.CancelInvitation)      // Cancel invitation

	// RBAC
	admin.Get("/rbac/permissions", d.OrganizationRBACCtrl.GetAllPermissions) // List all permissions
	admin.Get("/rbac/roles", d.OrganizationRBACCtrl.GetAllRoles)             // List all roles

	// ─── Owner-level access ─────────────────────────────────────────────────────
	owner := org.Group("/:orgId")
	owner.Use(d.OrgMiddleware.CheckOrganizationMembership(), d.OrgMiddleware.IsOwner())

	owner.Delete("/", d.OrganizationCtrl.DeleteOrganization)                           // Delete organization
	owner.Post("/members/transfer-ownership", d.OrganizationMemCtrl.TransferOwnership) // Transfer ownership
}
