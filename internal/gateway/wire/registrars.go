package wire

import (
	"github.com/vyolayer/vyolayer/internal/gateway/handlers"
	"github.com/vyolayer/vyolayer/internal/gateway/server"
	"github.com/vyolayer/vyolayer/internal/gateway/service"
	"github.com/vyolayer/vyolayer/pkg/jwt"
)

func NewRegistrars(
	clients *Clients,
	cookieSrv *service.AccountTokenService,
	accountJWT jwt.AccountJWT,
	iamCookieSrv *service.IAMCookieService,
	iamJWT jwt.IamJWT,
) []server.RouteRegistrar {
	return []server.RouteRegistrar{
		handlers.NewAccountHandler(clients.AccountClient, cookieSrv, accountJWT),
		handlers.NewIAMAuthGatewayHandler(clients.IamAuthClient, clients.IamUserClient, iamCookieSrv, iamJWT),
	}
}
