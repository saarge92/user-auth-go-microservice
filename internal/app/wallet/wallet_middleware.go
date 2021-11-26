package wallet

import (
	"context"
	sharedServices "go-user-microservice/internal/app/user/services"
	"go-user-microservice/pkg/protobuf/wallet"
	"google.golang.org/grpc"
	"reflect"
)

type GrpcWalletMiddleware struct {
	authContextService        *sharedServices.UserAuthContextService
	messageTypesAuthenticated []interface{}
}

func NewWalletGrpcServerMiddleware(
	authContextService *sharedServices.UserAuthContextService,
) *GrpcWalletMiddleware {
	messageTypesAuthenticated := []interface{}{
		&wallet.CreateWalletRequest{},
		&wallet.MyWalletsRequest{},
	}
	return &GrpcWalletMiddleware{
		authContextService:        authContextService,
		messageTypesAuthenticated: messageTypesAuthenticated,
	}
}
func (m *GrpcWalletMiddleware) WalletsRequestsAuthenticated(
	ctx context.Context,
	request interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	for _, messageType := range m.messageTypesAuthenticated {
		requestReflectType := reflect.TypeOf(request)
		messageReflectType := reflect.TypeOf(messageType)
		if requestReflectType == messageReflectType {
			newContext, e := m.authContextService.VerifyUserFromRequest(ctx)
			if e != nil {
				return nil, e
			}
			return handler(newContext, request)
		}
	}
	return handler(ctx, request)
}
