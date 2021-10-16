package middlewares

import (
	"context"
	"go-user-microservice/internal/app/errorlists"
	"go-user-microservice/internal/app/services/member"
	"go-user-microservice/pkg/protobuf/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserGrpcMiddleware struct {
	authService *member.AuthService
}

func NewUserGrpcMiddleware(
	authService *member.AuthService,
) *UserGrpcMiddleware {
	return &UserGrpcMiddleware{
		authService: authService,
	}
}

func (m *UserGrpcMiddleware) IsAuthenticatedMiddleware(
	ctx context.Context,
	request interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := request.(*wallet.CreateWalletMessage); ok {
		if headers, ok := metadata.FromIncomingContext(ctx); ok {
			var tokenInfo []string
			if tokenInfo, ok = headers["token"]; !ok {
				return nil, status.Error(codes.Unauthenticated, errorlists.UserNotFound)
			}
			if len(tokenInfo) == 0 {
				return nil, status.Error(codes.Unauthenticated, errorlists.UserNotFound)
			}
			_, user, e := m.authService.VerifyAndReturnPayloadToken(tokenInfo[0])
			if e != nil {
				return nil, e
			}
			newCtx := context.WithValue(ctx, "user_id", user.ID)
			return handler(newCtx, request)
		}
		return nil, status.Error(codes.Unauthenticated, errorlists.UserNotFound)
	}
	return handler(ctx, request)
}
