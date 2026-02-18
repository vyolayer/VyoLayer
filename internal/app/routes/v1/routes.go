package routes

import (
	"worklayer/internal/app/controller"
	"worklayer/internal/app/middleware"
	"worklayer/internal/config"
	"worklayer/internal/repository"
	"worklayer/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Routes interface {
	SetupRoutes()
}

type routesV1 struct {
	router fiber.Router
	cfg    *config.Config
	db     *gorm.DB
}

func NewV1Routes(router fiber.Router, config *config.Config, db *gorm.DB) Routes {
	return &routesV1{
		router: router,
		cfg:    config,
		db:     db,
	}
}

func (router *routesV1) SetupRoutes() {
	// Initialize repositories
	repo := repository.NewRegistry(router.db)

	// Core services
	authService := service.NewAuthService(repo.User)
	sessionService := service.NewSessionService(repo.User, repo.Session)
	userService := service.NewUserService(repo.User)

	// Organization services
	orgService := service.NewOrganizationService(repo.Organization, repo.User)
	orgMemberService := service.NewOrganizationMemberService(repo.OrganizationMember)
	orgInvitationService := service.NewOrganizationMemberInvitationService(repo.OrganizationMemberInvitation, repo.OrganizationMember, repo.User)
	orgRBACService := service.NewOrganizationRBACService(repo.OrganizationRBAC)

	// Utility services
	tokenService := service.NewTokenService(router.cfg.Auth)

	// Controller
	authController := controller.NewAuthController(authService, tokenService, sessionService)
	userController := controller.NewUserController(userService)
	orgController := controller.NewOrganizationController(orgService)
	orgMemberController := controller.NewOrganizationMemberController(orgMemberService)
	orgInvitationController := controller.NewOrganizationMemberInvitationController(orgInvitationService)
	orgRBACController := controller.NewOrganizationRBACController(orgRBACService)

	// Middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService)
	orgMiddleware := middleware.NewOrganizationMiddleware(orgMemberService)

	// Register routes
	//
	// Auth routes
	authRouter := router.router.Group("/auth")
	authRouter.Post("/register", authController.RegisterUser)
	authRouter.Post("/login", authController.LoginUser)
	authRouter.Post("/refresh", authController.RefreshSession)

	authRouter.Use(authMiddleware.JwtValidated())
	authRouter.Post("/validate", authController.ValidateSession)
	authRouter.Post("/logout", authController.LogoutUser)

	// User routes
	userRouter := router.router.Group("/users")
	userRouter.Use(authMiddleware.JwtValidated())
	userRouter.Get("/me", userController.GetMe)

	// Organization routes
	orgRouter := router.router.Group("/organizations")
	orgRouter.Use(authMiddleware.JwtValidated())

	// Access by user
	orgRouter.Post("/", orgController.CreateOrganization)             // Create organization
	orgRouter.Get("/", orgController.ListOrganizations)               // List organizations
	orgRouter.Post("/onboarding", orgController.OnboardOrganization)  // Onboard organization
	orgRouter.Get("/slug/:slug", orgController.GetOrganizationBySlug) // Get organization by slug

	// Organization invitation routes (user)
	orgRouter.Get("/invitations/pending", orgInvitationController.GetPendingInvitations)
	orgRouter.Post("/invitations/accept", orgInvitationController.AcceptInvitation)
	orgRouter.Delete("/invitations/:invitationId", orgInvitationController.CancelInvitation)

	// Access by organization member (all members)
	orgRouter.Get("/:orgId", orgMiddleware.CheckOrganizationMembership(), orgController.GetOrganizationByID)
	orgRouter.Get("/:orgId/members/me", orgMiddleware.CheckOrganizationMembership(), orgMemberController.CurrentMember)

	// Access by organization owner and admin
	orgAdminRouter := orgRouter.Group("/:orgId")
	orgAdminRouter.Use(orgMiddleware.CheckOrganizationMembership(), orgMiddleware.IsAdmin())
	// Organization member routes (owner and admin)
	orgAdminRouter.Get("/members", orgMemberController.GetAllMembersByOrgID)
	orgAdminRouter.Get("/members/:memberId", orgMemberController.GetMemberByOrgIDAndMemberID)

	// Organization invitation routes (owner and admin)
	orgAdminRouter.Post("/invitations", orgInvitationController.CreateInvitation)
	orgAdminRouter.Get("/invitations", orgInvitationController.ListInvitations)

	// Organization rbac routes
	orgAdminRouter.Get("/rbac/permissions", orgRBACController.GetAllPermissions)
	orgAdminRouter.Get("/rbac/roles", orgRBACController.GetAllRoles)
}
