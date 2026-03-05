package v1

import "github.com/gofiber/fiber/v2"

func (r *routes) registerProjectRoutes(router fiber.Router, d *dependencies) {
	// All project routes are nested under organizations
	projects := router.Group("/organizations/:orgId/projects")
	projects.Use(d.AuthMiddleware.JwtValidated())
	projects.Use(d.OrgMiddleware.CheckOrganizationMembership())

	// ─── Project CRUD ────────────────────────────────────────────────────────
	projects.Post("/", d.ProjectCtrl.CreateProject) // Create project
	projects.Get("/", d.ProjectCtrl.ListProjects)   // List projects

	// ─── Project-specific routes ─────────────────────────────────────────────
	project := projects.Group("/:projectId")

	// Member-level access (any project member)
	project.Get("/", d.ProjectCtrl.GetProjectByID)                   // Get project
	project.Get("/members", d.ProjectMemberCtrl.ListMembers)         // List members
	project.Get("/members/me", d.ProjectMemberCtrl.GetCurrentMember) // Get my membership
	project.Post("/members/leave", d.ProjectMemberCtrl.LeaveProject) // Leave project
	project.Get("/api-keys", d.ApiKeyCtrl.ListKeys)                  // List API keys
	project.Get("/api-keys/:apiKeyId", d.ApiKeyCtrl.GetKeyByID)      // Get API key

	// Admin-level access
	project.Patch("/", d.ProjectCtrl.UpdateProject)                          // Update project
	project.Post("/archive", d.ProjectCtrl.ArchiveProject)                   // Archive project
	project.Post("/restore", d.ProjectCtrl.RestoreProject)                   // Restore project
	project.Delete("/", d.ProjectCtrl.DeleteProject)                         // Delete project
	project.Post("/members", d.ProjectMemberCtrl.AddMember)                  // Add member
	project.Patch("/members/:memberId/role", d.ProjectMemberCtrl.ChangeRole) // Change role
	project.Delete("/members/:memberId", d.ProjectMemberCtrl.RemoveMember)   // Remove member
	project.Post("/api-keys", d.ApiKeyCtrl.GenerateKey)                      // Generate API key
	project.Delete("/api-keys/:apiKeyId", d.ApiKeyCtrl.RevokeKey)            // Revoke API key
}
