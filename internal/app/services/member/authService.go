package member

import (
	dto2 "go-user-microservice/internal/app/dto"
	entites2 "go-user-microservice/internal/app/entites"
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
) (*entites2.User, string, error) {
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
	f *user.SignIn, chanResp chan<- interface{},
) (*entites2.User, string, error) {
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

func (s *AuthService) VerifyAndReturnPayloadToken(token string) (*dto2.UserPayLoad, *entites2.User, error) {
	return s.JwtService.VerifyAndReturnPayloadToken(token)
}
