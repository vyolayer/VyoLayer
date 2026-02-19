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

	// Middleware
	AuthMiddleware *middleware.AuthMiddleware
	OrgMiddleware  *middleware.OrganizationMiddleware
}

func (r *routes) buildDependencies() *dependencies {
	// Repositories
	repo := repository.NewRegistry(r.db)

	// Iam Services
	authService := service.NewAuthService(repo.User)
	sessionService := service.NewSessionService(repo.User, repo.Session)
	userService := service.NewUserService(repo.User)

	// Organization services
	orgService := service.NewOrganizationService(repo.Organization, repo.User)
	orgMemberService := service.NewOrganizationMemberService(repo.OrganizationMember)
	orgInvitationService := service.NewOrganizationMemberInvitationService(
		repo.OrganizationMemberInvitation,
		repo.OrganizationMember,
		repo.User,
	)
	orgRBACService := service.NewOrganizationRBACService(repo.OrganizationRBAC)

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

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	orgMiddleware := middleware.NewOrganizationMiddleware(orgMemberService)

	return &dependencies{
		// Controllers
		HealthCtrl:             healthCtrl,
		AuthCtrl:               authCtrl,
		UserCtrl:               userCtrl,
		OrganizationCtrl:       orgCtrl,
		OrganizationMemCtrl:    orgMemCtrl,
		OrganizationMemInvCtrl: orgMemInvCtrl,
		OrganizationRBACCtrl:   orgRBACCtrl,
		// Middleware
		AuthMiddleware: &authMiddleware,
		OrgMiddleware:  &orgMiddleware,
	}
}
