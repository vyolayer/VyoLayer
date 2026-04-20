package wire

import (
	"github.com/vyolayer/vyolayer/internal/gateway/handlers"
	"github.com/vyolayer/vyolayer/internal/gateway/handlers/console"
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
		handlers.NewHealthHandler(),

		// Account routes
		handlers.NewAccountHandler(
			clients.AccountClient,
			cookieSrv,
			accountJWT,
		),

		// IAM routes
		handlers.NewIAMAuthGatewayHandler(
			clients.IamAuthClient,
			clients.IamUserClient,
			iamCookieSrv,
			iamJWT,
		),

		// Tenant routes
		// Tenant Organization routes
		handlers.NewOrganizationHandler(
			logger,
			clients.TenantOrganizationClient,
			iamJWT,
		),

		// Tenant Organization Member routes
		handlers.NewOrganizationMemberHandler(
			logger,
			clients.TenantOrganizationMemClient,
			iamJWT,
		),

		// Tenant Organization Invitation routes
		handlers.NewOrganizationInvitationHandler(
			logger,
			clients.TenantOrganizationInvClient,
			iamJWT,
		),

		// Tenant Project & Project Member routes
		handlers.NewProjectHandler(
			logger,
			clients.TenantProjectClient,
			iamJWT,
		),

		console.NewProjectServiceHandler(
			logger,
			clients.ConsoleProjectServiceManifestClient,
			iamJWT,
		),
	}
}
