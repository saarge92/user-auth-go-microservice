package member

import (
	"go-user-microservice/internal/app/dto"
	"go-user-microservice/internal/app/entites"
	"go-user-microservice/internal/app/forms/user"
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
	f *user.SignUp, chanResp chan<- interface{},
) (*entites.User, string, error) {
	userEntity, e := s.UserService.SignUp(f)
	defer close(chanResp)
	if e != nil {
		return nil, "", e
	}
	token, e := s.JwtService.CreateToken(userEntity.Login)
	if e != nil {
		return nil, "", e
	}
	return userEntity, token, nil
}

func (s *AuthService) SignIn(
	f *user.SignIn, chanResp chan<- interface{},
) (*entites.User, string, error) {
	userEntity, e := s.UserService.SignIn(f)
	defer close(chanResp)
	if e != nil {
		return nil, "", e
	}
	token, e := s.JwtService.CreateToken(userEntity.Login)
	if e != nil {
		return nil, "", e
	}
	return userEntity, token, nil
}

func (s *AuthService) VerifyAndReturnPayloadToken(token string) (*dto.UserPayLoad, *entites.User, error) {
	return s.JwtService.VerifyAndReturnPayloadToken(token)
}
