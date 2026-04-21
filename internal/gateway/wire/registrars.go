package wire

import (
	accounthandler "github.com/vyolayer/vyolayer/internal/gateway/handlers/account"
	consolehandler "github.com/vyolayer/vyolayer/internal/gateway/handlers/console"
	healthhandler "github.com/vyolayer/vyolayer/internal/gateway/handlers/health"
	iamhandler "github.com/vyolayer/vyolayer/internal/gateway/handlers/iam"
	tenanthandler "github.com/vyolayer/vyolayer/internal/gateway/handlers/tenant"
	"github.com/vyolayer/vyolayer/internal/gateway/server"
	"github.com/vyolayer/vyolayer/internal/gateway/service"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
)

func NewRegistrars(
	logger *logger.AppLogger,
	clients *Clients,
	cookieSrv *service.AccountTokenService,
	accountJWT jwt.AccountJWT,
	iamCookieSrv *service.IAMCookieService,
	iamJWT jwt.IamJWT,
) []server.RouteRegistrar {
	return []server.RouteRegistrar{
		healthhandler.NewHealthHandler(),

		// Account routes
		accounthandler.NewAccountHandler(
			clients.AccountClient,
			cookieSrv,
			accountJWT,
			logger,
		),

		// IAM routes
		iamhandler.NewIAMAuthGatewayHandler(
			clients.IamAuthClient,
			clients.IamUserClient,
			iamCookieSrv,
			iamJWT,
			logger,
		),

		// Tenant routes
		// Tenant Organization routes
		tenanthandler.NewOrganizationHandler(
			logger,
			clients.TenantOrganizationClient,
			iamJWT,
		),

		// Tenant Organization Member routes
		tenanthandler.NewOrganizationMemberHandler(
			logger,
			clients.TenantOrganizationMemClient,
			iamJWT,
		),

		// Tenant Organization Invitation routes
		tenanthandler.NewOrganizationInvitationHandler(
			logger,
			clients.TenantOrganizationInvClient,
			iamJWT,
		),

		// Tenant Project & Project Member routes
		tenanthandler.NewProjectHandler(
			logger,
			clients.TenantProjectClient,
			iamJWT,
		),

		consolehandler.NewProjectServiceHandler(
			logger,
			clients.ConsoleProjectServiceManifestClient,
			iamJWT,
		),
	}
}
