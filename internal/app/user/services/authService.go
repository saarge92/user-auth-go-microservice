package services

import (
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/app/user/forms"
)

type AuthService struct {
	UserService *ServiceUser
	jwtService  *JwtService
}

func NewAuthService(
	userService *ServiceUser,
	jwtService *JwtService,
) *AuthService {
	return &AuthService{
		UserService: userService,
		jwtService:  jwtService,
	}
}

func (s *AuthService) SignUp(
	f *forms.SignUp, chanResp chan<- interface{},
) (*entities.User, string, error) {
	userEntity, e := s.UserService.SignUp(f)
	defer close(chanResp)
	if e != nil {
		return nil, "", e
	}

	token, e := s.jwtService.CreateToken(userEntity.Login)
	if e != nil {
		return nil, "", e
	}
	return userEntity, token, nil
}

func (s *AuthService) SignIn(
	f *forms.SignIn, chanResp chan<- interface{},
) (*entities.User, string, error) {
	userEntity, e := s.UserService.SignIn(f)
	defer close(chanResp)
	if e != nil {
		return nil, "", e
	}
	token, e := s.jwtService.CreateToken(userEntity.Login)
	if e != nil {
		return nil, "", e
	}
	return userEntity, token, nil
}

func (s *AuthService) VerifyAndReturnPayloadToken(token string) (*entities.User, error) {
	return s.jwtService.VerifyTokenAndReturnUser(token)
}
