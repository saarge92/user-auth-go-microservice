package user

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-user-microservice/internal/config"
	"go-user-microservice/internal/contracts/repositories"
	"go-user-microservice/internal/dto"
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type JwtService struct {
	config         *config.Config
	userRepository repositories.UserRepositoryInterface
}

func NewJwtService(
	config *config.Config,
	userRepository repositories.UserRepositoryInterface,
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

func (s *JwtService) VerifyAndReturnPayloadToken(token string) (*dto.UserPayLoad, *entites.User, error) {
	jwtToken, e := jwt.ParseWithClaims(token, &dto.UserPayLoad{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JwtKey), nil
	})
	if e != nil {
		return nil, nil, status.Error(codes.InvalidArgument, errorlists.TokenInvalid)
	}
	payloadClaims := jwtToken.Claims.(*dto.UserPayLoad)
	user, e := s.checkClaims(payloadClaims)
	if e != nil {
		return nil, nil, e
	}

	return payloadClaims, user, nil
}

func (s *JwtService) checkClaims(claims *dto.UserPayLoad) (*entites.User, error) {
	tokenInvalidError := status.Error(codes.InvalidArgument, errorlists.TokenInvalid)
	login := claims.UserName
	user, e := s.userRepository.GetUser(login)
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
