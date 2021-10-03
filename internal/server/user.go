package server

import (
	"context"
	"go-user-microservice/internal/entites"
	"go-user-microservice/internal/forms/builders"
	"go-user-microservice/internal/services"
	"go-user-microservice/pkg/protobuf/user"
)

type UserGrpcServer struct {
	authService     *services.AuthService
	userFormBuilder *builders.UserFormBuilder
}

func NewUserGrpcServer(
	userFormBuilder *builders.UserFormBuilder,
	authService *services.AuthService,
) *UserGrpcServer {
	return &UserGrpcServer{
		userFormBuilder: userFormBuilder,
		authService:     authService,
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
	var tokenResponse string
	var errorResponse error
	go func() {
		userResponse, tokenResponse, errorResponse = s.authService.SignUp(form, channelResponse)
	}()
	<-channelResponse
	if errorResponse != nil {
		return nil, errorResponse
	}
	return &user.SignUpResponse{
		Id:    userResponse.ID,
		Token: tokenResponse,
	}, nil
}

func (s *UserGrpcServer) VerifyToken(
	_ context.Context,
	request *user.VerifyMessage,
) (*user.VerifyMessageResponse, error) {
	userPayload, userEntity, e := s.authService.VerifyAndReturnPayloadToken(request.Token)
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
	var token string
	go func() {
		userResponse, token, signInError = s.authService.SignIn(form, signInChan)
	}()
	<-signInChan
	if signInError != nil {
		return nil, signInError
	}
	return &user.SignInResponse{
		Id:    uint64(userResponse.ID),
		Token: token,
	}, nil
}
