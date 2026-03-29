package grpc

import (
	"github.com/vyolayer/vyolayer/internal/iam/usecase"
	iAMV1 "github.com/vyolayer/vyolayer/proto/iam/v1"
)

// IAMAuthHandler implements iAMV1.AuthServiceServer.
type IAMAuthHandler struct {
	iAMV1.UnimplementedAuthServiceServer
	au *usecase.AuthUsecase
}

// IAMUserHandler implements iAMV1.UserServiceServer.
type IAMUserHandler struct {
	iAMV1.UnimplementedUserServiceServer
	uu *usecase.UserUsecase
}
