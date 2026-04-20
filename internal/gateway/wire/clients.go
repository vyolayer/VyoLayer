package wire

import (
	"log"
	"time"

	"github.com/vyolayer/vyolayer/internal/gateway/config"
	"github.com/vyolayer/vyolayer/pkg/grpcutil"
	"github.com/vyolayer/vyolayer/pkg/logger"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
	consoleV1 "github.com/vyolayer/vyolayer/proto/console/v1"
	iamV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
	tenantV1 "github.com/vyolayer/vyolayer/proto/tenant/v1"
	"google.golang.org/grpc"
)

// Clients holds all gRPC client connections for the gateway.
type Clients struct {
	AccountClient accountV1.AccountServiceClient
	IamAuthClient iamV1.AuthServiceClient
	IamUserClient iamV1.UserServiceClient

	// Console
	ConsoleProjectServiceManifestClient consoleV1.ProjectServiceManifestClient

	// Tenant
	TenantOrganizationClient    tenantV1.OrganizationServiceClient
	TenantOrganizationInvClient tenantV1.OrganizationInvitationServiceClient
	TenantOrganizationMemClient tenantV1.OrganizationMemberServiceClient
	TenantProjectClient         tenantV1.ProjectServiceClient

	// Keep references to close them later
	accountConn *grpc.ClientConn
	iamConn     *grpc.ClientConn
	tenantConn  *grpc.ClientConn
	consoleConn *grpc.ClientConn
}

func NewClients(logger *logger.AppLogger, cfg *config.Config, grpcTimeout time.Duration) (*Clients, error) {
	logger = logger.WithContext("GRPC Clients")

	// Connect to Account Service
	accountConn, err := grpcutil.NewClient(grpcutil.ClientConfig{
		Address: cfg.AccountServiceAddr,
		Timeout: grpcTimeout,
	})
	if err != nil {
		logger.Warn("Failed to connect to account service", err)
		return nil, err
	}

	// Connect to IAM Service
	iamConn, err := grpcutil.NewClient(grpcutil.ClientConfig{
		Address: cfg.IAMServiceAddr,
		Timeout: grpcTimeout,
	})
	if err != nil {
		logger.Warn("Failed to connect to iam service", err)
		return nil, err
	}

	// Connent to tenant service
	tenantConn, tenanterr := grpcutil.NewClient(grpcutil.ClientConfig{
		Address: cfg.TenantServiceAddr,
		Timeout: grpcTimeout,
	})
	if tenanterr != nil {
		logger.Warn("Failed to connect to tenant service", tenanterr)
		return nil, tenanterr
	}

	// Connent to console service
	consoleConn, consoleerr := grpcutil.NewClient(grpcutil.ClientConfig{
		Address: cfg.ConsoleServiceAddr,
		Timeout: grpcTimeout,
	})
	if consoleerr != nil {
		logger.Warn("Failed to connect to console service", consoleerr)
		return nil, consoleerr
	}

	return &Clients{
		AccountClient:                       accountV1.NewAccountServiceClient(accountConn),
		IamAuthClient:                       iamV1.NewAuthServiceClient(iamConn),
		IamUserClient:                       iamV1.NewUserServiceClient(iamConn),
		TenantOrganizationClient:            tenantV1.NewOrganizationServiceClient(tenantConn),
		TenantOrganizationInvClient:         tenantV1.NewOrganizationInvitationServiceClient(tenantConn),
		TenantOrganizationMemClient:         tenantV1.NewOrganizationMemberServiceClient(tenantConn),
		TenantProjectClient:                 tenantV1.NewProjectServiceClient(tenantConn),
		ConsoleProjectServiceManifestClient: consoleV1.NewProjectServiceManifestClient(consoleConn),

		accountConn: accountConn,
		iamConn:     iamConn,
		tenantConn:  tenantConn,
		consoleConn: consoleConn,
	}, nil
}

func (c *Clients) Close() {
	if c.accountConn != nil {
		if err := c.accountConn.Close(); err != nil {
			log.Printf("error closing account connection: %v", err)
		}
	}

	if c.iamConn != nil {
		if err := c.iamConn.Close(); err != nil {
			log.Printf("error closing iam connection: %v", err)
		}
	}

	if c.tenantConn != nil {
		if err := c.tenantConn.Close(); err != nil {
			log.Printf("error closing tenant connection: %v", err)
		}
	}
}
