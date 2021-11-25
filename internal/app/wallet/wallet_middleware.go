package wallet

import (
	"context"
	sharedServices "go-user-microservice/internal/app/user/services"
	"go-user-microservice/pkg/protobuf/wallet"
	"google.golang.org/grpc"
)

type GrpcWalletMiddleware struct {
	authContextService *sharedServices.UserAuthContextService
}

func NewWalletGrpcServerMiddleware(
	authContextService *sharedServices.UserAuthContextService,
) *GrpcWalletMiddleware {
	return &GrpcWalletMiddleware{
		authContextService: authContextService,
	}
}

func (m *GrpcWalletMiddleware) CreateWalletAuthenticated(
	ctx context.Context,
	request interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if _, ok := request.(*wallet.CreateWalletRequest); ok {
		newContext, e := m.authContextService.VerifyUserFromRequest(ctx)
		if e != nil {
			return nil, e
		}
		return handler(newContext, request)
	}
	return handler(ctx, request)
}
