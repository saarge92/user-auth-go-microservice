package services

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	sharedRepoInterfaces "go-user-microservice/internal/app/user/domain"
	"go-user-microservice/internal/app/user/dto"
	userErrors "go-user-microservice/internal/app/user/errors"
	"go-user-microservice/internal/pkg/config"
	"go-user-microservice/internal/pkg/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type JwtService struct {
	config         *config.Config
	userRepository sharedRepoInterfaces.UserRepository
}

func NewJwtService(
	config *config.Config,
	userRepository sharedRepoInterfaces.UserRepository,
) *JwtService {
	return &JwtService{
		config:         config,
		userRepository: userRepository,
	}
}

func (s *JwtService) CreateToken(userName string) (string, error) {
	jwtExpirationTime := time.Duration(s.config.JwtExpiration)
	issuedAt := time.Now().UTC()
	expiredAt := issuedAt.Add(time.Minute * jwtExpirationTime)
	claims := &dto.UserPayLoad{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  issuedAt.Unix(),
			ExpiresAt: expiredAt.Unix(),
			Issuer:    s.config.JwtAudience,
		},
		UserName: userName,
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, e := jwtToken.SignedString([]byte(s.config.JwtKey))
	if e != nil {
		return "", e
	}
	return stringToken, nil
}

func (s *JwtService) VerifyTokenAndReturnUserRoleData(ctx context.Context, token string) (*dto.UserRole, error) {
	jwtToken, e := jwt.ParseWithClaims(token, &dto.UserPayLoad{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JwtKey), nil
	})
	if e != nil {
		return nil, status.Error(codes.InvalidArgument, errorlists.TokenInvalid)
	}
	payloadClaims := jwtToken.Claims.(*dto.UserPayLoad)
	user, e := s.checkClaims(ctx, payloadClaims)
	if e != nil {
		return nil, e
	}

	return user, nil
}

func (s *JwtService) checkClaims(ctx context.Context, claims *dto.UserPayLoad) (*dto.UserRole, error) {
	tokenInvalidError := status.Error(codes.InvalidArgument, errorlists.TokenInvalid)
	user, e := s.userRepository.GetUserWithRoles(ctx, claims.UserName)
	if e != nil {
		if errors.Is(e, userErrors.ErrUserNotFound) {
			return nil, tokenInvalidError
		}
		return nil, e
	}

	now := time.Now()
	if claims.ExpiresAt < now.Unix() {
		return nil, tokenInvalidError
	}
	if claims.Issuer != s.config.JwtAudience {
		return nil, tokenInvalidError
	}

	return user, nil
}
