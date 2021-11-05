package user

import (
	"context"
	"go-user-microservice/internal/app/user/entities"
	userFormBuilders "go-user-microservice/internal/app/user/forms/builders"
	"go-user-microservice/internal/app/user/services"
	"go-user-microservice/pkg/protobuf/user_server"
)

type GrpcUserServer struct {
	authService     *services.AuthService
	userFormBuilder *userFormBuilders.UserFormBuilder
}

func NewUserGrpcServer(
	userFormBuilder *userFormBuilders.UserFormBuilder,
	authService *services.AuthService,
) *GrpcUserServer {
	return &GrpcUserServer{
		userFormBuilder: userFormBuilder,
		authService:     authService,
	}
}

func (s *GrpcUserServer) Signup(
	_ context.Context,
	request *user_server.SignUpMessage,
) (*user_server.SignUpResponse, error) {
	form := s.userFormBuilder.Signup(request)
	if e := form.Validate(); e != nil {
		return nil, e
	}
	channelResponse := make(chan interface{})
	var userResponse *entities.User
	var tokenResponse string
	var errorResponse error
	go func() {
		userResponse, tokenResponse, errorResponse = s.authService.SignUp(form, channelResponse)
	}()
	<-channelResponse
	if errorResponse != nil {
		return nil, errorResponse
	}
	return &user_server.SignUpResponse{
		Id:    userResponse.ID,
		Token: tokenResponse,
	}, nil
}

func (s *GrpcUserServer) VerifyToken(
	_ context.Context,
	request *user_server.VerifyMessage,
) (*user_server.VerifyMessageResponse, error) {
	userEntity, e := s.authService.VerifyAndReturnPayloadToken(request.Token)
	if e != nil {
		return nil, e
	}
	return &user_server.VerifyMessageResponse{
		User: &user_server.UserMessageResponse{
			Login: userEntity.Login,
			Id:    userEntity.ID,
			Roles: nil,
		},
	}, nil
}

func (s *GrpcUserServer) SignIn(
	_ context.Context,
	request *user_server.SignInMessage,
) (*user_server.SignInResponse, error) {
	form := s.userFormBuilder.SignIn(request)
	signInChan := make(chan interface{})
	var signInError error
	var userResponse *entities.User
	var token string
	go func() {
		userResponse, token, signInError = s.authService.SignIn(form, signInChan)
	}()
	<-signInChan
	if signInError != nil {
		return nil, signInError
	}
	return &user_server.SignInResponse{
		Id:    userResponse.ID,
		Token: token,
	}, nil
}
