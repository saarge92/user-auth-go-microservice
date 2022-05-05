package services

import (
	"go-user-microservice/internal/app/user/entities"
	"go-user-microservice/internal/app/user/forms"
)

type Auth struct {
	UserService *User
	jwtService  *JwtService
}

func NewAuthService(
	userService *User,
	jwtService *JwtService,
) *Auth {
	return &Auth{
		UserService: userService,
		jwtService:  jwtService,
	}
}

func (s *Auth) SignUp(
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

func (s *Auth) SignIn(
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

func (s *Auth) VerifyAndReturnPayloadToken(token string) (*entities.User, error) {
	return s.jwtService.VerifyTokenAndReturnUser(token)
}
