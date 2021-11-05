package middlewares

import (
	"context"
	"go-user-microservice/internal/app/user/services"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/errorlists"
	"go-user-microservice/pkg/protobuf/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type WalletGrpcMiddleware struct {
	jwtService *services.JwtService
}

func NewUserGrpcMiddleware(
	jwtService *services.JwtService,
) *WalletGrpcMiddleware {
	return &WalletGrpcMiddleware{
		jwtService: jwtService,
	}
}

func (m *WalletGrpcMiddleware) IsAuthenticatedMiddleware(
	ctx context.Context,
	request interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if _, ok := request.(*wallet.CreateWalletMessage); ok {
		authError := status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
		if headers, ok := metadata.FromIncomingContext(ctx); ok {
			var tokenInfo []string
			if tokenInfo, ok = headers["token"]; !ok {
				return nil, authError
			}
			if len(tokenInfo) == 0 {
				return nil, authError
			}
			userData, e := m.jwtService.VerifyTokenAndReturnUser(tokenInfo[0])
			if e != nil {
				return nil, e
			}
			newCtx := context.WithValue(ctx, dictionary.User, userData)
			return handler(newCtx, request)
		}
		return nil, authError
	}
	return handler(ctx, request)
}
