package services

import (
	"github.com/dgrijalva/jwt-go"
	"go-user-microservice/internal/config"
	"go-user-microservice/internal/dto"
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
		Roles:    []string{"User"},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, e := jwtToken.SignedString([]byte(s.config.JwtKey))
	if e != nil {
		return "", e
	}
	return stringToken, nil
}
