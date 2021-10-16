package middlewares

import (
	"context"
	"fmt"
	"go-user-microservice/internal/app/errorlists"
	"go-user-microservice/pkg/protobuf/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UserUnaryInterceptor(
	ctx context.Context,
	request interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	if message, ok := request.(*wallet.CreateWalletMessage); ok {
		fmt.Println(message)
		if headers, ok := metadata.FromIncomingContext(ctx); ok {
			token := headers["token"]
			fmt.Println(token)
		}
		return nil, status.Error(codes.Unauthenticated, errorlists.UserNotFound)
	}
	return handler(ctx, request)
}
