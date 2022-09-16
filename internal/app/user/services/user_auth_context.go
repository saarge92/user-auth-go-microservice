package services

import (
	"context"
	"go-user-microservice/internal/pkg/dictionary"
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type UserAuthContextService struct {
	jwtService *JwtService
}

func NewUserAuthContextService(jwtService *JwtService) *UserAuthContextService {
	return &UserAuthContextService{
		jwtService: jwtService,
	}
}

func (s *UserAuthContextService) VerifyRetrieveNewUserContext(ctx context.Context) (context.Context, error) {
	tokenInfo, e := s.retrieveTokenFromContext(ctx)
	if e != nil {
		return nil, e
	}
	userData, e := s.jwtService.VerifyTokenAndReturnUserRoleData(ctx, tokenInfo)
	if e != nil {
		return nil, e
	}
	if len(userData.Roles) == 0 {
		return nil, status.Error(codes.PermissionDenied, "User has no roles")
	}
	newCtx := context.WithValue(ctx, dictionary.CurrentUser, userData.User)
	return newCtx, nil
}

func (s *UserAuthContextService) retrieveTokenFromContext(ctx context.Context) (string, error) {
	authError := errorlists.ErrUserUnAuthenticated
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		var tokenInfo []string
		if tokenInfo, ok = headers["token"]; !ok {
			return "", authError
		}
		if len(tokenInfo) == 0 {
			return "", authError
		}
		return tokenInfo[0], nil
	}
	return "", authError
}
