package server

import (
	"context"
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/forms/builders"
	"go-user-microservice/internal/services"
	"go-user-microservice/pkg/protobuf/user"
)

type UserGrpcServer struct {
	userService     *services.UserService
	jwtService      *services.JwtService
	userFormBuilder *builders.UserFormBuilder
}

func NewUserGrpcServer(
	userService *services.UserService,
	userFormBuilder *builders.UserFormBuilder,
	jwtService *services.JwtService,
) *UserGrpcServer {
	return &UserGrpcServer{
		userService:     userService,
		userFormBuilder: userFormBuilder,
		jwtService:      jwtService,
	}
}

func (s *UserGrpcServer) Signup(
	_ context.Context,
	request *user.SignUpMessage,
) (*user.SignUpResponse, error) {
	form := s.userFormBuilder.Signup(request)
	if e := form.Validate(); e != nil {
		return nil, e
	}
	channelResponse := make(chan interface{})
	var userResponse *entites.User
	var errorResponse error
	go func() {
		userResponse, errorResponse = s.userService.SignUp(form, channelResponse)
	}()
	<-channelResponse
	if errorResponse != nil {
		return nil, errorResponse
	}
	token, e := s.jwtService.CreateToken(userResponse.Login)
	if e != nil {
		return nil, e
	}
	return &user.SignUpResponse{
		Id:    userResponse.ID,
		Token: token,
	}, nil
}

func (s *UserGrpcServer) VerifyToken(
	_ context.Context,
	request *user.VerifyMessage,
) (*user.VerifyMessageResponse, error) {
	userPayload, userEntity, e := s.jwtService.VerifyAndReturnPayloadToken(request.Token)
	if e != nil {
		return nil, e
	}
	return &user.VerifyMessageResponse{
		User: &user.UserMessageResponse{
			Login: userPayload.UserName,
			Id:    uint64(userEntity.ID),
			Roles: nil,
		},
	}, nil
}

func (s *UserGrpcServer) SignIn(
	_ context.Context,
	request *user.SignInMessage,
) (*user.SignInResponse, error) {
	form := s.userFormBuilder.SignIn(request)
	signInChan := make(chan interface{})
	var signInError error
	var userResponse *entites.User
	go func() {
		userResponse, signInError = s.userService.SignIn(form, signInChan)
	}()
	<-signInChan
	if signInError != nil {
		return nil, signInError
	}
	token, e := s.jwtService.CreateToken(userResponse.Login)
	if e != nil {
		return nil, e
	}
	return &user.SignInResponse{
		Id:    uint64(userResponse.ID),
		Token: token,
	}, nil
}
