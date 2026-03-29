package wire

import (
	"log"
	"time"

	"github.com/vyolayer/vyolayer/internal/gateway/config"
	"github.com/vyolayer/vyolayer/pkg/grpcutil"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
	iamV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
	"google.golang.org/grpc"
)

// Clients holds all gRPC client connections for the gateway.
type Clients struct {
	AccountClient  accountV1.AccountServiceClient
	IamAuthClient  iamV1.AuthServiceClient
	IamUserClient  iamV1.UserServiceClient

	// Keep references to close them later
	accountConn *grpc.ClientConn
	iamConn     *grpc.ClientConn
}

func NewClients(cfg *config.Config) (*Clients, error) {
	// Connect to Account Service
	accountConn, err := grpcutil.NewClient(grpcutil.ClientConfig{
		Address: cfg.AccountServiceAddr,
		Timeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	// Connect to IAM Service
	iamConn, err := grpcutil.NewClient(grpcutil.ClientConfig{
		Address: cfg.IAMServiceAddr,
		Timeout: 5 * time.Second,
	})
	if err != nil {
		accountConn.Close()
		return nil, err
	}

	return &Clients{
		AccountClient: accountV1.NewAccountServiceClient(accountConn),
		IamAuthClient: iamV1.NewAuthServiceClient(iamConn),
		IamUserClient: iamV1.NewUserServiceClient(iamConn),

		accountConn: accountConn,
		iamConn:     iamConn,
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
}
