package card

import (
	"context"
	"go-user-microservice/internal/app/user/services"
	"go-user-microservice/pkg/protobuf/card"
	"google.golang.org/grpc"
)

type GrpcCardMiddleware struct {
	authUserContextService *services.UserAuthContextService
}

func NewGrpcCardMiddleware(authUserContextService *services.UserAuthContextService) *GrpcCardMiddleware {
	return &GrpcCardMiddleware{
		authUserContextService: authUserContextService,
	}
}

func (m *GrpcCardMiddleware) CreateCardAuthenticated(
	ctx context.Context,
	request interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := request.(*card.CreateCardRequest); ok {
		newContext, e := m.authUserContextService.VerifyUserFromRequest(ctx)
		if e != nil {
			return nil, e
		}
		return handler(newContext, request)
	}
	return handler(ctx, request)
}
