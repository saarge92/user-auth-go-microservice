package middlewares

import (
	"context"
	"go-user-microservice/internal/app/errorlists"
	"go-user-microservice/internal/app/middlewares/dictionary"
	"go-user-microservice/internal/app/services/member"
	"go-user-microservice/pkg/protobuf/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserGrpcMiddleware struct {
	jwtService *member.JwtService
}

func NewUserGrpcMiddleware(
	jwtService *member.JwtService,
) *UserGrpcMiddleware {
	return &UserGrpcMiddleware{
		jwtService: jwtService,
	}
}

func (m *UserGrpcMiddleware) IsAuthenticatedMiddleware(
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
			_, user, e := m.jwtService.VerifyAndReturnPayloadToken(tokenInfo[0])
			if e != nil {
				return nil, e
			}
			newCtx := context.WithValue(ctx, dictionary.UserID, user.ID)
			return handler(newCtx, request)
		}
		return nil, authError
	}
	return handler(ctx, request)
}
