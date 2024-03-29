package card

import (
	"context"
	"go-user-microservice/internal/app/user/services"
	"go-user-microservice/pkg/protobuf/core"
	"google.golang.org/grpc"
	"reflect"
)

type GrpcCardMiddleware struct {
	authUserContextService    *services.UserAuthContextService
	messageTypesAuthenticated []interface{}
}

func NewGrpcCardMiddleware(authUserContextService *services.UserAuthContextService) *GrpcCardMiddleware {
	messageTypesAuthenticated := []interface{}{
		&core.CreateCardRequest{},
		&core.MyCardsRequest{},
	}
	return &GrpcCardMiddleware{
		authUserContextService:    authUserContextService,
		messageTypesAuthenticated: messageTypesAuthenticated,
	}
}

func (m *GrpcCardMiddleware) CardsRequestAuthenticated(
	ctx context.Context,
	request interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if _, isCardServer := info.Server.(*GrpcServerCard); isCardServer {
		for _, messageType := range m.messageTypesAuthenticated {
			requestReflectType := reflect.TypeOf(request)
			messageReflectType := reflect.TypeOf(messageType)
			if requestReflectType == messageReflectType {
				newContext, e := m.authUserContextService.VerifyRetrieveNewUserContext(ctx)
				if e != nil {
					return nil, e
				}
				return handler(newContext, request)
			}
		}
		return handler(ctx, request)
	}
	return handler(ctx, request)
}
