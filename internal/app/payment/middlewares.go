package payment

import (
	"context"
	sharedServices "go-user-microservice/internal/app/user/services"
	"go-user-microservice/pkg/protobuf/core"
	"google.golang.org/grpc"
	"reflect"
)

type GrpcPaymentMiddleware struct {
	authContextService        *sharedServices.UserAuthContextService
	messageTypesAuthenticated []interface{}
}

func NewGrpcPaymentMiddleware(
	authContextService *sharedServices.UserAuthContextService,
) *GrpcPaymentMiddleware {
	messagesTypes := []interface{}{
		&core.DepositRequest{},
		&core.ListRequest{},
	}
	return &GrpcPaymentMiddleware{
		authContextService:        authContextService,
		messageTypesAuthenticated: messagesTypes,
	}
}

func (m *GrpcPaymentMiddleware) PaymentsRequestsAuthenticated(
	ctx context.Context,
	request interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if _, isCardServer := info.Server.(*GrpcServerPayment); isCardServer {
		for _, messageType := range m.messageTypesAuthenticated {
			requestReflectType := reflect.TypeOf(request)
			messageReflectType := reflect.TypeOf(messageType)
			if requestReflectType == messageReflectType {
				newContext, e := m.authContextService.VerifyRetrieveNewUserContext(ctx)
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
