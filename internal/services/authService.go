package services

import (
	"go-user-microservice/internal/dto"
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/forms"
)

type AuthService struct {
	UserService *UserService
	JwtService  *JwtService
}

func NewAuthService(
	userService *UserService,
	jwtService *JwtService,
) *AuthService {
	return &AuthService{
		UserService: userService,
		JwtService:  jwtService,
	}
}

func (s *AuthService) SignUp(
	f *forms.SignUp, chanResp chan<- interface{},
) (*entites.User, string, error) {
	user, e := s.UserService.SignUp(f)
	defer close(chanResp)
	if e != nil {
		return nil, "", e
	}
	token, e := s.JwtService.CreateToken(user.Login)
	if e != nil {
		return nil, "", e
	}
	return user, token, nil
}

func (s *AuthService) SignIn(
	f *forms.SignIn, chanResp chan<- interface{},
) (*entites.User, string, error) {
	user, e := s.UserService.SignIn(f)
	defer close(chanResp)
	if e != nil {
		return nil, "", e
	}
	token, e := s.JwtService.CreateToken(user.Login)
	if e != nil {
		return nil, "", e
	}
	return user, token, nil
}

func (s *AuthService) VerifyAndReturnPayloadToken(token string) (*dto.UserPayLoad, *entites.User, error) {
	return s.JwtService.VerifyAndReturnPayloadToken(token)
}
