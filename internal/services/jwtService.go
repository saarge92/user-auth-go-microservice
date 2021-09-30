package services

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"go-user-microservice/internal/config"
	"go-user-microservice/internal/dto"
	"go-user-microservice/internal/entites/dictionary"
	"go-user-microservice/internal/errorlists"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type JwtService struct {
	config *config.Config
}

func NewJwtService(config *config.Config) *JwtService {
	return &JwtService{config: config}
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
		Roles:    []string{string(dictionary.UserRole)},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, e := jwtToken.SignedString([]byte(s.config.JwtKey))
	if e != nil {
		return "", e
	}
	return stringToken, nil
}

func (s *JwtService) VerifyAndReturnPayloadToken(token string) (*dto.UserPayLoad, error) {
	jwtToken, e := jwt.ParseWithClaims(token, &dto.UserPayLoad{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.config.JwtKey), nil
	})
	if e != nil {
		return nil, status.Error(codes.InvalidArgument, errorlists.UserNotFound)
	}
	payloadClaims := jwtToken.Claims.(*dto.UserPayLoad)
	return payloadClaims, nil
}
