package member

import (
	dto2 "go-user-microservice/internal/app/dto"
	entites2 "go-user-microservice/internal/app/entites"
	forms2 "go-user-microservice/internal/app/forms"
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
	f *forms2.SignUp, chanResp chan<- interface{},
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
	f *forms2.SignIn, chanResp chan<- interface{},
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
