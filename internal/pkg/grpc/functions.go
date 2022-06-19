package grpc

import (
	"context"
	userDto "go-user-microservice/internal/app/user/dto"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetUserWithRolesFromContext(ctx context.Context) (*userDto.UserRole, error) {
	var user *userDto.UserRole
	var ok bool
	if user, ok = ctx.Value(dictionary.User).(*userDto.UserRole); !ok {
		return nil, status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	}

	return user, nil
}
