package services

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	sharedRepoInterfaces "go-user-microservice/internal/app/user/domain"
	"go-user-microservice/internal/app/user/dto"
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

func (s *JwtService) VerifyTokenAndReturnUser(ctx context.Context, token string) (*dto.UserRole, error) {
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
	login := claims.UserName
	user, e := s.userRepository.GetUserWithRoles(ctx, login)
	if e != nil {
		return nil, e
	}
	if user == nil {
		return nil, tokenInvalidError
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
