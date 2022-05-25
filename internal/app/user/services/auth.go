package services

import (
	"context"
	"go-user-microservice/internal/app/user/dto"
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
	ctx context.Context,
	formRequest *forms.SignUp,
	signUpResponseChannel chan<- interface{},
) (*entities.User, string, error) {
	userEntity, e := s.UserService.SignUp(ctx, formRequest)
	defer close(signUpResponseChannel)
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
	ctx context.Context,
	formRequest *forms.SignIn,
	signInResponseChannel chan<- interface{},
) (*dto.UserRole, string, error) {
	defer close(signInResponseChannel)
	userEntity, e := s.UserService.SignIn(ctx, formRequest)
	if e != nil {
		return nil, "", e
	}
	token, e := s.jwtService.CreateToken(userEntity.User.Login)
	if e != nil {
		return nil, "", e
	}
	return userEntity, token, nil
}

func (s *Auth) VerifyAndReturnPayloadToken(ctx context.Context, token string) (*dto.UserRole, error) {
	return s.jwtService.VerifyTokenAndReturnUser(ctx, token)
}
