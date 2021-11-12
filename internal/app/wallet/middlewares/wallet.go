package middlewares

import (
	"context"
	sharedServices "go-user-microservice/internal/pkg/services"
	"go-user-microservice/pkg/protobuf/wallet"
	"google.golang.org/grpc"
)

type WalletGrpcMiddleware struct {
	authContextService *sharedServices.UserAuthContextService
}

func NewWalletGrpcServerMiddleware(
	authContextService *sharedServices.UserAuthContextService,
) *WalletGrpcMiddleware {
	return &WalletGrpcMiddleware{
		authContextService: authContextService,
	}
}

func (m *WalletGrpcMiddleware) CreateWalletAuthenticated(
	ctx context.Context,
	request interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if _, ok := request.(*wallet.CreateWalletMessage); ok {
		newContext, e := m.authContextService.VerifyUserFromRequest(ctx)
		if e != nil {
			return nil, e
		}
		return handler(newContext, request)
	}
	return handler(ctx, request)
}
