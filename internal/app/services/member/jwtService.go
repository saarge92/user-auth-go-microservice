package member

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	config2 "go-user-microservice/internal/app/config"
	repositories2 "go-user-microservice/internal/app/domain/repositories"
	dto2 "go-user-microservice/internal/app/dto"
	entites2 "go-user-microservice/internal/app/entites"
	errorlists2 "go-user-microservice/internal/app/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type JwtService struct {
	config         *config2.Config
	userRepository repositories2.UserRepositoryInterface
}

func NewJwtService(
	config *config2.Config,
	userRepository repositories2.UserRepositoryInterface,
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
	claims := &dto2.UserPayLoad{
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

func (s *JwtService) VerifyAndReturnPayloadToken(token string) (*dto2.UserPayLoad, *entites2.User, error) {
	jwtToken, e := jwt.ParseWithClaims(token, &dto2.UserPayLoad{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JwtKey), nil
	})
	if e != nil {
		return nil, nil, status.Error(codes.InvalidArgument, errorlists2.TokenInvalid)
	}
	payloadClaims := jwtToken.Claims.(*dto2.UserPayLoad)
	user, e := s.checkClaims(payloadClaims)
	if e != nil {
		return nil, nil, e
	}

	return payloadClaims, user, nil
}

func (s *JwtService) checkClaims(claims *dto2.UserPayLoad) (*entites2.User, error) {
	tokenInvalidError := status.Error(codes.InvalidArgument, errorlists2.TokenInvalid)
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
