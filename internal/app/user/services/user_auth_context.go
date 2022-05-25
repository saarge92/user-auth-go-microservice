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

func (s *UserAuthContextService) VerifyUserFromRequest(ctx context.Context) (context.Context, error) {
	authError := status.Error(codes.Unauthenticated, errorlists.UserUnAuthenticated)
	if headers, ok := metadata.FromIncomingContext(ctx); ok {
		var tokenInfo []string
		if tokenInfo, ok = headers["token"]; !ok {
			return nil, authError
		}
		if len(tokenInfo) == 0 {
			return nil, authError
		}
		userData, e := s.jwtService.VerifyTokenAndReturnUser(ctx, tokenInfo[0])
		if e != nil {
			return nil, e
		}
		newCtx := context.WithValue(ctx, dictionary.User, userData)
		return newCtx, nil
	}
	return nil, authError
}
