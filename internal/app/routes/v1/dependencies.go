package v1

import (
	"worklayer/internal/app/controller"
	"worklayer/internal/app/middleware"
	"worklayer/internal/repository"
	"worklayer/internal/service"
)

type dependencies struct {
	// Controllers
	HealthCtrl             *controller.HealthController
	AuthCtrl               controller.AuthController
	UserCtrl               controller.UserController
	OrganizationCtrl       controller.OrganizationController
	OrganizationMemCtrl    controller.OrganizationMemberController
	OrganizationMemInvCtrl controller.OrganizationMemberInvitationController
	OrganizationRBACCtrl   controller.OrganizationRBACController
	ProjectCtrl            controller.ProjectController
	ProjectMemberCtrl      controller.ProjectMemberController
	ApiKeyCtrl             controller.ApiKeyController

	// Middleware
	AuthMiddleware    *middleware.AuthMiddleware
	OrgMiddleware     *middleware.OrganizationMiddleware
	ProjectMiddleware *middleware.ProjectMiddleware
}

func (r *routes) buildDependencies() *dependencies {
	// Repositories
	repo := repository.NewRegistry(r.db)

	// Iam Services
	authService := service.NewAuthService(repo.User)
	sessionService := service.NewSessionService(repo.User, repo.Session)
	userService := service.NewUserService(repo.User)

	// Organization services
	orgService := service.NewOrganizationService(repo.Organization, repo.User, repo.OrganizationMember, repo.AuditLog)
	orgMemberService := service.NewOrganizationMemberService(repo.OrganizationMember, repo.AuditLog, repo.OrganizationRBAC)
	orgInvitationService := service.NewOrganizationMemberInvitationService(
		repo.OrganizationMemberInvitation,
		repo.OrganizationMember,
		repo.User,
	)
	orgRBACService := service.NewOrganizationRBACService(repo.OrganizationRBAC)

	// Project services
	projectService := service.NewProjectService(repo.Project, repo.ProjectMember, repo.Organization, repo.AuditLog)
	projectMemberService := service.NewProjectMemberService(repo.ProjectMember, repo.Project, repo.AuditLog)
	apiKeyService := service.NewApiKeyService(repo.ApiKey, repo.ProjectMember, repo.Project)

	// Utility services
	tokenService := service.NewTokenService(r.cfg.Auth)

	// Controllers
	healthCtrl := controller.NewHealthController()
	authCtrl := controller.NewAuthController(authService, tokenService, sessionService)
	userCtrl := controller.NewUserController(userService)
	orgCtrl := controller.NewOrganizationController(orgService)
	orgMemCtrl := controller.NewOrganizationMemberController(orgMemberService)
	orgMemInvCtrl := controller.NewOrganizationMemberInvitationController(orgInvitationService)
	orgRBACCtrl := controller.NewOrganizationRBACController(orgRBACService)
	projectCtrl := controller.NewProjectController(projectService)
	projectMemberCtrl := controller.NewProjectMemberController(projectMemberService)
	apiKeyCtrl := controller.NewApiKeyController(apiKeyService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	orgMiddleware := middleware.NewOrganizationMiddleware(orgMemberService)
	projectMiddleware := middleware.NewProjectMiddleware(projectMemberService)

	return &dependencies{
		// Controllers
		HealthCtrl:             healthCtrl,
		AuthCtrl:               authCtrl,
		UserCtrl:               userCtrl,
		OrganizationCtrl:       orgCtrl,
		OrganizationMemCtrl:    orgMemCtrl,
		OrganizationMemInvCtrl: orgMemInvCtrl,
		OrganizationRBACCtrl:   orgRBACCtrl,
		ProjectCtrl:            projectCtrl,
		ProjectMemberCtrl:      projectMemberCtrl,
		ApiKeyCtrl:             apiKeyCtrl,
		// Middleware
		AuthMiddleware:    &authMiddleware,
		OrgMiddleware:     &orgMiddleware,
		ProjectMiddleware: &projectMiddleware,
	}
}
