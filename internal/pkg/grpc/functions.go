package grpc

import (
	"context"
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetUserFromContext(ctx context.Context) (*entities.User, error) {
	var user *entities.User
	var ok bool
	if user, ok = ctx.Value(dictionary.User).(*entities.User); !ok {
		return nil, status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	}

	return user, nil
}
