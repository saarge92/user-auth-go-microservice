package server

import (
	"context"
	entites2 "go-user-microservice/internal/app/entites"
	builders2 "go-user-microservice/internal/app/forms/builders"
	user2 "go-user-microservice/internal/app/services/member"
	"go-user-microservice/pkg/protobuf/user"
)

type UserGrpcServer struct {
	authService     *user2.AuthService
	userFormBuilder *builders2.UserFormBuilder
}

func NewUserGrpcServer(
	userFormBuilder *builders2.UserFormBuilder,
	authService *user2.AuthService,
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
	var userResponse *entites2.User
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
			Id:    userEntity.ID,
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
	var userResponse *entites2.User
	var token string
	go func() {
		userResponse, token, signInError = s.authService.SignIn(form, signInChan)
	}()
	<-signInChan
	if signInError != nil {
		return nil, signInError
	}
	return &user.SignInResponse{
		Id:    userResponse.ID,
		Token: token,
	}, nil
}